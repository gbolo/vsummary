package poller

import (
	"context"
	"errors"
	"net/url"

	"github.com/op/go-logging"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

var log = logging.MustGetLogger("vsummary")

const apiVersion = "2"

type PollerConfig struct {
	URL         string
	UserName    string
	Password    string
	IntervalMin int

	// Don't validate TLS Cert
	Insecure bool
}

type Poller struct {
	VmwareClient *govmomi.Client
	Config       *PollerConfig
}

func NewPoller() *Poller {
	return &Poller{}
}

func (p *Poller) Configure(config *PollerConfig) {
	p.Config = config
}

func (p *Poller) Connect(ctx *context.Context) (err error) {

	if p.Config.URL == "" {
		err = errors.New("vsphere URL cannot be empty")
		return
	}

	vUrl, err := soap.ParseURL(p.Config.URL)
	if err != nil {
		return
	}

	vUrl.User = url.UserPassword(p.Config.UserName, p.Config.Password)

	p.VmwareClient, err = govmomi.NewClient(*ctx, vUrl, p.Config.Insecure)
	return

}

func (p *Poller) PollAllEndpoints() {

	for k, _ := range vSummaryEndpoints {
		logPollingResult(p.PollThenSend(k))
	}
}

func logPollingResult(err error) {
	if err == nil {
		log.Info("poll completed successfully")
	} else {
		log.Warningf("poll was not successful: %s", err)
	}
}
