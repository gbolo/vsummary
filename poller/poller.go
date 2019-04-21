package poller

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/crypto"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

var log = logging.MustGetLogger("vsummary")

const (
	// the default interval we refresh the list of pollers from the backend db
	defaultRefreshInterval = 60 * time.Second
)

func init() {
	// seed the random package with current time with nano-second precision
	rand.Seed(time.Now().UTC().UnixNano())
}

// Poller can poll a single endpoint
type Poller struct {
	Name         string
	Enabled      bool
	VmwareClient *govmomi.Client
	Config       *common.Poller
}

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

// NewEmptyPoller returns an empty Poller
func NewEmptyPoller() *Poller {
	return &Poller{}
}

// NewPoller returns a Poller based from a common.Poller
func NewPoller(c common.Poller) (p *Poller) {
	p = NewEmptyPoller()
	p.Configure(c)
	return
}

// Configure allows you to configure a Poller based from a common.Poller
func (p *Poller) Configure(c common.Poller) {
	if c.PlainTextPassword == "" {
		decryptedPassword, err := crypto.Decrypt(c.EncryptedPassword)
		if err != nil {
			log.Warningf("failed to decrypt password for: %s", c.VcenterHost)
			return
		}
		c.PlainTextPassword = decryptedPassword
	}
	p.Config = &c
	p.Config.VcenterURL = fmt.Sprintf("https://%s/sdk", c.VcenterHost)
	p.Name = c.VcenterName
	p.Enabled = c.Enabled
}

// Connect establishes a connection to the vmware endpoint
func (p *Poller) Connect(ctx *context.Context) (err error) {

	// construct vCenter URL
	if p.Config.VcenterURL == "" {
		err = errors.New("vCenter URL cannot be empty")
		return
	}

	vUrl, err := soap.ParseURL(p.Config.VcenterURL)
	if err != nil {
		return
	}
	vUrl.User = url.UserPassword(p.Config.Username, p.Config.PlainTextPassword)

	// configure default vmware client
	p.VmwareClient, err = govmomi.NewClient(*ctx, vUrl, true)
	if err != nil {
		log.Errorf("error setting vmware client: %v", err)
		return
	}

	// load any defined CAs
	caFiles := viper.GetString("poller.vcenter_cafile")
	if caFiles != "" {
		errLoadingCAs := p.VmwareClient.SetRootCAs(caFiles)
		if errLoadingCAs != nil {
			log.Errorf("error loading custom CA(s): %v", errLoadingCAs)
		} else {
			// since we were able to load CA(s), we should now enforce it
			log.Infof("loaded additional CA file(s) to validate vCenter URL(s): %v", caFiles)
			p.VmwareClient.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = false
		}
	}
	return
}

// GetPollResults will return pollResults along with all errors encountered during the polling
func (p *Poller) GetPollResults() (r pollResults, errors []error) {

	// create context and connect to vsphere
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	err := p.Connect(&ctx)
	if err != nil {
		log.Errorf("failed to connect to: %s %s ", p.Config.VcenterURL, err)
		appendIfError(&errors, err)
		return
	}
	defer p.VmwareClient.Logout(ctx)

	// do the polls
	r.Vcenter, err = p.GetVcenter()
	if err != nil {
		// if we can't get vcenter info, we might as well just quit here...
		appendIfError(&errors, err)
		return
	}

	// if we got past the vcenter poll, we can do the rest now
	r.Esxi, _, r.VSwitch, r.StdPortgroup, err = p.GetEsxi()
	appendIfError(&errors, err)
	r.Virtualmachine, r.VDisk, r.Vnic, err = p.GetVirtualMachines()
	appendIfError(&errors, err)
	r.Datacenter, err = p.GetDatacenters()
	appendIfError(&errors, err)
	r.Cluster, err = p.GetClusters()
	appendIfError(&errors, err)
	r.Datastore, err = p.GetDatastores()
	appendIfError(&errors, err)
	r.Dvs, err = p.GetDVS()
	appendIfError(&errors, err)
	r.DvsPortGroup, err = p.GetDVSPortgroups()
	appendIfError(&errors, err)
	r.ResourcePool, err = p.GetResourcepools()
	appendIfError(&errors, err)
	r.Folder, err = p.GetFolders()
	appendIfError(&errors, err)

	return
}

// TestConnection will test if we can successfully log into the provided vcenter server
func TestConnection(p common.Poller) (err error) {
	poller := NewPoller(p)

	// create context and connect to vsphere
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = poller.Connect(&ctx); err != nil {
		log.Errorf("failed to connect to: %s %s ", poller.Config.VcenterURL, err)
		return
	}
	return
}

// randomizedWait sleeps for a random amount of seconds between
// the specified upper and lower integers provided.
func randomizedWait(lower, upper int) {
	s := rand.Intn(upper-lower) + lower
	time.Sleep(time.Duration(s) * time.Second)
}

// appendIfError adds an err to the slice if err is not nil
func appendIfError(errors *[]error, err error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}
