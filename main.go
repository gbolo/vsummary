// this main package is for dev testing only
// real main packages will come later...

package main

import (
	"fmt"
	"time"

	"github.com/gbolo/vsummary/config"
	back "github.com/gbolo/vsummary/db"
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

	// init config and logging
	config.ConfigInit("")

	// init backend
	b, err := back.InitBackend()
	handleErr(err)

	// apply backend schemas
	err = b.ApplySchemas()
	handleErr(err)

	// start vsummary server
	go server.Start()

	time.Sleep(3 * time.Second)
	fmt.Println("---------------------------------------------------------------------------------------")

	// start internalCollector
	i := poller.NewEmptyInternalCollector()
	i.SetBackend(*b)
	i.Run()
}
