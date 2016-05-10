/*

vSummary
MYSQL SCHEMA

Any database schema changes made for this project will be reflected in this file.

TODO:
 - create more views
 - further normalization of the data
 - add foreign keys/constraints
 - add a logging table
 - maybe a last_modified date column?
 - verify correct types are being used for each column

*/

CREATE TABLE vm
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
moref VARCHAR(16),
vmx_path VARCHAR(255),
vcpu SMALLINT UNSIGNED,
memory_mb INT UNSIGNED,
config_guest_os VARCHAR(128),
config_version VARCHAR(16),
smbios_uuid VARCHAR(36),
instance_uuid VARCHAR(36),
config_change_version VARCHAR(64),
guest_tools_version VARCHAR(32),
guest_tools_running VARCHAR(32),
guest_hostname VARCHAR(128),
guest_ip VARCHAR(255),
stat_cpu_usage INT UNSIGNED,
stat_host_memory_usage INT UNSIGNED,
stat_guest_memory_usage INT UNSIGNED,
stat_uptime_sec INT UNSIGNED,
power_state TINYINT UNSIGNED,
esxi_id VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE esxi
(
id VARCHAR(32) PRIMARY KEY,	
name VARCHAR(128),
moref VARCHAR(16),
max_evc VARCHAR(64),
current_evc VARCHAR(64),
status VARCHAR(32),
power_state TINYINT UNSIGNED,
in_maintenance_mode INT,
vendor VARCHAR(64),
model VARCHAR(64),
uuid VARCHAR(36),
memory_bytes BIGINT UNSIGNED,
cpu_model VARCHAR(64),
cpu_mhz INT UNSIGNED,
cpu_sockets SMALLINT UNSIGNED,
cpu_cores SMALLINT UNSIGNED,
cpu_threads SMALLINT UNSIGNED,
nics SMALLINT UNSIGNED,
hbas SMALLINT UNSIGNED,
version VARCHAR(32),
build VARCHAR(32),
stat_cpu_usage INT UNSIGNED,
stat_memory_usage BIGINT UNSIGNED,
stat_uptime_sec INT UNSIGNED,
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE datastore
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
moref VARCHAR(16),
status VARCHAR(32),
capacity_bytes BIGINT UNSIGNED,
free_bytes BIGINT UNSIGNED,
uncommitted_bytes BIGINT UNSIGNED,
type VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE vdisk
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
capacity_bytes BIGINT UNSIGNED,
path VARCHAR(255),
thin_provisioned TINYINT UNSIGNED,
datastore_id VARCHAR(32),
uuid VARCHAR(128),
disk_object_id VARCHAR(32),
vm_id VARCHAR(32),
esxi_id VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE vcenter
(
id VARCHAR(36) PRIMARY KEY NOT NULL,
fqdn VARCHAR(128),
short_name VARCHAR(64),
present TINYINT DEFAULT 1
);


CREATE TABLE pnic
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
mac VARCHAR(17),
link_speed SMALLINT UNSIGNED,
driver VARCHAR(45),
esxi_id VARCHAR(32),
vswitch_id VARCHAR(32) DEFAULT null,
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE vswitch
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
type VARCHAR(64),
version VARCHAR(32) DEFAULT null,
max_mtu SMALLINT UNSIGNED DEFAULT 0,
ports SMALLINT UNSIGNED DEFAULT 0,
esxi_id VARCHAR(32) DEFAULT null,
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE portgroup
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
type VARCHAR(32),
vlan VARCHAR(128),
vlan_type VARCHAR(64),
vswitch_id VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);

CREATE TABLE vnic
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(64),
mac VARCHAR(17),
type VARCHAR(45),
connected VARCHAR(16),
status VARCHAR(16),
vm_id VARCHAR(32),
portgroup_id VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);






CREATE TABLE vmknic
(
id VARCHAR(32) PRIMARY KEY,
name VARCHAR(128),
mac VARCHAR(17),
mtu SMALLINT UNSIGNED,
ip VARCHAR(45),
netmask VARCHAR(32),
portgroup_id VARCHAR(32),
esxi_id VARCHAR(32),
vcenter_id VARCHAR(36),
present TINYINT DEFAULT 1
);






/*

CREATE VIEWS TO SIMPLIFY QUEIRES IN APPLICATION

*/


CREATE VIEW view_vm AS
SELECT  
  vm.*,
  esxi.name AS esxi_name,
  esxi.current_evc AS esxi_current_evc,
  esxi.status AS esxi_status,
  esxi.cpu_model AS esxi_cpu_model,
  coalesce(COUNT(distinct vdisk.id),0) AS vdisks,
  coalesce(COUNT(distinct vnic.id),0) AS vnics,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    vm
LEFT JOIN
        vdisk
ON      vm.id = vdisk.vm_id
    AND vm.present = 1
    AND vdisk.present = 1
LEFT JOIN
        vnic
ON      vm.id = vnic.vm_id
    AND vm.present = 1
    AND vnic.present = 1
LEFT JOIN
        esxi
ON      vm.esxi_id = esxi.id
LEFT JOIN
        vcenter
ON      vm.vcenter_id = vcenter.id
GROUP BY
        vm.id;


CREATE VIEW view_vnic AS
SELECT  
  vnic.*,
  vm.name AS vm_name, 
  esxi.name AS esxi_name, 
  coalesce(portgroup.name,"ORPHANED") AS portgroup_name,
  portgroup.vlan,
  coalesce(vswitch.name,"ORPHANED") AS vswitch_name, 
  vswitch.type AS vswitch_type,
  vswitch.max_mtu AS vswitch_max_mtu,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    vnic
LEFT JOIN
        portgroup
ON      vnic.portgroup_id = portgroup.id
    AND vnic.present = 1
    AND portgroup.present = 1
LEFT JOIN
        vm
ON      vnic.vm_id = vm.id
    AND vm.present = 1
LEFT JOIN
        esxi
ON      vm.esxi_id = esxi.id
    AND esxi.present = 1
LEFT JOIN
        vswitch
ON      portgroup.vswitch_id = vswitch.id
LEFT JOIN
        vcenter
ON      vm.vcenter_id = vcenter.id
GROUP BY
        vnic.id;


CREATE VIEW view_esxi AS
SELECT 
  esxi.*, 
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name,
  ( SELECT coalesce(sum(vm.vcpu),0)
    FROM vm
    WHERE vm.esxi_id = esxi.id AND vm.power_state = 1 AND vm.present = 1) vcpus_powered_on,
  ( SELECT coalesce(sum(vm.memory_mb),0)
    FROM vm
    WHERE vm.esxi_id = esxi.id AND vm.power_state = 1 AND vm.present = 1) vmemory_mb_powered_on,
  ( SELECT coalesce(count(vm.id),0)
    FROM vm 
    WHERE vm.esxi_id = esxi.id AND vm.power_state = 1 AND vm.present = 1) vms_powered_on,
  ( SELECT coalesce(count(pnic.id),0)
    FROM pnic 
    WHERE pnic.esxi_id = esxi.id AND pnic.present = 1) pnics
FROM esxi 
LEFT JOIN
    vcenter
ON  esxi.vcenter_id = vcenter.id;


CREATE VIEW view_datastore AS
SELECT  
  datastore.*,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    datastore
LEFT JOIN
        vcenter
ON      datastore.vcenter_id = vcenter.id
    AND datastore.present = 1
GROUP BY
        datastore.id;


CREATE VIEW view_vdisk AS
SELECT  
  vdisk.*,
  vm.name AS vm_name,
  vm.power_state AS vm_power_state,
  datastore.name AS datastore_name,
  datastore.type AS datastore_type,
  esxi.name AS esxi_name,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    vdisk
LEFT JOIN
        vm
ON      vdisk.vm_id = vm.id
    AND vdisk.present = 1
LEFT JOIN
        datastore
ON      vdisk.datastore_id = datastore.id
    AND vdisk.present = 1
LEFT JOIN
        esxi
ON      vdisk.esxi_id = esxi.id
    AND vdisk.present = 1
LEFT JOIN
        vcenter
ON      vdisk.vcenter_id = vcenter.id
    AND vdisk.present = 1
GROUP BY
        vdisk.id;


CREATE VIEW view_portgroup AS
SELECT  
  portgroup.*,
  vswitch.name AS vswitch_name,
  vswitch.type AS vswitch_type,
  vswitch.max_mtu AS vswitch_max_mtu,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    portgroup
LEFT JOIN
        vswitch
ON      portgroup.vswitch_id = vswitch.id
    AND portgroup.present = 1
LEFT JOIN
        vcenter
ON      portgroup.vcenter_id = vcenter.id
    AND portgroup.present = 1
GROUP BY
        portgroup.id;

/* 
TESTING
*/

CREATE VIEW view_esxi AS
SELECT  
  esxi.*,
  coalesce(COUNT(distinct vm.id),0) AS vms_powered_on,
  coalesce(SUM(vm.vcpu),0) AS vcpus_powered_on, 
  coalesce(SUM(vm.memory_mb),0) AS vmemory_mb_powered_on,
  coalesce(COUNT(distinct pnic.id),0) AS pnics,
  vcenter.fqdn AS vcenter_fqdn,
  vcenter.short_name AS vcenter_short_name
FROM    esxi
LEFT JOIN
        vm
ON      esxi.id = vm.esxi_id
    AND esxi.present = 1
    AND vm.present = 1
    AND vm.power_state = 1
LEFT JOIN
        pnic
ON      esxi.id = pnic.esxi_id
    AND pnic.present = 1
LEFT JOIN
        vcenter
ON      vm.vcenter_id = vcenter.id
GROUP BY
        esxi.id;
