package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertVDisk = `
	INSERT INTO vdisk (
		id,
		name,
		capacity_bytes,
		path,
		thin_provisioned,
		datastore_id,
		uuid,
		disk_object_id,
		vm_id,
		esxi_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:capacity_bytes,
		:path,
		:thin_provisioned,
		:datastore_id,
		:uuid,
		:disk_object_id,
		:vm_id,
		:esxi_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		capacity_bytes=VALUES(capacity_bytes),
		path=VALUES(path),
		thin_provisioned=VALUES(thin_provisioned),
		datastore_id=VALUES(datastore_id),
		uuid=VALUES(uuid),
		disk_object_id=VALUES(disk_object_id),
		vm_id=VALUES(vm_id),
		esxi_id=VALUES(esxi_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertVDisks(vdisks []common.VDisk) (err error) {

	if len(vdisks) == 0 {
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
	tx.MustExec("UPDATE vdisk SET present = 0 WHERE present = 1 AND vcenter_id=?", vdisks[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vdisk := range vdisks {

		// Fill in some required Ids
		vdisk.Id = common.ComputeId(fmt.Sprintf("%s%s%s", vdisk.VcenterId, vdisk.DiskObjectId, vdisk.Path))
		vdisk.DatastoreId = common.ComputeId(fmt.Sprintf("%s%s", vdisk.VcenterId, vdisk.DatastoreMoref))
		vdisk.VirtualMachineId = common.ComputeId(fmt.Sprintf("%s%s", vdisk.VcenterId, vdisk.VirtualmachineMoref))
		vdisk.EsxiId = common.ComputeId(fmt.Sprintf("%s%s", vdisk.VcenterId, vdisk.EsxiMoref))

		// Determine capacity in bytes if not present
		if vdisk.CapacityBytes == 0 && vdisk.CapacityKb > 0 {
			vdisk.CapacityBytes = vdisk.CapacityKb * 1024
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertVDisk, &vdisk)

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
