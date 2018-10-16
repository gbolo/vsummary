package common

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// returns a human readable string from any number of bytes
// example: 1855425871872 will return 1.9 TB
func BytesHumanReadable(bytes string) string {
	if bytes == "" {
		return "0"
	}
	// ignore numbers after a possible decimal
	bytesSplit := strings.Split(bytes, ".")
	b, err := strconv.ParseInt(bytesSplit[0], 10, 64)
	if err != nil {
		log.Errorf("parse int err: %s", err)
		return "000"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// returns a human readable string from any number of megabytes
func MegaBytesHumanReadable(megaBytes string) string {
	if megaBytes == "" {
		return "0"
	}
	// ignore numbers after a possible decimal
	megaBytesSplit := strings.Split(megaBytes, ".")
	b, _ := strconv.ParseInt(megaBytesSplit[0], 10, 64)
	return BytesHumanReadable(fmt.Sprintf("%d", (b * 1000 * 1000)))
}

// converts seconds to days
func SecondsToHuman(secondsString string) string {
	seconds, _ := strconv.ParseInt(secondsString, 10, 64)
	days := math.Floor(float64(seconds) / 86400)
	hours := math.Floor(float64(seconds%86400) / 3600)
	minutes := math.Floor(float64(seconds%86400%3600) / 60)

	if seconds == 0 {
		return "nil"
	} else if days < 1 {
		return fmt.Sprintf("%dh, %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%v days", days)
	}
}
