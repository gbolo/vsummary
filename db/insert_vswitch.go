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

const insertDVS = `
	INSERT INTO vswitch (
		id,
		name,
       	type,
		version,
		max_mtu,
		ports,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:type,
		:version,
		:max_mtu,
		:ports,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		type=VALUES(type),
		version=VALUES(version),
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
	tx.MustExec("UPDATE vswitch SET present = 0 WHERE present = 1 AND type=? AND vcenter_id=?", vswitches[0].Type, vswitches[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vswitch := range vswitches {

		insertStatement := ""
		// Fill in some required Ids
		if vswitch.Type == "DVS" {
			insertStatement = insertDVS
			vswitch.Id = common.ComputeId(fmt.Sprintf("%s%s", vswitch.VcenterId, vswitch.Moref))
		} else if vswitch.Type == "vSwitch" {
			insertStatement = insertVswitch
			vswitch.Id = common.ComputeId(fmt.Sprintf("%s%s%s", vswitch.VcenterId, vswitch.EsxiMoref, vswitch.Name))
			vswitch.EsxiId = common.ComputeId(fmt.Sprintf("%s%s", vswitch.VcenterId, vswitch.EsxiMoref))
		} else {
			err = fmt.Errorf("incorrect vswitch type: %s", vswitch.Type)
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertStatement, &vswitch)

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
