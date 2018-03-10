package poller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/spf13/viper"
)

var (

	// global http client for calls to vsummary server api
	vSummaryClient *http.Client

	// object types and thier endpoints
	vSummaryEndpoints = map[string]string{
		"virtualmachines": "/virtualmachine",
		"clusters":        "/cluster",
		"datacenters":     "/datacenter",
		"esxi":            "/esxi",
	}
)

// Initializes the shared vSummaryClient
// TODO: add error conditions
func initHttpClient() (err error) {

	// return right away if not nil
	if vSummaryClient != nil {
		log.Debug("vSummaryClient already initialized")
		return
	}

	// TODO: add more options and TLS as well
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

	return
}

// sends an api request to vsummary server api
func sendResults(endpoint string, jsonBody []byte) (err error) {

	// endpoint can't be empty
	if endpoint == "" {
		err = fmt.Errorf("endpoint is not defined")
		return
	}

	// init client (no error for now)
	initHttpClient()

	// construct url
	url := fmt.Sprintf("%s/api/v%s%s",
		viper.GetString("poller.url"),
		apiVersion,
		endpoint,
	)

	// send request
	log.Debugf("sending results to: %s", url)
	//log.Debugf("jsonBody: %s", string(jsonBody))
	res, err := vSummaryClient.Post(url, "application/json", bytes.NewReader(jsonBody))

	// this means the vsummary server api is unreachable
	if err != nil {
		log.Errorf("vsummary api is unreachable: %s error %s", url, err)
		return
	}

	// we only accept 202 as success
	if res.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("recieved %d response code from %", res.StatusCode, url)
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

// does a poll then sends the results to the vsummary server api
func (p *Poller) PollThenSend(objectType string) (err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), fmt.Sprintf("pollThenSend: %s", objectType))

	// poll the object type
	var o interface{}

	switch objectType {

	case "virtualmachines":
		o, err = p.GetVirtualMachines()
	case "datacenters":
		o, err = p.GetDatacenters()
	case "clusters":
		o, err = p.GetClusters()
	case "esxi":
		o, err = p.GetEsxi()

	default:
		err = fmt.Errorf("invalid endpoint: %s", objectType)
		return

	}

	if err != nil {
		log.Debugf("failed to poll %s: %s", objectType, err)
		return
	}

	// marshal, then send the results
	log.Infof("poller sending summary of %s", objectType)

	jsonObj, err := json.Marshal(o)
	if err != nil {
		log.Errorf("invalid json %s: %s", objectType, err)
		return
	}

	err = sendResults(vSummaryEndpoints[objectType], jsonObj)
	if err != nil {
		log.Errorf("error sending %s: %s", objectType, err)
		return
	}

	return
}
