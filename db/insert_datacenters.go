package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertDatacenter = `
	INSERT INTO datacenter (
		id,
		name,
		vm_folder_id,
		esxi_folder_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:vm_folder_id,
		:esxi_folder_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		vm_folder_id=VALUES(vm_folder_id),
		esxi_folder_id=VALUES(esxi_folder_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// InsertVMs inserts a vm into database
func (b *Backend) InsertDatacenters(dcs []common.Datacenter) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE datacenter SET present = 0 WHERE present = 1 AND vcenter_id=?", dcs[0].VcenterId)
	var rowsAffected int64 = 0

	for _, dc := range dcs {

		// Fill in some required Ids
		dc.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", dc.VcenterId, dc.Moref))
		dc.EsxiFolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", dc.VcenterId, dc.EsxiFolderMoref))
		dc.VmFolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", dc.VcenterId, dc.VmFolderMoref))

		// Store the record in the DB
		res, err := tx.NamedExec(insertDatacenter, &dc)

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
