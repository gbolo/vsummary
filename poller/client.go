package poller

import (
	"net/http"
	"time"
	"fmt"
	"github.com/spf13/viper"
	"bytes"
	"io"
	"io/ioutil"
)

var vSummaryClient *http.Client

// Initializes the shared vSummaryClient
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
	res, err := vSummaryClient.Post(url,"application/json", bytes.NewReader(jsonBody))

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