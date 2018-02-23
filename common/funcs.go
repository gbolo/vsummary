package common

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/op/go-logging"
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
