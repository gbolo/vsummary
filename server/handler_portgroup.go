package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"gopkg.in/go-playground/validator.v9"
)

var portgroupView = UiView{
	Title:        "PortGroup",
	AjaxEndpoint: "/api/dt/portgroups",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"type", "Type"},
		{"vlan", "Vlan"},
		{"vlan_type", "VlanType"},
		{"vswitch_name", "vSwitch"},
		{"vswitch_type", "vSwitchType"},
		{"vswitch_max_mtu", "vSwitchMTU"},
		{"vnics", "vNics"},
		{"vcenter_fqdn", "vCenter"},
		{"vcenter_short_name", "VC-ENV"},
	},
}

func handlerUiPortgroup(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiPortgroup")

	// output the page
	writeSummaryPage(w, &portgroupView)

	return
}

func handlerDtPortgroup(w http.ResponseWriter, req *http.Request) {
	handlerDatatables(w, req, "view_portgroup")
}

func handlerPortgroups(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerPortgroups")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var reqStruct []common.Portgroup
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
	err = backend.InsertPortgroups(reqStruct)
	if err != nil {
		log.Errorf("failed to insert portgroups: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
