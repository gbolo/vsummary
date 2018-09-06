package db

import (
	"github.com/gbolo/vsummary/common"
)

// GetPollers returns a list of pollers
func (b *Backend) GetPollers() (pollers []common.Poller, err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// attempt to retrieve pollers
	err = b.db.Select(&pollers, "SELECT * from poller")
	if err != nil {
		return
	}

	return

}
