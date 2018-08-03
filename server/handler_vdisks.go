package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"gopkg.in/go-playground/validator.v9"
)

var vDiskView = UiView{
	Title:        "vDisks",
	AjaxEndpoint: "/api/dt/vdisks",
	TableHeaders: []tableColumnMap{
		{"vm_name", "VM"},
		{"name", "Disk"},
		{"capacity_bytes", "Capacity"},
		{"path", "Path"},
		{"thin_provisioned", "ThinProvisioned"},
		{"vm_power_state", "VMpowerstate"},
		{"datastore_name", "Datastore"},
		{"datastore_type", "DatastoreType"},
		{"esxi_name", "ESXi"},
		{"vcenter_fqdn", "vCenter"},
	},
}

func handlerUiVDisk(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiVDisk")

	// output the page
	writeSummaryPage(w, &vDiskView)

	return
}

func handlerDtVDisk(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_vdisk")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["capacity_bytes"] = bytesHumanReadable(row["capacity_bytes"])
		row["vm_power_state"] = decorateCell(row["vm_power_state"])
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

func handlerVDisks(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerVDisks")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var reqStruct []common.VDisk
	err = json.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	// TODO: fail if any obj is invalid?
	validate := validator.New()
	for _, o := range reqStruct {
		err = validate.Struct(o)
		if err != nil {
			log.Errorf("failed to validate body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// insert to backend
	err = backend.InsertVDisks(reqStruct)
	if err != nil {
		log.Errorf("failed to insert vdisks: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
