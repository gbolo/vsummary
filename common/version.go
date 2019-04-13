package common

import (
	"fmt"

	"runtime"
)

var (
	// !! These variables defined by the Makefile and passed in with ldflags !!
	// !! DO NOT CHANGE THESE DEFAULT VALUES !!

	// Version of application
	Version = "devel"
	// CommitSHA is the short SHA hash of the git commit
	CommitSHA = "unknown"
	// BuildDate is the date this application was compiled
	BuildDate = "unknown"
)

// PrintVersion prints the current version information to stdout
func PrintVersion() {
	fmt.Printf(`------- version info -------
  version     : %s
  build date  : %s
  git hash    : %s
  go version  : %s
  go compiler : %s
  platform    : %s/%s
`, Version, BuildDate, CommitSHA, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
