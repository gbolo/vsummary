package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertPortgroup = `
	INSERT INTO portgroup (
		id,
		name,
       	type,
		vlan,
		vlan_type,
		vswitch_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
       	:type,
		:vlan,
		:vlan_type,
		:vswitch_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		type=VALUES(type),
		vlan=VALUES(vlan),
		vlan_type=VALUES(vlan_type),
		vswitch_id=VALUES(vswitch_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertPortgroups(portgroups []common.Portgroup) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE portgroup SET present = 0 WHERE present = 1 AND type=? AND vcenter_id=?", portgroups[0].Type, portgroups[0].VcenterId)
	var rowsAffected int64 = 0

	for _, portgroup := range portgroups {

		// standard vswitch portgroup
		if portgroup.Type == "vSwitch" {
			// vSwitch portgroup doesnt have a moref
			portgroup.Id = common.ComputeId(fmt.Sprintf("%s%s%s", portgroup.VcenterId, portgroup.EsxiMoref, portgroup.Name))
			portgroup.VswitchId = common.ComputeId(fmt.Sprintf("%s%s%s", portgroup.VcenterId, portgroup.EsxiMoref, portgroup.VswitchName))
			portgroup.VlanType = "single"

		} else if portgroup.Type == "DVS" {
			// Fill in some required Ids
			portgroup.Id = common.ComputeId(fmt.Sprintf("%s%s", portgroup.VcenterId, portgroup.Moref))
			portgroup.VswitchId = common.ComputeId(fmt.Sprintf("%s%s", portgroup.VcenterId, portgroup.VswitchMoref))

			switch portgroup.VlanType {
			case "VmwareDistributedVirtualSwitchVlanIdSpec":
				portgroup.VlanType = "single"
			case "VmwareDistributedVirtualSwitchTrunkVlanSpec":
				portgroup.VlanType = "trunk"
			default:
				portgroup.VlanType = "TypeNotImplemented"
			}
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertPortgroup, &portgroup)

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
