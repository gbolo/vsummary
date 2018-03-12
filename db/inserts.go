package db

import (
	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/crypto"
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
