package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/crypto"
)

func LoadPollers(pollers []common.Poller) {

	log.Debugf("starting %d pollers", len(pollers))

	for _, p := range pollers {

		decryptedPassword, err := crypto.Decrypt(p.Password)
		if err != nil {
			log.Warningf("failed to decrypt password for: %s", p.VcenterHost)
			break
		}

		poller := NewPoller()
		poller.Config = &PollerConfig{
			URL:         fmt.Sprintf("https://%s/sdk", p.VcenterHost),
			UserName:    p.Username,
			Password:    decryptedPassword,
			IntervalMin: p.IntervalMin,
			Insecure:    true,
		}

		log.Infof("starting poller loop for: %s", p.VcenterHost)
		go pollerLoop(poller)
	}
}

// testing poller loop
func pollerLoop(p *Poller) (err error) {

	// create context and connect to vsphere
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = p.Connect(&ctx)
	if err != nil {
		log.Errorf("failed to connect to: %s %s ", p.Config.URL, err)
		return
	}

	defer p.VmwareClient.Logout(ctx)

	// do a poll immediately
	p.PollAllEndpoints()

	timeout := time.After(60 * time.Minute)
	tick := time.Tick(time.Duration(p.Config.IntervalMin) * time.Minute)
	//tick := time.Tick(10 * time.Second)

	// loop
	log.Debugf("connection to %s successful, polling interval: %d min", p.Config.URL, p.Config.IntervalMin)
	for {
		select {
		case <-timeout:
			// exit when timeout reached
			log.Debug("exiting poller loop")
			return
		case <-tick:
			log.Debug("ticker reached, polling now")
			p.PollAllEndpoints()
		}
	}

}

// TestConnection will test if we can successfully log into the provided vcenter server
func TestConnection(p PollerConfig) (err error) {
	poller := NewPoller()
	poller.Config = &p

	// create context and connect to vsphere
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = poller.Connect(&ctx); err != nil {
		log.Errorf("failed to connect to: %s %s ", poller.Config.URL, err)
		return
	}
	return
}
