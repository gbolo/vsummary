package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

const insertVm = `
	INSERT INTO vm (
		id,
		name,
		moref,
		vmx_path,
		vcpu,
		memory_mb,
		template,
		config_guest_os,
		config_version,
		smbios_uuid,
		instance_uuid,
		config_change_version,
		guest_tools_version,
		guest_tools_running,
		guest_hostname,
		guest_ip,
		guest_os,
		stat_cpu_usage,
		stat_host_memory_usage,
		stat_guest_memory_usage,
		stat_uptime_sec,
		power_state,
		folder_id,
		vapp_id,
		resourcepool_id,
		esxi_id,
		vcenter_id
		)
	VALUES (
		:id,
		:name,
		:moref,
		:vmx_path,
		:vcpu,
		:memory_mb,
		:template,
		:config_guest_os,
		:config_version,
		:smbios_uuid,
		:instance_uuid,
		:config_change_version,
		:guest_tools_version,
		:guest_tools_running,
		:guest_hostname,
		:guest_ip,
		:guest_os,
		:stat_cpu_usage,
		:stat_host_memory_usage,
		:stat_guest_memory_usage,
		:stat_uptime_sec,
		:power_state,
		:folder_id,
		:vapp_id,
		:resourcepool_id,
		:esxi_id,
		:vcenter_id
		)
	ON DUPLICATE KEY UPDATE
		id=VALUES(id),
		name=VALUES(name),
		moref=VALUES(moref),
		vmx_path=VALUES(vmx_path),
		vcpu=VALUES(vcpu),
		memory_mb=VALUES(memory_mb),
		template=VALUES(template),
		config_guest_os=VALUES(config_guest_os),
		config_version=VALUES(config_version),
		smbios_uuid=VALUES(smbios_uuid),
		instance_uuid=VALUES(instance_uuid),
		config_change_version=VALUES(config_change_version),
		guest_tools_version=VALUES(guest_tools_version),
		guest_tools_running=VALUES(guest_tools_running),
		guest_hostname=VALUES(guest_hostname),
		guest_ip=VALUES(guest_ip),
		guest_os=VALUES(guest_os),
		stat_cpu_usage=VALUES(stat_cpu_usage),
		stat_host_memory_usage=VALUES(stat_host_memory_usage),
		stat_guest_memory_usage=VALUES(stat_guest_memory_usage),
		stat_uptime_sec=VALUES(stat_uptime_sec),
		power_state=VALUES(power_state),
		folder_id=VALUES(folder_id),
		vapp_id=VALUES(vapp_id),
		resourcepool_id=VALUES(resourcepool_id),
		esxi_id=VALUES(esxi_id),
		vcenter_id=VALUES(vcenter_id),
		present=1;`

// InsertVirtualmachines inserts a vm into database
func (b *Backend) InsertVirtualmachines(vms []common.VirtualMachine) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this. better way to detect vcenter.
	tx.MustExec("UPDATE vm SET present = 0 WHERE present = 1 AND vcenter_id=?", vms[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vm := range vms {
		// fill in missing data
		// folder may not exist
		if vm.FolderMoref != "" && vm.FolderMoref != "vapp" {
			vm.FolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.FolderMoref))
		}

		// vapps may not exist
		if vm.VappMoref != "none" {
			vm.VappId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.VappMoref))
		}
		// resourcepool may not exist
		if vm.ResourcePoolMoref != "" {
			vm.ResourcePoolId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.ResourcePoolId))
		}

		// Fill in some required Ids
		vm.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.Moref))
		vm.EsxiId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.EsxiMoref))

		// Store the user record in the DB
		res, err := tx.NamedExec(insertVm, &vm)

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
