package db

import (
	"fmt"
	"strings"

	"github.com/gbolo/vsummary/common"
)

const insertResourcepool = `
	INSERT INTO resourcepool (
		id,
		moref,
		name,
		type,
		status,
		vapp_state,
		parent,
		parent_moref,
		cluster_id,
		configured_memory_mb,
		cpu_reservation,
		cpu_limit,
		mem_reservation,
		mem_limit,
		vcenter_id
		)
	VALUES (
		:id,
		:moref,
		:name,
		:type,
		:status,
		:vapp_state,
		:parent,
		:parent_moref,
		:cluster_id,
		:configured_memory_mb,
		:cpu_reservation,
		:cpu_limit,
		:mem_reservation,
		:mem_limit,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		moref=VALUES(moref),
		type=VALUES(type),
		status=VALUES(status),
		vapp_state=VALUES(vapp_state),
		parent=VALUES(parent),
		parent_moref=VALUES(parent_moref),
		cluster_id=VALUES(cluster_id),
		configured_memory_mb=VALUES(configured_memory_mb),
		cpu_limit=VALUES(cpu_limit),
		cpu_reservation=VALUES(cpu_reservation),
		cpu_limit=VALUES(cpu_limit),
		mem_reservation=VALUES(mem_reservation),
		mem_limit=VALUES(mem_limit),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertResourcepools(resourcepools []common.ResourcePool) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE resourcepool SET present = 0 WHERE present = 1 AND vcenter_id=?", resourcepools[0].VcenterId)
	var rowsAffected int64 = 0

	for _, resourcepool := range resourcepools {

		// Fill in some required Ids
		resourcepool.Id = common.ComputeId(fmt.Sprintf("%s%s", resourcepool.VcenterId, resourcepool.Moref))
		resourcepool.ClusterId = common.ComputeId(fmt.Sprintf("%s%s", resourcepool.VcenterId, resourcepool.ClusterMoref))

		// parent information
		if strings.HasPrefix(resourcepool.ParentMoref, "domain-") {
			resourcepool.Parent = "cluster"
		} else {
			resourcepool.Parent = common.ComputeId(fmt.Sprintf("%s%s", resourcepool.VcenterId, resourcepool.ParentMoref))
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertResourcepool, &resourcepool)

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
