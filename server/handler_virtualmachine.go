package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"gopkg.in/go-playground/validator.v9"
	//"github.com/thoas/stats"
)

func handlerVirtualmachine(w http.ResponseWriter, req *http.Request) {

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
	var VmReq []common.VirtualMachine
	err = json.Unmarshal(reqBody, &VmReq)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	// TODO: fail if any vm is invalid?
	validate := validator.New()
	for _, vm := range VmReq {
		err = validate.Struct(vm)
		if err != nil {
			log.Errorf("failed to validate body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// insert to backend
	err = backend.InsertVirtualmachines(VmReq)
	if err != nil {
		log.Errorf("failed to insert vm: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
