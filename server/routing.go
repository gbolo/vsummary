package server

import (
	"fmt"
	"net/http"

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
		appendRequestPrefix("/virtualmachine"),
		handlerVirtualmachine,
	},
	Route{
		"Datacenter",
		"POST",
		appendRequestPrefix("/datacenter"),
		handlerDatacenter,
	},
	Route{
		"Cluster",
		"POST",
		appendRequestPrefix("/cluster"),
		handlerCluster,
	},
	Route{
		"Esxi",
		"POST",
		appendRequestPrefix("/esxi"),
		handlerEsxi,
	},
	Route{
		"ResourcePool",
		"POST",
		appendRequestPrefix("/resourcepool"),
		handlerResourcepool,
	},
	Route{
		"VDisks",
		"POST",
		appendRequestPrefix("/vdisk"),
		handlerVDisks,
	},
	Route{
		"VNics",
		"POST",
		appendRequestPrefix("/vnic"),
		handlerVNics,
	},
	Route{
		"Folders",
		"POST",
		appendRequestPrefix("/folder"),
		handlerFolders,
	},
	Route{
		"VSwitch",
		"POST",
		appendRequestPrefix("/vswitch"),
		handlerVswitch,
	},
	Route{
		"vCenter",
		"POST",
		appendRequestPrefix("/vcenter"),
		handlerVcenter,
	},
	Route{
		"Poller",
		"POST",
		appendRequestPrefix("/poller"),
		handlerPoller,
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
	//Route{
	//	"Stats",
	//	"GET",
	//	appendRequestPrefix("/stats"),
	//	handlerStats,
	//},
}

// appends prefix to route path
func appendRequestPrefix(route string) string {

	return fmt.Sprintf("/api/v%s%s", apiVersion, route)
}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		//handler = accessLog(handler, route.Name)

		// add routes to mux
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
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
