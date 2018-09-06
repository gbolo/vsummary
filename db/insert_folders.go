package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
	"strings"
)

const insertFolders = `
	INSERT INTO folder (
		id,
		name,
		moref,
       	type,
		parent,
		parent_datacenter_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:moref,
       	:type,
		:parent,
		:parent_datacenter_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		moref=VALUES(moref),
		type=VALUES(type),
		parent=VALUES(parent),
		parent_datacenter_id=VALUES(parent_datacenter_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertFolders(folders []common.Folder) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE folder SET present = 0 WHERE present = 1 AND vcenter_id=?", folders[0].VcenterId)
	var rowsAffected int64 = 0

	for _, folder := range folders {

		// Fill in some required Ids
		folder.Id = common.ComputeId(fmt.Sprintf("%s%s", folder.VcenterId, folder.Moref))

		// type information
		if strings.Contains(folder.Type, "VirtualMachine") {
			folder.Type = "VirtualMachine"
		} else {
			// dont handle non-vm folders for now
			folder.Type = "not_vm"
		}

		// parent information
		if strings.HasPrefix(folder.ParentMoref, "datacenter-") {
			folder.Parent = "datacenter"
			folder.ParentDatacenterId = common.ComputeId(fmt.Sprintf("%s%s", folder.VcenterId, folder.ParentMoref))
		} else {
			folder.Parent = common.ComputeId(fmt.Sprintf("%s%s", folder.VcenterId, folder.ParentMoref))
			folder.ParentDatacenterId = "n/a"
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertFolders, &folder)

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

	// TODO: we need to update folder full path after after we instered all the folders above
	// call some function to do that here

	return

}
