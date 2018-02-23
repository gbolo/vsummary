package main

import (
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jmoiron/sqlx"
	"github.com/gbolo/vsummary/config"
	back "github.com/gbolo/vsummary/db"
	"github.com/gbolo/vsummary/poller"
	"github.com/gbolo/vsummary/server"
	//"github.com/gbolo/vsummary/crypto"

	//"context"
	//"encoding/json"
	//"fmt"
	//"log"

	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/op/go-logging"
	//"net/http"
	"time"
	//"os"
	//"bytes"
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

	// test
	pollers, err := b.GetPollers()
	fmt.Println(pollers)

	//// vmware section ----------------------------------------------------------------------
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//
	//// connect and login to ESX or vCenter
	//vPoller := poller.NewPoller()
	//vPoller.Config = &poller.PollerConfig{
	//		URL: os.Getenv("VSUMMARY_TEST_URL"),
	//		UserName: os.Getenv("VSUMMARY_TEST_USER"),
	//		Password: os.Getenv("VSUMMARY_TEST_PASSWORD"),
	//		Insecure: true,
	//}
	//
	//err = vPoller.Connect(&ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer vPoller.VmwareClient.Logout(ctx)

	go server.Start()

	time.Sleep(3 * time.Second)

	fmt.Println("---------------------------------------------------------------------------------------")

	poller.LoadPollers(pollers)

	time.Sleep(70 * time.Minute)

}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
