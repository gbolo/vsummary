package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertVNics = `
	INSERT INTO vnic (
		id,
		name,
		mac,
       	type,
		connected,
		status,
		vm_id,
		portgroup_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:mac,
       	:type,
		:connected,
		:status,
		:vm_id,
		:portgroup_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		mac=VALUES(mac),
		type=VALUES(type),
		connected=VALUES(connected),
		status=VALUES(status),
		vm_id=VALUES(vm_id),
		portgroup_id=VALUES(portgroup_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertVNics(vnics []common.VNic) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE vnic SET present = 0 WHERE present = 1 AND vcenter_id=?", vnics[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vnic := range vnics {

		// Fill in some required Ids
		vnic.Id = common.GetMD5Hash(fmt.Sprintf("%s%s%s", vnic.VcenterId, vnic.VirtualmachineMoref, vnic.Name))
		vnic.VirtualmachineId = common.GetMD5Hash(fmt.Sprintf("%s%s", vnic.VcenterId, vnic.VirtualmachineMoref))

		// determine unique portgroup id
		if vnic.VswitchType == "HostVirtualSwitch" {
			// standard vswitch
			vnic.PortgroupId = common.GetMD5Hash(fmt.Sprintf("%s%s%s", vnic.VcenterId, vnic.EsxiMoref, vnic.PortgroupName))
		} else if vnic.VswitchType == "VmwareDistributedVirtualSwitch" {
			vnic.PortgroupId = common.GetMD5Hash(fmt.Sprintf("%s%s", vnic.VcenterId, vnic.PortgroupMoref))
		} else {
			vnic.PortgroupId = "ORPHANED"
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertVNics, &vnic)

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
