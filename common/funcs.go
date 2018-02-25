package common

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"strconv"

	"github.com/op/go-logging"
	"github.com/buger/jsonparser"
	"encoding/json"
)

var log = logging.MustGetLogger("vsummary")

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ExecutionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %s", name, elapsed)
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// When working with dynamic json
// ignore errors, but log them
func GetInt(o interface{}, keys ...string) (i int64){

	b, _ := json.Marshal(o)

	i, err := jsonparser.GetInt(b, keys...)
	if err != nil {
		log.Infof("error parsing json: %s", err)
	}

	return

}

