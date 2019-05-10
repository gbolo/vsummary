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

var vNicView = UiView{
	Title:        "vNics",
	AjaxEndpoint: "/api/dt/vnics",
	TableHeaders: []tableColumnMap{
		{"vm_name", "VM"},
		{"name", "Device"},
		{"mac", "MAC"},
		{"type", "Tpye"},
		{"connected", "Connected"},
		{"status", "Status"},
		{"esxi_name", "ESXi"},
		{"vlan", "Vlan(s)"},
		{"portgroup_name", "Portgroup"},
		{"vswitch_name", "vSwitch"},
		{"vswitch_type", "Backing"},
		{"vswitch_max_mtu", "MaxMTU"},
		{"vcenter_fqdn", "vCenter"},
		{"vcenter_short_name", "Site"},
	},
}

func handlerUiVNic(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiVNic")

	// output the page
	writeSummaryPage(w, &vNicView)

	return
}

func handlerDtVNic(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_vnic")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["connected"] = decorateCell(row["connected"])
		row["status"] = decorateCell(row["status"])
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

func handlerVNics(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerVNics")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var reqStruct []common.VNic
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
	err = backend.InsertVNics(reqStruct)
	if err != nil {
		log.Errorf("failed to insert vnics: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
