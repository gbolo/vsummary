package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
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
		row["memory_bytes"] = bytesHumanReadable(row["memory_bytes"])
		row["memory_mb"] = megaBytesHumanReadable(row["memory_mb"])
		row["stat_guest_memory_usage"] = megaBytesHumanReadable(row["stat_guest_memory_usage"])
		row["stat_host_memory_usage"] = megaBytesHumanReadable(row["stat_host_memory_usage"])
		row["stat_cpu_usage"] = row["stat_cpu_usage"] + " MHz"
		row["stat_uptime_sec"] = secondsToHuman(row["stat_uptime_sec"])
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
		row["capacity_bytes"] = bytesHumanReadable(row["capacity_bytes"])
		row["free_bytes"] = bytesHumanReadable(row["free_bytes"])
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

// returns a human readable string from any number of bytes
// example: 1855425871872 will return 1.9 TB
func bytesHumanReadable(bytes string) string {
	// ignore numbers after a possible decimal
	bytesSplit := strings.Split(bytes, ".")
	b, err := strconv.ParseInt(bytesSplit[0], 10, 64)
	if err != nil {
		log.Errorf("parse int err: %s", err)
		return "000"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// returns a human readable string from any number of megabytes
func megaBytesHumanReadable(megaBytes string) string {
	// ignore numbers after a possible decimal
	megaBytesSplit := strings.Split(megaBytes, ".")
	b, _ := strconv.ParseInt(megaBytesSplit[0], 10, 64)
	return bytesHumanReadable(fmt.Sprintf("%d", (b * 1000 * 1000)))
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

// converts seconds to days
func secondsToHuman(secondsString string) string {
	seconds, _ := strconv.ParseInt(secondsString, 10, 64)
	days := math.Floor(float64(seconds) / 86400)
	hours := math.Floor(float64(seconds%86400) / 3600)
	minutes := math.Floor(float64(seconds%86400%3600) / 60)

	if seconds == 0 {
		return "nil"
	} else if days < 1 {
		return fmt.Sprintf("%dh, %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%v days", days)
	}
}
