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

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleDt")

	// test
	di, err := ParseDatatablesRequest(req)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Errorf("error parsing datatables request: %v", err)
		return
	}

	di.SetDbX(backend.GetDB())

	response, err := di.fetchDataForResponse("vm")
	if err != nil {
		log.Errorf("error getting datatables response: %v", err)
		return
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintf(w,string(b))

	return
}
