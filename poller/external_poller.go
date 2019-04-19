package poller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"

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
	if err == nil {
		e.vSummaryApiUrl = u
	}
	return
}

// constructUrl will create the desired vsummary api url
func (e *ExternalPoller) constructUrl(endpoint string) (urlEndpont string, err error) {
	if e.vSummaryApiUrl != "" && endpoint != "" {
		urlEndpont = fmt.Sprintf("%s%s", e.vSummaryApiUrl, endpoint)
	} else {
		err = fmt.Errorf("vSummaryApiUrl or endpoint is empty: [%s] [%s]", e.vSummaryApiUrl, endpoint)
		return
	}
	_, err = url.ParseRequestURI(urlEndpont)
	return
}

// sendResult does an http post request to the vsummary api server to process the poll result
func (e *ExternalPoller) sendResult(endpoint string, o interface{}) (err error) {
	// convertproccess object to json bytes
	jsonBody, err := json.Marshal(o)
	if err != nil {
		log.Errorf("invalid json %s", err)
		return
	}

	// determine url
	url, err := e.constructUrl(endpoint)
	if err != nil {
		return
	}

	// send request
	log.Debugf("sending results to: %s", url)
	res, err := vSummaryClient.Post(url, "application/json", bytes.NewReader(jsonBody))

	// this means the vsummary server api is unreachable
	if err != nil {
		log.Errorf("vsummary api is unreachable: %s error %s", url, err)
		return
	}

	// we only accept 202 as success
	if res.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("recieved %d response code from %v", res.StatusCode, url)
		return
	}

	// To ensure KeepAlive:
	// Read until Response is complete (i.e. ioutil.ReadAll(rep.Body))
	// Call Body.Close()
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	log.Infof("api call successful: %d %s", res.StatusCode, url)
	return
}

// SendPollResults will attempt to send the polling results to the vSummary API server
func (e *ExternalPoller) SendPollResults(r pollResults) (err []error) {
	appendIfError(&err, e.sendResult(common.EndpointVCenter, r.Vcenter))
	appendIfError(&err, e.sendResult(common.EndpointESXi, r.Esxi))
	appendIfError(&err, e.sendResult(common.EndpointDatastore, r.Datastore))
	appendIfError(&err, e.sendResult(common.EndpointVirtualMachine, r.Virtualmachine))
	appendIfError(&err, e.sendResult(common.EndpointVSwitch, r.VSwitch))
	appendIfError(&err, e.sendResult(common.EndpointVSwitch, r.Dvs))
	// need to insert portgroups here...
	appendIfError(&err, e.sendResult(common.EndpointVNIC, r.Vnic))
	appendIfError(&err, e.sendResult(common.EndpointVDisk, r.VDisk))
	appendIfError(&err, e.sendResult(common.EndpointResourcepool, r.ResourcePool))
	appendIfError(&err, e.sendResult(common.EndpointDatacenter, r.Datacenter))
	appendIfError(&err, e.sendResult(common.EndpointFolder, r.Folder))
	appendIfError(&err, e.sendResult(common.EndpointCluster, r.Cluster))

	return
}

// PollThenSend will poll all endpoints then send results to vSummary API server
func (e *ExternalPoller) PollThenSend() (errs []error) {
	r, errs := e.GetPollResults()
	if len(errs) > 0 {
		log.Warningf(
			"will not send poll results since %d error(s) occurred during polling of: %s",
			len(errs),
			e.Config.URL,
		)
		return
	}
	errs = e.SendPollResults(r)
	if len(errs) > 0 {
		log.Warningf(
			"there were %d errors during sending polling results of: %s",
			len(errs),
			e.Config.URL,
		)
		log.Debugf("polling errors: %v", errs)
	}
	return
}

// Daemonize is a blocking loop which continues to PollThenSend indefinitely
func (e *ExternalPoller) Daemonize() {
	// TODO: global polling interval is use for now.
	// in future versions we can try and support an interval per poller
	t := time.Tick(time.Duration(viper.GetInt("poller.interval")) * time.Minute)
	log.Infof("start interval polling (%dm) of %s", viper.GetInt("poller.interval"), e.Config.URL)

	// this prevents all pollers to go off at the exact same time
	randomizedWait(1, 10)
	e.PollThenSend()

	for {
		select {
		case <-t:
			if e.Enabled {
				// this prevents all pollers to go off at the exact same time
				randomizedWait(1, 120)
				log.Debugf("executing poll of %s", e.Config.URL)
				e.PollThenSend()
			} else {
				log.Infof("stopping polling of %s", e.Config.URL)
				return
			}
		}
	}
}

// GetExternalPollersFromConfig returns preconfigured list of ExternalPoller(s) found in config
func GetExternalPollersFromConfig() (externalPollers []*ExternalPoller) {
	var pollersInConfig []common.Poller
	err := viper.UnmarshalKey("poller.vcenters", &pollersInConfig)
	if err != nil || len(pollersInConfig) < 1 {
		return
	}

	for _, poller := range pollersInConfig {
		poller.Enabled = true
		externalPollers = append(externalPollers, NewExternalPoller(poller))
	}
	return
}
