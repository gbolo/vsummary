// TODO: this main package is for dev testing only
// real main packages will come later...

package main

import (
	"flag"
	"os"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/config"
	"github.com/gbolo/vsummary/db"
	"github.com/gbolo/vsummary/poller"
	"github.com/gbolo/vsummary/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("vsummary")

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// parse flags
	cfgFile := flag.String("config", "", "path to config file")
	outputVersion := flag.Bool("version", false, "prints version then exits")
	flag.Parse()

	// print version and exit if flag is present
	common.PrintVersion()
	if *outputVersion {
		os.Exit(0)
	}

	// init config and logging
	config.ConfigInit(*cfgFile)

	// init backend
	b, err := db.InitBackend()
	handleErr(err)

	// apply backend schemas
	err = b.ApplySchemas()
	handleErr(err)

	// start vsummary server
	go server.Start()

	// configure and start built-in internalCollector
	poller.BuiltInCollector.SetBackend(*b)
	poller.BuiltInCollector.Run()
}
