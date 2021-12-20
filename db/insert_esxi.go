package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertEsxi = `
	INSERT INTO esxi (
		id,
		name,
		cluster_id,
		max_evc,
		current_evc,
		power_state,
		vendor,
		model,
		uuid,
		memory_bytes,
		cpu_model,
		cpu_mhz,
		cpu_sockets,
		cpu_cores,
		cpu_threads,
		nics,
		hbas,
		version,
		build,
		stat_cpu_usage,
		stat_memory_usage,
		stat_uptime_sec,
		status,
		in_maintenance_mode,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:cluster_id,
		:max_evc,
		:current_evc,
		:power_state,
		:vendor,
		:model,
		:uuid,
		:memory_bytes,
		:cpu_model,
		:cpu_mhz,
		:cpu_sockets,
		:cpu_cores,
		:cpu_threads,
		:nics,
		:hbas,
		:version,
		:build,
		:stat_cpu_usage,
		:stat_memory_usage,
		:stat_uptime_sec,
		:status,
		:in_maintenance_mode,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		cluster_id=VALUES(cluster_id),
		max_evc=VALUES(max_evc),
		current_evc=VALUES(current_evc),
		power_state=VALUES(power_state),
		vendor=VALUES(vendor),
		model=VALUES(model),
		uuid=VALUES(uuid),
		memory_bytes=VALUES(memory_bytes),
		cpu_model=VALUES(cpu_model),
		cpu_mhz=VALUES(cpu_mhz),
		cpu_sockets=VALUES(cpu_sockets),
		cpu_cores=VALUES(cpu_cores),
		cpu_threads=VALUES(cpu_threads),
		nics=VALUES(nics),
		hbas=VALUES(hbas),
		version=VALUES(version),
		build=VALUES(build),
		stat_cpu_usage=VALUES(stat_cpu_usage),
		stat_memory_usage=VALUES(stat_memory_usage),
		stat_uptime_sec=VALUES(stat_uptime_sec),
		status=VALUES(status),
		in_maintenance_mode=VALUES(in_maintenance_mode),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// Insert into database
func (b *Backend) InsertEsxi(esxis []common.Esxi) (err error) {

	if len(esxis) == 0 {
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
	tx.MustExec("UPDATE esxi SET present = 0 WHERE present = 1 AND vcenter_id=?", esxis[0].VcenterId)
	var rowsAffected int64 = 0

	for _, esxi := range esxis {

		// Fill in some required Ids
		esxi.Id = common.ComputeId(fmt.Sprintf("%s%s", esxi.VcenterId, esxi.Moref))
		esxi.ClusterId = common.ComputeId(fmt.Sprintf("%s%s", esxi.VcenterId, esxi.ClusterMoref))

		if esxi.CurrentEvc == "" {
			esxi.CurrentEvc = "NULL"
		}

		// Store the record in the DB
		res, err := tx.NamedExec(insertEsxi, &esxi)

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
