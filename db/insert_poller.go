package db

import (
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/crypto"
)

const (
	insertPoller = `
	INSERT INTO poller (
		id,
		vcenter_host,
		vcenter_name,
		enabled,
		user_name,
		password,
		interval_min,
		internal
		)
	VALUES (
		:id,
		:vcenter_host,
		:vcenter_name,
		:enabled,
		:user_name,
		:password,
		:interval_min,
		:internal
		)
	ON DUPLICATE KEY UPDATE
		vcenter_name=VALUES(vcenter_name),
		enabled=VALUES(enabled),
		user_name=VALUES(user_name),
		password=VALUES(password),
		interval_min=VALUES(interval_min),
		internal=VALUES(internal);`

	updatePollDate   = "UPDATE poller SET last_poll=:last_poll WHERE id=:id"
	selectPollerById = "SELECT * FROM poller WHERE id=?"
	deletePollerById = "DELETE FROM poller WHERE id=?"
)

// InsertPoller inserts a poller into database
func (b *Backend) InsertPoller(poller common.Poller) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// attempt to encrypt the provided password before storing to database
	log.Debug("encrypting password before database insert/update")
	encryptedPassword, err := crypto.Encrypt(poller.Password)

	if err != nil {
		return
	}

	poller.Password = encryptedPassword

	// Create an Id
	poller.Id = common.ComputeId(poller.VcenterHost)

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	var rowsAffected int64 = 0

	// Store the user record in the DB
	res, err := tx.NamedExec(insertPoller, &poller)

	if err != nil {
		return
	}

	// tally up rows affected for logging
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debugf("total combined affected rows: %d", rowsAffected)

	return

}

func (b *Backend) UpdateLastPollDate(poller common.Poller) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// Create an Id and date
	poller.Id = common.ComputeId(poller.VcenterHost)
	currentTime := time.Now()
	poller.LastPoll = currentTime.Format("2006-01-02 3:4 pm")

	tx := b.db.MustBegin()
	var rowsAffected int64 = 0

	res, err := tx.NamedExec(updatePollDate, &poller)
	if err != nil {
		return
	}

	// tally up rows affected for logging
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debugf("total combined affected rows: %d", rowsAffected)

	return

}

// SelectPoller returns a single poller
func (b *Backend) SelectPoller(pollerId string) (poller common.Poller, err error) {
	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// do select
	err = b.db.Get(&poller, selectPollerById, pollerId)
	return
}

// RemovePoller removes a specified poller by ID
func (b *Backend) RemovePoller(pollerId string) (err error) {
	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	res, err := b.db.Exec(deletePollerById, pollerId)
	if err == nil {
		rowsAffected, _ := res.RowsAffected()
		log.Debugf("total combined affected rows: %d", rowsAffected)
	}
	return
}
