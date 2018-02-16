package db

const (

	// PREPARED INSERT STATEMENTS --------------------------------------------------------------------------------------

	insertVm = `
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

	insertPoller = `
INSERT INTO poller (
	vcenter_host,
	vcenter_name,
	enabled,
	user_name,
	password,
	interval_min
	)
VALUES (
	:vcenter_host,
	:vcenter_name,
	:enabled,
	:user_name,
	:password,
	:interval_min
	)
ON DUPLICATE KEY UPDATE
	vcenter_name=VALUES(vcenter_name),
	enabled=VALUES(enabled),
	user_name=VALUES(user_name),
	password=VALUES(password),
	interval_min=VALUES(interval_min);`
)
