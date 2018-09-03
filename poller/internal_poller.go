package poller

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/crypto"
	"github.com/gbolo/vsummary/db"
)

// pollResults stores the results of a full poll
type pollResults struct {
	Vcenter        common.VCenter
	Esxi           []common.Esxi
	Virtualmachine []common.VirtualMachine
	Datastore      []common.Datastore
	VSwitch        []common.VSwitch
	StdPortgroup   []common.Portgroup
	Dvs            []common.VSwitch
	DvsPortGroup   []common.Portgroup
	Vnic           []common.VNic
	VDisk          []common.VDisk
	ResourcePool   []common.ResourcePool
	Datacenter     []common.Datacenter
	Folder         []common.Folder
	Cluster        []common.Cluster
}

const (
	// the default interval a poller will poll endpoints
	defaultPollingInterval = 30 * time.Minute
	// the default interval we refresh the list of pollers from the backend db
	defaultRefreshInterval = 60 * time.Second
)

func init() {
	// seed the random package with current time with nano-second precision
	rand.Seed(time.Now().UTC().UnixNano())
}

// poller with a channel to send a stop signal
type channelPoller struct {
	Poller      common.Poller
	stopPolling chan bool
}

// internal poller that contains a list of pollers as well as a backend db
type internalPoller struct {
	ActivePollers []*channelPoller
	backend       db.Backend
}

// SetBackend allows internalPoller to connect to backend database
func (i *internalPoller) SetBackend(backend db.Backend) {
	i.backend = backend
}

// refreshPollers gets a list of pollers from backend database
// then populates internalPoller list of ActivePollers with it.
func (i *internalPoller) refreshPollers() {
	pollers, err := i.backend.GetPollers()
	if err != nil {
		return
	}

	for _, p := range pollers {
		i.addIfUnique(channelPoller{
			Poller:      p,
			stopPolling: make(chan bool),
		})
	}
}

// addIfUnique will spawn a new poller thread for a given poller, if one doe not already exist
// it will also stop a running poller if it notices that poller should be disabled
func (i *internalPoller) addIfUnique(c channelPoller) {
	spawnPoller := true
	uniquePoller := true
	for k, p := range i.ActivePollers {
		// TODO: instead of host, we should use vcenter UUID to determine if it's truly unique
		if p.Poller.VcenterHost == c.Poller.VcenterHost {
			uniquePoller = false
			spawnPoller = false
			// stop the poller if it marked as disabled
			if c.Poller.Enabled == false && p.Poller.Enabled {
				log.Infof("poller state has changed to disabled for %s", p.Poller.VcenterName)
				i.ActivePollers[k].Poller.Enabled = false
				i.ActivePollers[k].stopPolling <- true
			}
			// spawn a new poller since it was disabled
			if c.Poller.Enabled && p.Poller.Enabled == false {
				log.Infof("poller state has changed to enabled for %s", p.Poller.VcenterName)
				i.ActivePollers[k].Poller.Enabled = true
				spawnPoller = true
			}
			continue
		}
	}

	if spawnPoller {
		if uniquePoller {
			log.Infof("spawning new poller for %s", c.Poller.VcenterName)
		} else {
			log.Infof("respawning poller for %s", c.Poller.VcenterName)
		}
		i.ActivePollers = append(i.ActivePollers, &c)
		// spwan a go routine for this poller
		go c.Daemonize(i.backend)
	}
}

// Daemonize will take a poller and poll it periodically until either
// the channel is closed or poller is marked as disabled in database.
func (c *channelPoller) Daemonize(b db.Backend) {
	t := time.Tick(defaultPollingInterval)
	log.Infof("start polling of %s", c.Poller.VcenterName)
	// this prevents all pollers to go off at the exact same time
	randomizedWait(1, 120)
	DoInternalPoll(c.Poller, b)

	// start infinite loop until we receive a false from our channel
	for {
		select {
		case <-t:
			if c.Poller.Enabled {
				// this prevents all pollers to go off at the exact same time
				randomizedWait(1, 120)
				log.Debugf("executing poll of %s", c.Poller.VcenterName)
				DoInternalPoll(c.Poller, b)
			} else {
				log.Infof("stopping polling of %s", c.Poller.VcenterName)
				return
			}
		case <-c.stopPolling:
			log.Infof("channel signal received: stop polling of %s", c.Poller.VcenterName)
			return
		}
	}
}

// RunInternalPoller is a blocking loop. This should only be executed once.
// refreshing of the pollers is also handled in this function.
func RunInternalPoller(backend db.Backend) {
	i := internalPoller{backend: backend}
	tick := time.Tick(defaultRefreshInterval)
	i.refreshPollers()
	// refresh pollers forever
	for {
		select {
		case <-tick:
			log.Debug("refreshing pollers from backend")
			i.refreshPollers()
		}
	}
}

// randomizedWait sleeps for a random amount of seconds between
// the specified upper and lower integers provided.
func randomizedWait(lower, upper int) {
	s := rand.Intn(upper-lower) + lower
	time.Sleep(time.Duration(s) * time.Second)
}

// storePollResults will send results directly to backend db and not to vsummary API server
func storePollResults(r pollResults, b db.Backend) (err []error) {

	appendIfError(&err, b.InsertVcenter(r.Vcenter))
	appendIfError(&err, b.InsertEsxi(r.Esxi))
	appendIfError(&err, b.InsertDatastores(r.Datastore))
	appendIfError(&err, b.InsertVirtualmachines(r.Virtualmachine))
	appendIfError(&err, b.InsertVSwitch(r.VSwitch))
	appendIfError(&err, b.InsertVSwitch(r.Dvs))
	// need to insert portgroups here...
	appendIfError(&err, b.InsertVNics(r.Vnic))
	appendIfError(&err, b.InsertVDisks(r.VDisk))
	appendIfError(&err, b.InsertResourcepools(r.ResourcePool))
	appendIfError(&err, b.InsertDatacenters(r.Datacenter))
	appendIfError(&err, b.InsertFolders(r.Folder))
	appendIfError(&err, b.InsertClusters(r.Cluster))

	return
}

// appendIfError adds an err to the slice if err is not nil
func appendIfError(errors *[]error, err error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}

// does full poll from a common.Poller
func DoInternalPoll(p common.Poller, b db.Backend) (err error) {

	decryptedPassword, err := crypto.Decrypt(p.Password)
	if err != nil {
		log.Warningf("failed to decrypt password for: %s", p.VcenterHost)
		return
	}

	poller := NewPoller()
	poller.Config = &PollerConfig{
		URL:         fmt.Sprintf("https://%s/sdk", p.VcenterHost),
		UserName:    p.Username,
		Password:    decryptedPassword,
		IntervalMin: p.IntervalMin,
		Insecure:    true,
	}

	// create context and connect to vsphere
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = poller.Connect(&ctx)
	if err != nil {
		log.Errorf("failed to connect to: %s %s ", poller.Config.URL, err)
		return
	}

	defer poller.VmwareClient.Logout(ctx)

	pollResults, _ := poller.GetPollResults()
	errors := storePollResults(pollResults, b)
	if len(errors) > 0 {
		log.Warningf("there were %v error(s) storing poll results to backed", len(errors))
		log.Debugf("storePollResults errors: %v", errors)
	}

	return
}
