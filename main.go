package main

import (

	_ "github.com/go-sql-driver/mysql"
	//"github.com/jmoiron/sqlx"
	back "github.com/gbolo/vsummary/db"
	"github.com/gbolo/vsummary/poller"
	"github.com/gbolo/vsummary/server"
	"github.com/gbolo/vsummary/config"
	"github.com/gbolo/vsummary/crypto"

	"context"
	//"encoding/json"
	//"fmt"
	//"log"

	"crypto/md5"
	"encoding/hex"
	"github.com/op/go-logging"
	"fmt"
	//"net/http"
	"time"
	"os"
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



	// vmware section ----------------------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	// connect and login to ESX or vCenter
	vPoller := poller.NewPoller()
	vPoller.Config = &poller.PollerConfig{
			URL: os.Getenv("VSUMMARY_TEST_URL"),
			UserName: os.Getenv("VSUMMARY_TEST_USER"),
			Password: os.Getenv("VSUMMARY_TEST_PASSWORD"),
			Insecure: true,
	}

	err = vPoller.Connect(&ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer vPoller.VmwareClient.Logout(ctx)

	//// get list of VMs
	//vmList, err := vPoller.GetVMs()
	//
	//// print list of vms
	//jsonVms, err := json.Marshal(vmList[0])
	//if err == nil {
	//	fmt.Println(string(jsonVms))
	//} else {
	//	fmt.Println("Error:", err)
	//}

	go server.Start()

	enc, err := crypto.Encrypt("some sample data to encrypt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Encrypted:", enc)

	decr, err := crypto.Decrypt("c_Ijfq9iqD_8acr2AQuxbVj9GXR17ZkJc8u8Gyn7LG84aTZZ3fGZL6tEfA==")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Decrypted:", decr)

	time.Sleep(3 * time.Second)

	fmt.Println("---------------------------------------------------------------------------------------")

	vPoller.Loop()

}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}