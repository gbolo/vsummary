package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("vsummary")

// return md5 hash of string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// logs a debug message indicating how long something took to execute
func ExecutionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %s", name, elapsed)
}

// converts a boolean to human readable string
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// When working with dynamic json
// ignore errors, but log them
func GetInt(o interface{}, keys ...string) (i int64) {

	b, _ := json.Marshal(o)

	i, err := jsonparser.GetInt(b, keys...)
	if err != nil {
		log.Infof("error parsing json: %s", err)
	}

	return

}

// json parser - returns string value of key
func GetString(o interface{}, keys ...string) (s string) {

	b, _ := json.Marshal(o)

	s, err := jsonparser.GetString(b, keys...)
	if err != nil {
		log.Infof("error parsing json: %s", err)
	}

	return
}

// json parser - returns boolean value of key
func GetBool(o interface{}, keys ...string) (l bool) {

	b, _ := json.Marshal(o)

	l, err := jsonparser.GetBoolean(b, keys...)
	if err != nil {
		log.Infof("error parsing json: %s", err)
	}

	return
}

// json parser - returns true if key exists
func CheckIfKeyExists(o interface{}, keys ...string) (e bool) {

	b, _ := json.Marshal(o)

	_, dataType, _, err := jsonparser.Get(b, keys...)

	if dataType != jsonparser.NotExist || err == nil {
		e = true
	}

	return
}
