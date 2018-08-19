package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/poller"
	"gopkg.in/go-playground/validator.v9"
	//"github.com/thoas/stats"
	//"github.com/codegangsta/negroni"
	"fmt"
)

func handlerUiPoller(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiPoller")

	// output the page
	writePollerPage(w, "pollers")

	return
}

func handlerUiFormPoller(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiFormPoller")

	// output the page
	writePollerPage(w, "form_add_poller")

	return
}

func handlerPoller(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handle")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var poller common.Poller
	err = json.Unmarshal(reqBody, &poller)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	validate := validator.New()

	err = validate.Struct(poller)
	if err != nil {
		log.Errorf("failed to validate body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// insert to backend
	err = backend.InsertPoller(poller)
	if err != nil {
		log.Errorf("failed to insert poller: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}


func handlerPutVcenter(w http.ResponseWriter, req *http.Request) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerPutVcenter")

	// parse the form
	req.ParseForm()

	if err := poller.TestConnection(poller.PollerConfig{
		URL: req.FormValue("host"),
		UserName: req.FormValue("user"),
		Password: req.FormValue("pass"),
		Insecure: true,
	}); err != nil {
		log.Errorf("could not connect to vCenter: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf("Could not connect to vCenter: %s", err))
		return
	}

	if err := backend.InsertPoller(common.Poller{
		VcenterHost: req.FormValue("host"),
		VcenterName: req.FormValue("short_name"),
		Username: req.FormValue("user"),
		Password: req.FormValue("pass"),
	}); err != nil {
		log.Errorf("could not add poller: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf("Could not add poller: %s", err))
		return
	}

	fmt.Fprint(w, "Successfuly tested and added connection")


}