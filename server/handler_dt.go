package server

import (
	"fmt"
	"net/http"
	"time"

	//"github.com/gbolo/go-util/lib/debugging"
	"github.com/gbolo/vsummary/common"
	//"github.com/toebes/go-datatables-serverside"
	"encoding/json"
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
	handlerDatatables(w, req, "view_datastore")
}

func handlerDatatables(w http.ResponseWriter, req *http.Request, dbTable string) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "dt api "+dbTable)

	// test
	di, err := ParseDatatablesRequest(req)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
		return
	}

	di.SetDbX(backend.GetDB())

	response, err := di.fetchDataForResponse(dbTable)
	if err != nil {
		log.Errorf("error getting datatables response: %v", err)
		return
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintf(w, string(b))

	return
}
