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

// can be used to store both internal and external pollers
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
	var reqPoller common.Poller
	err = json.Unmarshal(reqBody, &reqPoller)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	validate := validator.New()

	err = validate.Struct(reqPoller)
	if err != nil {
		log.Errorf("failed to validate body: %s", err)
		fmt.Fprint(w, "ALL fields must be populated")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqPoller.PlainTextPassword == "" && reqPoller.Internal {
		log.Error("poller cannot be marked as internal with an empty plain_password field")
		fmt.Fprint(w, "Password field MUST be set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// test connection if it's internal
	if reqPoller.Internal {
		if err := poller.TestConnection(reqPoller); err != nil {
			log.Errorf("could not connect to vCenter: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprintf("Could not connect to vCenter: %s", err))
			return
		}
	}

	// insert to backend
	err = backend.InsertPoller(reqPoller)
	if err != nil {
		log.Errorf("failed to insert poller: %s", err)
		fmt.Fprint(w, "Failed to insert poller. (database error, check logs)")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, `{"status": "OK", "message": "successful tested and added poller"}`)
	return
}

func handlerDeletePoller(w http.ResponseWriter, req *http.Request) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerDeletePoller")

	// get vars from request to determine poller id
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

func handlerPollNow(w http.ResponseWriter, req *http.Request) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerPollNow")

	// get vars from request to determine poller id
	vars := mux.Vars(req)
	pollerID := vars["id"]

	// validate poller id
	// TODO: we should be able to return 404 if it doesn't exist
	if len(pollerID) != 12 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "poller ID is not correctly specified")
		return
	}

	// poll the poller now
	errs := poller.BuiltInCollector.PollPollerById(pollerID)
	if len(errs) > 0 {
		log.Errorf("PollPollerById produced errors for poller: %s", pollerID)
	}

	fmt.Fprint(w, "OK")
}
