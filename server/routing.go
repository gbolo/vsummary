package server

import (
"net/http"

"github.com/gorilla/mux"
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
	Route{
		"Vm",
		"POST",
		appendRequestPrefix("/vm"),
		handlerVm,
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
		var handler http.Handler

		handler = route.HandlerFunc
		//handler = accessLog(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}