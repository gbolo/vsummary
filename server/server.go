package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gbolo/vsummary/db"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	//"github.com/thoas/stats"
	//"github.com/codegangsta/negroni"
)

const apiVersion = "2"

var log = logging.MustGetLogger("vsummary")
var backend *db.Backend

//var serverStats stats.Stats

type Server struct {
	//Version 	string
	HttpServer *http.Server
	//Backend		*db.Backend
}

func Start() (err error) {

	// Init backend database
	backend, err = db.InitBackend()
	if err != nil {
		log.Errorf("failed to connect to database: %s", err)
		return
	}

	// Init Stats
	//serverStats = *stats.New()

	// create routes
	mux := newRouter()

	// get server config
	srv := configureHttpServer(mux)

	// start the server
	log.Info("starting http server")
	err = srv.ListenAndServe()

	return
}

func configureHttpServer(mux *mux.Router) (httpServer *http.Server) {

	// apply standard http server settings
	address := fmt.Sprintf(
		"%s:%s",
		viper.GetString("server.bind_address"),
		viper.GetString("server.bind_port"),
	)

	httpServer = &http.Server{
		Addr: address,
		// set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,

		// the maximum amount of time to wait for the
		// next request when keep-alives are enabled
		IdleTimeout: time.Second * 60,
	}

	// explicitly enable keep-alives
	httpServer.SetKeepAlivesEnabled(true)

	// stdout access log enable/disable
	if viper.GetBool("server.access_log") {
		httpServer.Handler = handlers.CombinedLoggingHandler(os.Stdout, mux)
	} else {
		httpServer.Handler = mux
	}

	return
}
