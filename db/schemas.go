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

		// TODO: add views
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
     moref                   VARCHAR(32),
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
     vcenter_host   VARCHAR(64) PRIMARY KEY,
     vcenter_name   VARCHAR(64),
     enabled        TINYINT DEFAULT 1,
     user_name      VARCHAR(128),
     password       VARCHAR(256),
     interval_min   INT UNSIGNED
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
     moref                VARCHAR(16),
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
     parent_moref         VARCHAR(16),
     cluster_id           VARCHAR(32),
     vcenter_id           VARCHAR(36),
     present              TINYINT DEFAULT 1
  );`

	schemaFolder = `
CREATE TABLE IF NOT EXISTS folder
  (
     id                   VARCHAR(32) PRIMARY KEY,
     moref                VARCHAR(32),
     NAME                 VARCHAR(128),
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
     KEY 'vmid_ix' ('vm_id')
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
     KEY 'vmid_ix' ('vm_id')
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
     moref             VARCHAR(16),
     status            VARCHAR(32),
     capacity_bytes    BIGINT UNSIGNED,
     free_bytes        BIGINT UNSIGNED,
     uncommitted_bytes BIGINT UNSIGNED,
     type              VARCHAR(32),
     vcenter_id        VARCHAR(36),
     present           TINYINT DEFAULT 1
  );`
)
