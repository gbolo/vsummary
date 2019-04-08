package poller

import (
	"fmt"
	"time"

	"github.com/gbolo/vsummary/db"
)

// internal poller that contains a list of pollers as well as a backend db
type InternalCollector struct {
	ActivePollers []*InternalPoller
	backend       db.Backend
}

// NewEmptyInternalCollector returns an empty InternalCollector
func NewEmptyInternalCollector() *InternalCollector {
	return &InternalCollector{}
}

// SetBackend allows InternalCollector to connect to backend database
func (i *InternalCollector) SetBackend(backend db.Backend) {
	i.backend = backend
}

// addIfUnique will spawn a new poller thread for a given poller, if one doe not already exist
// it will also stop a running poller if it notices that poller should be disabled
func (i *InternalCollector) addIfUnique(p InternalPoller) {
	spawnPoller := true
	uniquePoller := true
	for k, a := range i.ActivePollers {
		// TODO: instead of host, we should use vcenter UUID to determine if it's truly unique
		if a.Config.URL == p.Config.URL {
			uniquePoller = false
			spawnPoller = false
			// stop the poller if it marked as disabled
			if p.Enabled == false && a.Poller.Enabled {
				log.Infof("poller state has changed to disabled for %s", a.Config.URL)
				i.ActivePollers[k].Enabled = false
				i.ActivePollers[k].StopPolling()
			}
			// spawn a new poller since it was disabled
			if p.Enabled && a.Enabled == false {
				log.Infof("poller state has changed to enabled for %s", a.Config.URL)
				i.ActivePollers[k].Enabled = true
				spawnPoller = true
			}
			continue
		}
	}

	if spawnPoller {
		if uniquePoller {
			log.Infof("spawning new poller for %s", p.Config.URL)
		} else {
			log.Infof("respawning poller for %s", p.Config.URL)
		}
		i.ActivePollers = append(i.ActivePollers, &p)
		// spawn a go routine for this poller
		go p.Daemonize()
	}
}

// RefreshPollers gets a list of pollers from backend database
// then populates internalPoller list of ActivePollers with it.
func (i *InternalCollector) RefreshPollers() {
	pollers, err := i.backend.GetPollers()
	log.Debugf("found %d pollers", len(pollers))
	if err != nil {
		log.Errorf("error getting pollers: %v", err)
		return
	}

	// add unique new pollers
	var backendPollerURLs []string
	for _, p := range pollers {
		internalPoller := NewInternalPoller(p)
		internalPoller.SetBackend(i.backend)
		i.addIfUnique(*internalPoller)
		backendPollerURLs = append(backendPollerURLs, fmt.Sprintf("https://%s/sdk", p.VcenterHost))
	}

	// remove pollers that are no longer present or disabled
	i.StopPollersByURL(difference(i.GetActivePollerURLs(), backendPollerURLs))
}

// GetActivePollerURLs returns a list of active pollers by URL
func (i *InternalCollector) GetActivePollerURLs() (urls []string) {
	for _, p := range i.ActivePollers {
		urls = append(urls, p.Config.URL)
	}
	return
}

// StopPollersByURL will stop active pollers that match the list of URLs
func (i *InternalCollector) StopPollersByURL(urls []string) {
	for _, url := range urls {
		for _, p := range i.ActivePollers {
			if p.Config.URL == url && p.Enabled {
				log.Warningf("poller URL is active in memory but no longer listed in backend: %v", url)
				p.StopPolling()
			}
		}
	}
}

// difference returns the elements in a that aren't in b
func difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// Run is a blocking loop. This should only be executed once.
// refreshing of the pollers is also handled in this function.
func (i *InternalCollector) Run() {
	tick := time.Tick(defaultRefreshInterval)
	i.RefreshPollers()
	// refresh pollers forever
	for {
		select {
		case <-tick:
			log.Debug("refreshing pollers from backend")
			i.RefreshPollers()
		}
	}
}

//// findPollerById returns a matching internal poller
//func (i *InternalCollector) findPollerById(id string) (foundPoller *InternalPoller) {
//	for _, p := range i.ActivePollers {
//		//
//	}
//}
