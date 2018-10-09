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

var clusterView = UiView{
	Title:        "Clusters",
	AjaxEndpoint: "/api/dt/cluster",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"status", "Status"},
		{"ha_enabled", "HA Enabled"},
		{"drs_enabled", "DRS Enabled"},
		{"total_vmotions", "vMotions"},
		{"num_hosts", "Esxi Hosts"},
		{"avg_memory_per_host", "Avg Memory per Host"},
		{"total_memory_bytes", "Total Memory"},
		{"total_memory_used", "Used Memory"},
		{"vms_on", "VMs PoweredOn"},
		{"avg_vcpu_per_vm", "Avg vCPU per VM"},
		{"avg_memory_per_vm", "Avg Memory per VM"},
		//{"vcenter_fqdn", "vCenter"},
	},
}

func handlerUiCluster(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiCluster")

	// output the page
	writeSummaryPage(w, &clusterView)

	return
}

func handlerDtCluster(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_cluster_capacity")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["status"] = decorateCell(row["status"])
		row["ha_enabled"] = decorateCell(row["ha_enabled"])
		row["drs_enabled"] = decorateCell(row["drs_enabled"])
		row["avg_memory_per_host"] = bytesHumanReadable(row["avg_memory_per_host"])
		row["total_memory_bytes"] = bytesHumanReadable(row["total_memory_bytes"])
		row["total_memory_used"] = bytesHumanReadable(row["total_memory_used"])
		row["avg_memory_per_vm"] = bytesHumanReadable(row["avg_memory_per_vm"])
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

func handlerCluster(w http.ResponseWriter, req *http.Request) {

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
	var reqStruct []common.Cluster
	err = json.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	// TODO: fail if any dc is invalid?
	validate := validator.New()
	for _, dc := range reqStruct {
		err = validate.Struct(dc)
		if err != nil {
			log.Errorf("failed to validate body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// insert to backend
	err = backend.InsertClusters(reqStruct)
	if err != nil {
		log.Errorf("failed to insert cluster: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
