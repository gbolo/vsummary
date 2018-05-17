package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertVswitch = `
	INSERT INTO vswitch (
		id,
		name,
       	type,
		esxi_id,
		max_mtu,
		ports,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:type,
		:esxi_id,
		:max_mtu,
		:ports,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		type=VALUES(type),
		esxi_id=VALUES(esxi_id),
		max_mtu=VALUES(max_mtu),
		ports=VALUES(ports),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertVSwitch(vswitches []common.VSwitch) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE vswitch SET present = 0 WHERE present = 1 AND type='SVS' AND vcenter_id=?", vswitches[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vswitch := range vswitches {

		// Fill in some required Ids
		vswitch.Id = common.GetMD5Hash(fmt.Sprintf("%s%s%s", vswitch.VcenterId, vswitch.EsxiMoref, vswitch.Name))
		vswitch.EsxiId = common.GetMD5Hash(fmt.Sprintf("%s%s", vswitch.VcenterId, vswitch.EsxiMoref))

		// Store the record in the DB
		res, err := tx.NamedExec(insertVswitch, &vswitch)

		if err != nil {
			log.Errorf("error storing record: %s", err)
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
