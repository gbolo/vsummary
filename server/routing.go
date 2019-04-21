package server

import (
	"net/http"

	"github.com/gbolo/vsummary/common"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//type vSummeryHandlerFunc func(Server, http.ResponseWriter, *http.Request)

type Routes []Route

// all defined api routes
var routes = Routes{

	// vSummary API endpoints
	Route{
		"VirtualMachine",
		"POST",
		common.EndpointVirtualMachine,
		handlerVirtualmachine,
	},
	Route{
		"Datacenter",
		"POST",
		common.EndpointDatacenter,
		handlerDatacenter,
	},
	Route{
		"Cluster",
		"POST",
		common.EndpointCluster,
		handlerCluster,
	},
	Route{
		"Esxi",
		"POST",
		common.EndpointESXi,
		handlerEsxi,
	},
	Route{
		"ResourcePool",
		"POST",
		common.EndpointResourcepool,
		handlerResourcepool,
	},
	Route{
		"Datastore",
		"POST",
		common.EndpointDatastore,
		handlerDatastore,
	},
	Route{
		"VDisks",
		"POST",
		common.EndpointVDisk,
		handlerVDisks,
	},
	Route{
		"VNics",
		"POST",
		common.EndpointVNIC,
		handlerVNics,
	},
	Route{
		"Folders",
		"POST",
		common.EndpointFolder,
		handlerFolders,
	},
	Route{
		"VSwitch",
		"POST",
		common.EndpointVSwitch,
		handlerVswitch,
	},
	Route{
		"vCenter",
		"POST",
		common.EndpointVCenter,
		handlerVcenter,
	},
	Route{
		"Poller",
		"POST",
		common.EndpointPoller,
		handlerPoller,
	},
	Route{
		"DeletePoller",
		"DELETE",
		common.EndpointPoller + "/{id}",
		handlerDeletePoller,
	},
	Route{
		"PollNow",
		"POST",
		common.EndpointPoller + "/{id}/pollnow",
		handlerPollNow,
	},

	// vSummary UI endpoints
	Route{
		"UiIndex",
		"GET",
		"/",
		handlerUiIndex,
	},
	Route{
		"UiVirtualmachines",
		"GET",
		"/ui/virtualmachines",
		handlerUiVirtualmachines,
	},
	Route{
		"UiESXi",
		"GET",
		"/ui/esxi",
		handlerUiEsxi,
	},
	Route{
		"UiPortgroups",
		"GET",
		"/ui/portgroups",
		handlerUiPortgroup,
	},
	Route{
		"UiDatastores",
		"GET",
		"/ui/datastores",
		handlerUiDatastore,
	},
	Route{
		"UiVNics",
		"GET",
		"/ui/vnics",
		handlerUiVNic,
	},
	Route{
		"UiVDisks",
		"GET",
		"/ui/vdisks",
		handlerUiVDisk,
	},
	Route{
		"UiClusters",
		"GET",
		"/ui/clusters",
		handlerUiCluster,
	},
	Route{
		"UiVCenters",
		"GET",
		"/ui/vcenters",
		handlerUiVCenters,
	},
	Route{
		"UiPoller",
		"GET",
		"/ui/pollers",
		handlerUiPoller,
	},
	Route{
		"UiFormAddPoller",
		"GET",
		"/ui/form/poller",
		handlerUiFormPoller,
	},
	Route{
		"UiFormEditPoller",
		"GET",
		"/ui/form/editpoller/{id}",
		handlerUiFormEditPoller,
	},
	Route{
		"UiFormRemovePoller",
		"GET",
		"/ui/form/removepoller/{id}",
		handlerUiFormRemovePoller,
	},

	// Datatables API endpoints
	Route{
		"DtVirtualMachines",
		"POST",
		"/api/dt/virtualmachines",
		handlerDtVirtualMachine,
	},
	Route{
		"DtESXi",
		"POST",
		"/api/dt/esxi",
		handlerDtEsxi,
	},
	Route{
		"DtPortgroups",
		"POST",
		"/api/dt/portgroups",
		handlerDtPortgroup,
	},
	Route{
		"DtDatastores",
		"POST",
		"/api/dt/datastores",
		handlerDtDatastore,
	},
	Route{
		"DtVNics",
		"POST",
		"/api/dt/vnics",
		handlerDtVNic,
	},
	Route{
		"DtVDisks",
		"POST",
		"/api/dt/vdisks",
		handlerDtVDisk,
	},
	Route{
		"DtClusters",
		"POST",
		clusterView.AjaxEndpoint,
		handlerDtCluster,
	},
	Route{
		"DtVCenters",
		"POST",
		vcenterView.AjaxEndpoint,
		handlerDtVCenter,
	},
	//Route{
	//	"Stats",
	//	"GET",
	//	appendRequestPrefix("/stats"),
	//	handlerStats,
	//},
}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		// add routes to mux
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// add route to mux to handle static files
	staticPath := viper.GetString("server.static_files_dir")
	if staticPath == "" {
		staticPath = "./static"
	}

	router.
		Methods("GET").
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	return router
}
