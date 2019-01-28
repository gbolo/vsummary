// this main package is for dev testing only
// real main packages will come later...

package main

import (
	"fmt"
	"time"

	"github.com/gbolo/vsummary/common"
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

	common.PrintVersion()
	fmt.Println("---------------------------------------------------------------------------------------")

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

	//b.SelectPoller("500bf0f86671")

	// test external poller
	//pollers, _ := b.GetPollers()
	//e := poller.NewExternalPoller(pollers[0])
	//err = e.SetApiUrl("http://127.0.0.1:8080")
	//if err != nil {
	//	log.Fatalf("error setiing api url",err)
	//}
	//e.Daemonize()
}
