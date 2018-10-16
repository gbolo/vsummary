package common

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("vsummary")
)

// ComputeId returns the first 12 characters from a SHA1 hash of the input text
func ComputeId(input string) string {
	sum := fmt.Sprintf("%x", sha1.Sum([]byte(input)))
	return sum[0:12]
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

// set default value if empty
func SetDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return fmt.Sprintf("%s", value)
}