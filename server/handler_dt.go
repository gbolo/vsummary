package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbolo/vsummary/common"
)

func handlerDtVirtualMachine(w http.ResponseWriter, req *http.Request) {
	handlerDatatables(w, req, "view_vm")
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
	b, err := strconv.ParseInt(bytes, 10, 64)
	if err != nil {
		return "000"
	}
	const unit = 1000
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
