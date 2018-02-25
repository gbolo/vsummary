package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertCluster = `
	INSERT INTO cluster (
		id,
		name,
		datacenter_id,
		current_balance,
		target_balance,
		total_cpu_threads,
		total_cpu_mhz,
		total_memory_bytes,
		total_vmotions,
		num_hosts,
		drs_enabled,
		drs_behaviour,
		ha_enabled,
		status,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:datacenter_id,
		:current_balance,
		:target_balance,
		:total_cpu_threads,
		:total_cpu_mhz,
		:total_memory_bytes,
		:total_vmotions,
		:num_hosts,
		:drs_enabled,
		:drs_behaviour,
		:ha_enabled,
		:status,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		datacenter_id=VALUES(datacenter_id),
		current_balance=VALUES(current_balance),
		target_balance=VALUES(target_balance),
		total_cpu_threads=VALUES(total_cpu_threads),
		total_cpu_mhz=VALUES(total_cpu_mhz),
		total_memory_bytes=VALUES(total_memory_bytes),
		total_vmotions=VALUES(total_vmotions),
		num_hosts=VALUES(num_hosts),
		drs_enabled=VALUES(drs_enabled),
		drs_behaviour=VALUES(drs_behaviour),
		ha_enabled=VALUES(ha_enabled),
		status=VALUES(status),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertClusters(clusters []common.Cluster) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE cluster SET present = 0 WHERE present = 1 AND vcenter_id=?", clusters[0].VcenterId)
	var rowsAffected int64 = 0

	for _, cluster := range clusters {

		// Fill in some required Ids
		cluster.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", cluster.VcenterId, cluster.Moref))
		cluster.DatacenterId = common.GetMD5Hash(fmt.Sprintf("%s%s", cluster.VcenterId, cluster.DatacenterMoref))

		// Store the record in the DB
		res, err := tx.NamedExec(insertCluster, &cluster)

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
