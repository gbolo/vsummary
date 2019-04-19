package db

type SqlSchema struct {
	Name     string
	SqlQuery string
}

func generateSqlSchemas() (schemas []SqlSchema) {

	schemas = append(
		schemas,

		// tables
		SqlSchema{"VirtualMachine", schemaVm},
		SqlSchema{"Datacenter", schemaDatacenter},
		SqlSchema{"Poller", schemaPoller},
		SqlSchema{"Cluster", schemaCluster},
		SqlSchema{"Esxi", schemaEsxi},
		SqlSchema{"vCenter", schemaVcenter},
		SqlSchema{"Cluster", schemaResourcepool},
		SqlSchema{"Folder", schemaFolder},
		SqlSchema{"vSwitch", schemaVswitch},
		SqlSchema{"vDisk", schemaVdisk},
		SqlSchema{"vNic", schemaVnic},
		SqlSchema{"pNic", schemaPnic},
		SqlSchema{"vmkNic", schemaVmknic},
		SqlSchema{"Portgroup", schmeaPortgroup},
		SqlSchema{"Datastore", schemaDatastore},

		// views
		SqlSchema{"ViewVirtualMachine", schemaViewVm},
		SqlSchema{"ViewCluster", schemaViewCluster},
		SqlSchema{"ViewClusterCapacity", schemaViewClusterCapacity},
		SqlSchema{"ViewDatastore", schemaViewDatastore},
		SqlSchema{"ViewEsxi", schemaViewEsxi},
		SqlSchema{"ViewPortgroup", schemaViewPortgroup},
		SqlSchema{"ViewVDisk", schemaViewVdisk},
		SqlSchema{"ViewVNic", schemaViewVnic},
		SqlSchema{"ViewVCenter", schemaViewVcenter},
	)

	return
}

// defined table schemas -----------------------------------------------------------------------------------------------
const (
	schemaVm = `
CREATE TABLE IF NOT EXISTS vm
  (
     id                      VARCHAR(32) PRIMARY KEY,
     name                    VARCHAR(128),
     moref                   VARCHAR(128),
     vmx_path                VARCHAR(255),
     vcpu                    SMALLINT UNSIGNED,
     memory_mb               INT UNSIGNED,
     config_guest_os         VARCHAR(128),
     config_version          VARCHAR(16),
     smbios_uuid             VARCHAR(36),
     instance_uuid           VARCHAR(36),
     config_change_version   VARCHAR(64),
     template                VARCHAR(16),
     guest_tools_version     VARCHAR(32),
     guest_tools_running     VARCHAR(32),
     guest_hostname          VARCHAR(128),
     guest_ip                VARCHAR(255),
     guest_os                VARCHAR(128),
     stat_cpu_usage          INT UNSIGNED,
     stat_host_memory_usage  INT UNSIGNED,
     stat_guest_memory_usage INT UNSIGNED,
     stat_uptime_sec         INT UNSIGNED,
     power_state             VARCHAR(16),
     folder_id               VARCHAR(32),
     vapp_id                 VARCHAR(32),
     resourcepool_id         VARCHAR(32),
     esxi_id                 VARCHAR(32),
     vcenter_id              VARCHAR(36),
     present                 TINYINT DEFAULT 1
  );`

	schemaDatacenter = `
CREATE TABLE IF NOT EXISTS datacenter
  (
     id             VARCHAR(32) PRIMARY KEY,
     vm_folder_id   VARCHAR(32),
     esxi_folder_id VARCHAR(32),
     name           VARCHAR(128),
     vcenter_id     VARCHAR(36),
     present        TINYINT DEFAULT 1
  );`

	schemaPoller = `
CREATE TABLE IF NOT EXISTS poller
  (
     id             VARCHAR(12) PRIMARY KEY,
     vcenter_host   VARCHAR(64),
     vcenter_name   VARCHAR(64),
     enabled        TINYINT DEFAULT 1,
     user_name      VARCHAR(128),
     password       VARCHAR(256),
     interval_min   INT UNSIGNED,
     internal       TINYINT DEFAULT 0,
     last_poll      VARCHAR(128) DEFAULT 'never'
  );`

	schemaVcenter = `
CREATE TABLE IF NOT EXISTS vcenter
  (
     id     VARCHAR(36) PRIMARY KEY,
     name   VARCHAR(64),
     host   VARCHAR(64)
  );`

	schemaCluster = `
CREATE TABLE IF NOT EXISTS cluster
  (
     id                 VARCHAR(32) PRIMARY KEY,
     name               VARCHAR(128),
     datacenter_id      VARCHAR(32),
     total_cpu_threads  INT UNSIGNED,
     total_cpu_mhz      BIGINT UNSIGNED,
     total_memory_bytes BIGINT UNSIGNED,
     total_vmotions     INT UNSIGNED,
     num_hosts          SMALLINT UNSIGNED,
     drs_enabled        VARCHAR(16),
     drs_behaviour      VARCHAR(64),
     ha_enabled         VARCHAR(16),
     current_balance    INT,
     target_balance     INT,
     status             VARCHAR(36),
     vcenter_id         VARCHAR(36),
     present            TINYINT DEFAULT 1
  );`

	schemaEsxi = `
CREATE TABLE IF NOT EXISTS esxi
  (
     id                  VARCHAR(32) PRIMARY KEY,
     name                VARCHAR(128),
     cluster_id          VARCHAR(32),
     max_evc             VARCHAR(64),
     current_evc         VARCHAR(64),
     power_state         VARCHAR(16),
     in_maintenance_mode VARCHAR(16),
     vendor              VARCHAR(64),
     model               VARCHAR(64),
     uuid                VARCHAR(36),
     memory_bytes        BIGINT UNSIGNED,
     cpu_model           VARCHAR(64),
     cpu_mhz             INT UNSIGNED,
     cpu_sockets         SMALLINT UNSIGNED,
     cpu_cores           SMALLINT UNSIGNED,
     cpu_threads         SMALLINT UNSIGNED,
     nics                SMALLINT UNSIGNED,
     hbas                SMALLINT UNSIGNED,
     version             VARCHAR(32),
     build               VARCHAR(32),
     stat_cpu_usage      INT UNSIGNED,
     stat_memory_usage   BIGINT UNSIGNED,
     stat_uptime_sec     INT UNSIGNED,
     status              VARCHAR(36),
     vcenter_id          VARCHAR(36),
     present             TINYINT DEFAULT 1
  );`

	schemaResourcepool = `
CREATE TABLE IF NOT EXISTS resourcepool
  (
     id                   VARCHAR(32) PRIMARY KEY,
     moref                VARCHAR(128),
     full_path            VARCHAR(512),
     name                 VARCHAR(128),
     type                 VARCHAR(64),
     status               VARCHAR(64),
     vapp_state           VARCHAR(64),
     vapp_in_path         TINYINT DEFAULT 0,
     configured_memory_mb BIGINT UNSIGNED,
     cpu_reservation      BIGINT UNSIGNED,
     cpu_limit            BIGINT,
     mem_reservation      BIGINT UNSIGNED,
     mem_limit            BIGINT,
     parent               VARCHAR(32),
     parent_moref         VARCHAR(128),
     cluster_id           VARCHAR(32),
     vcenter_id           VARCHAR(36),
     present              TINYINT DEFAULT 1
  );`

	schemaFolder = `
CREATE TABLE IF NOT EXISTS folder
  (
     id                   VARCHAR(32) PRIMARY KEY,
     moref                VARCHAR(128),
     name                 VARCHAR(128),
     type                 VARCHAR(64),
     full_path            VARCHAR(512),
     parent               VARCHAR(32),
     parent_datacenter_id VARCHAR(32),
     vcenter_id           VARCHAR(36),
     present              TINYINT DEFAULT 1
  );`

	schemaVdisk = `
CREATE TABLE IF NOT EXISTS vdisk
  (
     id               VARCHAR(32) PRIMARY KEY,
     name             VARCHAR(128),
     capacity_bytes   BIGINT UNSIGNED,
     path             VARCHAR(255),
     thin_provisioned VARCHAR(16),
     datastore_id     VARCHAR(32),
     uuid             VARCHAR(128),
     disk_object_id   VARCHAR(32),
     vm_id            VARCHAR(32),
     esxi_id          VARCHAR(32),
     vcenter_id       VARCHAR(36),
     present          TINYINT DEFAULT 1,
     INDEX vmid_ix (vm_id)
  );`

	schemaPnic = `
CREATE TABLE IF NOT EXISTS pnic
  (
     id         VARCHAR(32) PRIMARY KEY,
     name       VARCHAR(128),
     mac        VARCHAR(17),
     link_speed SMALLINT UNSIGNED,
     driver     VARCHAR(45),
     esxi_id    VARCHAR(32),
     vswitch_id VARCHAR(32) DEFAULT NULL,
     vcenter_id VARCHAR(36),
     present    TINYINT DEFAULT 1
  );`

	schemaVnic = `
CREATE TABLE IF NOT EXISTS vnic
  (
     id           VARCHAR(32) PRIMARY KEY,
     name         VARCHAR(64),
     mac          VARCHAR(17),
     type         VARCHAR(45),
     connected    VARCHAR(16),
     status       VARCHAR(16),
     vm_id        VARCHAR(32),
     portgroup_id VARCHAR(32),
     vcenter_id   VARCHAR(36),
     present      TINYINT DEFAULT 1,
     INDEX vmid_ix (vm_id)
  );`

	schemaVswitch = `
CREATE TABLE IF NOT EXISTS vswitch
  (
     id         VARCHAR(32) PRIMARY KEY,
     name       VARCHAR(128),
     type       VARCHAR(64),
     version    VARCHAR(32) DEFAULT NULL,
     max_mtu    SMALLINT UNSIGNED DEFAULT 0,
     ports      SMALLINT UNSIGNED DEFAULT 0,
     esxi_id    VARCHAR(32) DEFAULT NULL,
     vcenter_id VARCHAR(36),
     present    TINYINT DEFAULT 1
  );`

	schmeaPortgroup = `
CREATE TABLE IF NOT EXISTS portgroup
  (
     id         VARCHAR(32) PRIMARY KEY,
     name       VARCHAR(128),
     type       VARCHAR(32),
     vlan       VARCHAR(128),
     vlan_type  VARCHAR(64),
     vswitch_id VARCHAR(32),
     vcenter_id VARCHAR(36),
     present    TINYINT DEFAULT 1
  );`

	schemaVmknic = `
CREATE TABLE IF NOT EXISTS vmknic
  (
     id           VARCHAR(32) PRIMARY KEY,
     name         VARCHAR(128),
     mac          VARCHAR(17),
     mtu          SMALLINT UNSIGNED,
     ip           VARCHAR(45),
     netmask      VARCHAR(32),
     portgroup_id VARCHAR(32),
     esxi_id      VARCHAR(32),
     vcenter_id   VARCHAR(36),
     present      TINYINT DEFAULT 1
  );`

	schemaDatastore = `
CREATE TABLE IF NOT EXISTS datastore
  (
     id                VARCHAR(32) PRIMARY KEY,
     name              VARCHAR(128),
     moref             VARCHAR(128),
     status            VARCHAR(32),
     capacity_bytes    BIGINT UNSIGNED,
     free_bytes        BIGINT UNSIGNED,
     uncommitted_bytes BIGINT UNSIGNED,
     type              VARCHAR(32),
     vcenter_id        VARCHAR(36),
     present           TINYINT DEFAULT 1
  );`
)

// defined table schemas -----------------------------------------------------------------------------------------------
const (
	schemaViewVm = `
CREATE OR REPLACE VIEW view_vm AS
SELECT
  vm.*,
  (SELECT full_path FROM folder WHERE id = vm.folder_id) AS folder,
  e.name AS esxi_name,
  e.current_evc AS esxi_current_evc,
  e.status AS esxi_status,
  e.cpu_model AS esxi_cpu_model,
  e.cluster_id AS cluster_id,
  (SELECT COUNT(1) FROM vdisk WHERE vm_id = vm.id and present = 1) AS vdisks,
  (SELECT COUNT(1) FROM vnic WHERE vm_id = vm.id and present = 1) AS vnics,
  COALESCE((SELECT name FROM cluster WHERE id = e.cluster_id), 'n/a') AS cluster,
  COALESCE((SELECT full_path FROM resourcepool WHERE id = vm.resourcepool_id), 'n/a') AS pool,
  (SELECT name FROM datacenter WHERE esxi_folder_id = (SELECT datacenter_id FROM cluster WHERE id = e.cluster_id)) AS datacenter,
  vc.host AS vcenter_fqdn,
  vc.name AS vcenter_short_name
FROM
  vm,
  esxi e,
  vcenter vc
WHERE
  e.id = vm.esxi_id AND
  vc.id = vm.vcenter_id AND
  vm.present = 1;`

	schemaViewVnic = `
CREATE OR REPLACE VIEW view_vnic AS
SELECT
  vnic.*,
  vm.name AS vm_name,
  esxi.name AS esxi_name,
  coalesce(portgroup.name,"ORPHANED") AS portgroup_name,
  portgroup.vlan,
  coalesce(vswitch.name,"ORPHANED") AS vswitch_name,
  vswitch.type AS vswitch_type,
  vswitch.max_mtu AS vswitch_max_mtu,
  vcenter.host AS vcenter_fqdn,
  vcenter.name AS vcenter_short_name
FROM    vnic
LEFT JOIN
        portgroup
ON      vnic.portgroup_id = portgroup.id
LEFT JOIN
        vm
ON      vnic.vm_id = vm.id
LEFT JOIN
        esxi
ON      vm.esxi_id = esxi.id
LEFT JOIN
        vswitch
ON      portgroup.vswitch_id = vswitch.id
LEFT JOIN
        vcenter
ON      vm.vcenter_id = vcenter.id
WHERE   vnic.present = 1;`

	schemaViewEsxi = `
CREATE OR REPLACE VIEW view_esxi AS
SELECT
  esxi.*,
  vcenter.host AS vcenter_fqdn,
  vcenter.name AS vcenter_short_name,
  coalesce(cluster.name,'n/a') AS cluster,
  datacenter.name AS datacenter,
  ( SELECT coalesce(sum(vm.vcpu),0)
    FROM vm
    WHERE vm.esxi_id = esxi.id AND vm.power_state = "poweredOn" AND vm.present = 1) vcpus_powered_on,
  ( SELECT coalesce(sum(vm.memory_mb),0)
    FROM vm
    WHERE vm.esxi_id = esxi.id AND vm.power_state = "poweredOn" AND vm.present = 1) vmemory_mb_powered_on,
  ( SELECT coalesce(count(vm.id),0)
    FROM vm
    WHERE vm.esxi_id = esxi.id AND vm.power_state = "poweredOn" AND vm.present = 1) vms_powered_on,
  ( SELECT coalesce(count(pnic.id),0)
    FROM pnic
    WHERE pnic.esxi_id = esxi.id AND pnic.present = 1) pnics
FROM esxi
LEFT JOIN
        cluster
ON      esxi.cluster_id = cluster.id
LEFT JOIN
        datacenter
ON      cluster.datacenter_id = datacenter.esxi_folder_id
LEFT JOIN
    vcenter
ON  esxi.vcenter_id = vcenter.id
WHERE esxi.present = 1;`

	schemaViewDatastore = `
CREATE OR REPLACE VIEW view_datastore AS
SELECT
  datastore.*,
  vcenter.host AS vcenter_fqdn,
  vcenter.name AS vcenter_short_name
FROM    datastore
LEFT JOIN
        vcenter
ON      datastore.vcenter_id = vcenter.id
WHERE   datastore.present = 1
GROUP BY
        datastore.id;`

	schemaViewVdisk = `
CREATE OR REPLACE VIEW view_vdisk AS
SELECT
  vdisk.*,
  vm.name AS vm_name,
  vm.power_state AS vm_power_state,
  datastore.name AS datastore_name,
  datastore.type AS datastore_type,
  esxi.name AS esxi_name,
  vcenter.host AS vcenter_fqdn,
  vcenter.name AS vcenter_short_name
FROM    vdisk
LEFT JOIN
        vm
ON      vdisk.vm_id = vm.id
LEFT JOIN
        datastore
ON      vdisk.datastore_id = datastore.id
LEFT JOIN
        esxi
ON      vdisk.esxi_id = esxi.id
LEFT JOIN
        vcenter
ON      vdisk.vcenter_id = vcenter.id
WHERE   vdisk.present = 1
GROUP BY
        vdisk.id;`

	schemaViewPortgroup = `
CREATE OR REPLACE VIEW view_portgroup AS
SELECT DISTINCT
  portgroup.name,
  portgroup.type,
  portgroup.vlan,
  portgroup.vlan_type,
  vswitch.name AS vswitch_name,
  vswitch.type AS vswitch_type,
  vswitch.max_mtu AS vswitch_max_mtu,
  vcenter.host AS vcenter_fqdn,
  vcenter.name AS vcenter_short_name,
  ( SELECT coalesce(count(vnic.id),0)
    FROM vnic
    WHERE vnic.portgroup_id = portgroup.id AND vnic.present = 1) vnics
FROM    portgroup
LEFT JOIN
        vswitch
ON      portgroup.vswitch_id = vswitch.id
LEFT JOIN
        vcenter
ON      portgroup.vcenter_id = vcenter.id
WHERE   portgroup.present = 1
GROUP BY
        portgroup.id;`

	schemaViewVcenter = `
CREATE OR REPLACE VIEW view_vcenter AS
SELECT
  vcenter.*,
  ( SELECT coalesce(sum(vm.vcpu),0)
    FROM vm
    WHERE vm.vcenter_id = vcenter.id AND vm.power_state = "poweredOn" AND vm.present = 1) vms_vcpu_on,
  ( SELECT coalesce(sum(vm.memory_mb),0)
    FROM vm
    WHERE vm.vcenter_id = vcenter.id AND vm.power_state = "poweredOn" AND vm.present = 1) vms_memory_on,
  ( SELECT coalesce(count(vm.id),0)
    FROM vm
    WHERE vm.vcenter_id = vcenter.id AND vm.power_state = "poweredOn" AND vm.present = 1) vms_on,
  ( SELECT coalesce(count(vm.id),0)
    FROM vm
    WHERE vm.vcenter_id = vcenter.id AND vm.present = 1) vms,
  ( SELECT coalesce(count(datacenter.id),0)
    FROM datacenter
    WHERE datacenter.vcenter_id = vcenter.id AND datacenter.present = 1) datacenters,
  ( SELECT coalesce(count(cluster.id),0)
    FROM cluster
    WHERE cluster.vcenter_id = vcenter.id AND cluster.present = 1) clusters,
  ( SELECT coalesce(count(esxi.id),0)
    FROM esxi
    WHERE esxi.vcenter_id = vcenter.id AND esxi.present = 1) esxi_hosts,
  ( SELECT coalesce(sum(esxi.cpu_threads),0)
    FROM esxi
    WHERE esxi.vcenter_id = vcenter.id AND esxi.present = 1) esxi_cpu,
  ( SELECT coalesce(sum(esxi.memory_bytes),0)
    FROM esxi
    WHERE esxi.vcenter_id = vcenter.id AND esxi.present = 1) esxi_memory,
  ( SELECT coalesce(count(vnic.id),0)
    FROM vnic
    WHERE vnic.vcenter_id = vcenter.id AND vnic.present = 1) vnics,
  ( SELECT coalesce(count(vdisk.id),0)
    FROM vdisk
    WHERE vdisk.vcenter_id = vcenter.id AND vdisk.present = 1) vdisks,
  ( SELECT coalesce(count(datastore.id),0)
    FROM datastore
    WHERE datastore.vcenter_id = vcenter.id AND datastore.present = 1) datastores,
  ( SELECT coalesce(count(portgroup.id),0)
    FROM portgroup
    WHERE portgroup.vcenter_id = vcenter.id AND portgroup.present = 1) portgroups,
  ( SELECT coalesce(count(vswitch.id),0)
    FROM vswitch
    WHERE vswitch.vcenter_id = vcenter.id AND vswitch.present = 1) vswitches,
  ( SELECT coalesce(count(resourcepool.id),0)
    FROM resourcepool
    WHERE resourcepool.vcenter_id = vcenter.id AND resourcepool.present = 1) resourcepools
FROM vcenter;`

	schemaViewCluster = `
CREATE OR REPLACE VIEW view_cluster AS
SELECT
  cluster.*, vcenter.name AS vcenter_short_name,
  ( cluster.total_memory_bytes / cluster.num_hosts ) AS avg_memory_per_host,
  ( SELECT coalesce(sum(esxi.stat_memory_usage),0)*1024*1024
    FROM esxi
    WHERE esxi.cluster_id = cluster.id AND esxi.present = 1 AND esxi.power_state = "poweredOn" ) AS total_memory_used,
  ( SELECT coalesce(count(view_vm.id),0)
    FROM view_vm
    WHERE view_vm.cluster_id = cluster.id AND view_vm.power_state = "poweredOn" AND view_vm.present = 1) AS vms_on,
  ( SELECT coalesce(sum(view_vm.vcpu),0)
    FROM view_vm
    WHERE view_vm.cluster_id = cluster.id AND view_vm.power_state = "poweredOn" AND view_vm.present = 1) total_vms_vcpu_on
FROM cluster
LEFT JOIN
    vcenter
ON  cluster.vcenter_id = vcenter.id
WHERE cluster.present = 1;`

	schemaViewClusterCapacity = `CREATE OR REPLACE VIEW view_cluster_capacity AS
SELECT view_cluster.*,
  ( total_memory_used / vms_on ) AS avg_memory_per_vm,
  ( total_vms_vcpu_on / vms_on ) AS avg_vcpu_per_vm,
  ( total_memory_used / avg_memory_per_host ) AS ratio_memory,
  ( num_hosts - CEIL( total_memory_used / avg_memory_per_host ) ) AS supported_failures,
  ( total_memory_used / ( avg_memory_per_host * 0.8 ) ) AS ratio_memory_80,
  ( num_hosts - CEIL( total_memory_used / ( avg_memory_per_host * 0.8 ) ) ) AS supported_failures_80
FROM view_cluster;`
)
