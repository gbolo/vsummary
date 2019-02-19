package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/poller"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
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

func handlerUiFormEditPoller(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pollerId := vars["id"]
	log.Debugf("pollerId: %s", pollerId)

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiFormEditPoller")

	// output the page
	writePollerEditPage(w, "form_edit_poller", pollerId)

	return
}

func handlerUiFormRemovePoller(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pollerId := vars["id"]
	log.Debugf("pollerId: %s", pollerId)

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiFormRemovePoller")

	// output the page
	writePollerEditPage(w, "form_remove_poller", pollerId)

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

func handlerAddPoller(w http.ResponseWriter, req *http.Request) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerAddPoller")

	// parse the form
	req.ParseForm()

	// convert checkbox value
	checkboxEnabled := false
	if req.FormValue("enabled") == "on" {
		checkboxEnabled = true
	}

	if err := poller.TestConnection(poller.PollerConfig{
		URL:      req.FormValue("host"),
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
		Username:    req.FormValue("user"),
		Password:    req.FormValue("pass"),
		Enabled:     checkboxEnabled,
	}); err != nil {
		log.Errorf("could not add poller: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf("Could not add poller: %s", err))
		return
	}

	fmt.Fprint(w, "Successfuly tested and added connection")
}

func handlerDeletePoller(w http.ResponseWriter, req *http.Request) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerDeletePoller")

	// get vars from request to determine environment
	vars := mux.Vars(req)
	pollerID := vars["id"]

	// validate poller id
	if len(pollerID) != 12 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "poller ID is not correctly specified")
		return
	}

	// remove poller
	if err := backend.RemovePoller(pollerID); err != nil {
		log.Errorf("remove poller err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "There was an error removing poller. See logs")
		return
	}

	fmt.Fprint(w, "Successfuly removed poller")
}