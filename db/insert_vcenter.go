package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertVcenter = `
	INSERT INTO vcenter (
		id,
		name,
		host
		)
	VALUES (
		:id,
		:name,
		:host
		)
	ON DUPLICATE KEY UPDATE
		name=VALUES(name),
		host=VALUES(host);`

// InsertVMs inserts a vm into database
func (b *Backend) InsertVcenter(vcenter common.VCenter) (err error) {

	// never allow an empty id for vcenter
	if vcenter.Id == "" {
		log.Errorf("cannot insert vcenter without an id")
		err = fmt.Errorf("vcenter is missing an id")
		return
	}

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	_, err = tx.NamedExec(insertVcenter, &vcenter)

	if err != nil {
		return
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debug("vcenter inserted successfully")

	return

}
