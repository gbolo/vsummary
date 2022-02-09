package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertDatastore = `
	INSERT INTO datastore (
		id,
		name,
		moref,
		status,
		capacity_bytes,
		free_bytes,
		uncommitted_bytes,
		type,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:moref,
		:status,
		:capacity_bytes,
		:free_bytes,
		:uncommitted_bytes,
		:type,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		moref=VALUES(moref),
		status=VALUES(status),
		capacity_bytes=VALUES(capacity_bytes),
		free_bytes=VALUES(free_bytes),
		uncommitted_bytes=VALUES(uncommitted_bytes),
		type=VALUES(type),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

func (b *Backend) InsertDatastores(dss []common.Datastore) (err error) {

	if len(dss) == 0 {
		return
	}

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE datastore SET present = 0 WHERE present = 1 AND vcenter_id=?", dss[0].VcenterId)
	var rowsAffected int64 = 0

	for _, ds := range dss {

		// Fill in some required Ids
		ds.Id = common.ComputeId(fmt.Sprintf("%s%s", ds.VcenterId, ds.Moref))

		// Store the record in the DB
		res, err := tx.NamedExec(insertDatastore, &ds)

		if err != nil {
			break
		}

		// tally up rows affected for logging
		numRowsAffected, err := res.RowsAffected()
		if err != nil {
			break
		}
		rowsAffected = rowsAffected + numRowsAffected

	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debugf("total combined affected rows: %d", rowsAffected)

	return

}
