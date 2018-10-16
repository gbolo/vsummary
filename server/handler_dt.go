package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gbolo/vsummary/common"
)

func handlerDtVirtualMachine(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_vm")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["status"] = decorateCell(row["status"])
		row["esxi_status"] = decorateCell(row["esxi_status"])
		row["power_state"] = decorateCell(row["power_state"])
		row["guest_tools_running"] = decorateCell(row["guest_tools_running"])
		row["memory_bytes"] = common.BytesHumanReadable(row["memory_bytes"])
		row["memory_mb"] = common.MegaBytesHumanReadable(row["memory_mb"])
		row["stat_guest_memory_usage"] = common.MegaBytesHumanReadable(row["stat_guest_memory_usage"])
		row["stat_host_memory_usage"] = common.MegaBytesHumanReadable(row["stat_host_memory_usage"])
		row["stat_cpu_usage"] = row["stat_cpu_usage"] + " MHz"
		row["stat_uptime_sec"] = common.SecondsToHuman(row["stat_uptime_sec"])
		row["folder"] = common.SetDefaultValue(row["folder"], "None")
		row["cluster"] = common.SetDefaultValue(row["cluster"], "None")
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

func handlerDtEsxi(w http.ResponseWriter, req *http.Request) {
	handlerDatatables(w, req, "view_esxi")
}

func handlerDtPortgroup(w http.ResponseWriter, req *http.Request) {
	handlerDatatables(w, req, "view_portgroup")
}

func handlerDtDatastore(w http.ResponseWriter, req *http.Request) {
	dtResponse, err := getDatatablesResponse(req, "view_datastore")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
	}

	// loop through data and make modifications
	data := dtResponse.Data[:0]
	for _, row := range dtResponse.Data {
		row["capacity_bytes"] = common.BytesHumanReadable(row["capacity_bytes"])
		row["free_bytes"] = common.BytesHumanReadable(row["free_bytes"])
		data = append(data, row)
	}

	// write a response
	dtResponse.Data = data
	b, _ := json.MarshalIndent(dtResponse, "", "  ")
	fmt.Fprintf(w, string(b))
}

// generic handler that will write the http response itself
func handlerDatatables(w http.ResponseWriter, req *http.Request, dbTable string) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "dt api "+dbTable)

	// ensure the request is a proper datatables request
	di, err := ParseDatatablesRequest(req)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
		return
	}

	// get the data for the request
	di.SetDbX(backend.GetDB())
	response, err := di.fetchDataForResponse(dbTable)
	if err != nil {
		log.Errorf("error getting datatables response: %v", err)
		return
	}

	// write a response
	b, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintf(w, string(b))
	return
}

// retrieves the DataTablesResponse based on the request body
func getDatatablesResponse(req *http.Request, dbTable string) (dtResponse DataTablesResponse, err error) {
	// log time on debug
	defer common.ExecutionTime(time.Now(), "dt api "+dbTable)

	// ensure the request is a proper datatables request
	di, err := ParseDatatablesRequest(req)
	if err != nil {
		return
	}

	// get the data for the request
	di.SetDbX(backend.GetDB())
	dtResponse, err = di.fetchDataForResponse(dbTable)
	if err != nil {
		log.Errorf("error getting datatables response: %v", err)
		return
	}
	return
}

// this will add some html styling to string
func decorateCell(value string) string {
	switch strings.ToLower(value) {
	case "green":
		return "<span class=\"label label-pill label-success\">green</span>"
	case "yellow":
		return "<span class=\"label label-pill label-warning\">yellow</span>"
	case "red":
		return "<span class=\"label label-pill label-danger\">red</span>"
	case "poweredon":
		return "<span class=\"label label-pill label-success\">poweredOn</span>"
	case "poweredoff":
		return "<span class=\"label label-pill label-danger\">poweredOff</span>"
	case "yes", "true", "guesttoolsrunning":
		return "<span class=\"label label-pill label-success\">Yes</span>"
	case "no", "false", "guesttoolsnotrunning":
		return "<span class=\"label label-pill label-danger\">No</span>"
	case "ok":
		return "<span class=\"label label-pill label-success\">OK</span>"
	default:
		return value
	}
}
