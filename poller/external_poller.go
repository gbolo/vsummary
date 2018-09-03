package poller

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gbolo/vsummary/common"
)

var (
	// global http client for calls to vsummary server api
	vSummaryClient *http.Client
)

func init() {
	// set sane defaults for vSummaryClient HTTP client
	vSummaryClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5,
			DisableCompression:    true,
			IdleConnTimeout:       10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

// ExternalPoller extends Poller with functionality relevant to
// sending results to a vSummary API server over http(s).
type ExternalPoller struct {
	Poller
	stopSignal     chan bool
	vSummaryApiUrl string
}

// NewEmptyExternalPoller returns a empty ExternalPoller
func NewEmptyExternalPoller() *ExternalPoller {
	return &ExternalPoller{
		stopSignal: make(chan bool),
	}
}

// NewExternalPoller returns a ExternalPoller based from a common.Poller
func NewExternalPoller(c common.Poller) (e *ExternalPoller) {
	e = NewEmptyExternalPoller()
	e.Configure(c)
	return
}

// SetEndpoint sets the vSummary API server url unless it's invalid
func (e *ExternalPoller) SetApiUrl(u string) (err error) {
	_, err = url.ParseRequestURI(u)
	if err != nil {
		e.vSummaryApiUrl = u
	}
	return
}
