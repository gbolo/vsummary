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

var vcenterView = UiView{
	Title:        "vCenters",
	AjaxEndpoint: "/api/dt/vcenter",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"host", "Address"},
		{"datacenters", "Datacenters"},
		{"clusters", "Clusters"},
		{"esxi_hosts", "Esxi Hosts"},
		{"esxi_cpu", "Total CPU Cores"},
		{"esxi_memory", "Total Memory"},
		{"vms_on", "VMs PoweredOn"},
		{"vms", "VMs"},
		{"vms_vcpu_on", "VMs PoweredOn vCPU"},
		{"vms_memory_on", "VMs PoweredOn Memory"},
		//{"vcenter_fqdn", "vCenter"},
	},
}

func handlerUiVCenters(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiCluster")

	// output the page
	writeSummaryPage(w, &vcenterView)

	return
}

func handlerDtVCenter(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_vcenter")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["esxi_memory"] = common.BytesHumanReadable(row["esxi_memory"])
		row["vms_memory_on"] = common.BytesHumanReadable(row["vms_memory_on"])
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

func handlerVcenter(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleVcenter")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var reqStruct common.VCenter
	err = json.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	validate := validator.New()

	err = validate.Struct(reqStruct)
	if err != nil {
		log.Errorf("failed to validate body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// insert to backend
	err = backend.InsertVcenter(reqStruct)
	if err != nil {
		log.Errorf("failed to insert vcenter: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
