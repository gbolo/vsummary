package common

type Vm struct {

	// These are part of BOTH API request AND db record
	VcenterId            string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name                 string `json:"name" db:"name" validate:"required"`
	Moref                string `json:"moref" db:"moref" validate:"required"`
	VmxPath              string `json:"vmx_path" db:"vmx_path" validate:"required"`
	Vcpu                 int32  `json:"vcpu" db:"vcpu" validate:"required"`
	MemoryMb             int32  `json:"memory_mb" db:"memory_mb" validate:"required"`
	ConfigGuestOs        string `json:"config_guest_os" db:"config_guest_os" validate:"required"`
	ConfigVersion        string `json:"config_version" db:"config_version" validate:"required"`
	SmbiosUuid           string `json:"smbios_uuid" db:"smbios_uuid" validate:"required"`
	InstanceUuid         string `json:"instance_uuid" db:"instance_uuid" validate:"required"`
	ConfigChangeVersion  string `json:"config_change_version" db:"config_change_version" validate:"required"`
	GuestToolsRunning    string `json:"guest_tools_running" db:"guest_tools_running" validate:"required"`
	GuestToolsVersion    string `json:"guest_tools_version" db:"guest_tools_version"`
	GuestHostname        string `json:"guest_hostname" db:"guest_hostname"`
	GuestIp              string `json:"guest_ip" db:"guest_ip"`
	GuestOs              string `json:"guest_os" db:"guest_os"`
	StatCpuUsage         int32  `json:"stat_cpu_usage" db:"stat_cpu_usage"`
	StatHostMemoryUsage  int32  `json:"stat_host_memory_usage" db:"stat_host_memory_usage"`
	StatGuestMemoryUsage int32  `json:"stat_guest_memory_usage" db:"stat_guest_memory_usage"`
	StatUptimeSec        int32  `json:"stat_uptime_sec" db:"stat_uptime_sec"`
	PowerState           string `json:"power_state" db:"power_state" validate:"required"`
	Template             bool   `json:"template" db:"template"`

	// These are part of API request ONLY
	ObjectType        string `json:"objecttype"`
	EsxiMoref         string `json:"esxi_moref" validate:"required"`
	FolderMoref       string `json:"folder_moref" validate:"required"`
	VappMoref         string `json:"vapp_moref"`
	ResourcePoolMoref string `json:"resourcepool_moref" validate:"required"`

	// These are part of db record ONLY
	Id             string `db:"id"`
	EsxiId         string `db:"esxi_id"`
	FolderId       string `db:"folder_id"`
	VappId         string `db:"vapp_id"`
	ResourcePoolId string `db:"resourcepool_id"`
}

type Datacenter struct {

	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id"`
	Name      string `json:"name" db:"name"`
	Moref     string `json:"moref" db:"moref"`

	// These are part of API request ONLY
	EsxiFolderMoref string `json:"esxi_folder_moref"`
	VmFolderMoref   string `json:"vm_folder_moref"`

	// These are part of db record ONLY
	Id           string `db:"id"`
	EsxiFolderId string `db:"esxi_folder_id"`
	VmFolderId   string `db:"vm_folder_id"`
}

type ResourcePool struct {

	// These are part of BOTH API request AND db record
	VcenterId          string `json:"vcenter_id" db:"vcenter_id"`
	Name               string `json:"name" db:"name"`
	Moref              string `json:"moref" db:"moref"`
	Type               string `json:"type" db:"type"`
	Status             string `json:"status" db:"status"`
	VappState          string `json:"vapp_state" db:"vapp_state"`
	ConfiguredMemoryMb string `json:"configured_memory_mb" db:"configured_memory_mb"`
	CpuReservation     string `json:"cpu_reservation" db:"cpu_reservation"`
	CpuLimit           string `json:"cpu_limit" db:"cpu_limit"`
	MemoryReservation  string `json:"mem_reservation" db:"mem_reservation"`
	MemoryLimit        string `json:"mem_limit" db:"mem_limit"`

	// These are part of API request ONLY
	ParentMoref  string `json:"parent_moref"`
	ClusterMoref string `json:"cluster_moref"`

	// These are part of db record ONLY
	Id         string `db:"id"`
	ClusterId  string `db:"cluster_id"`
	FullPath   string `db:"full_path"`
	VappInPath int    `db:"vapp_in_path"`
}

type VSwitch struct {

	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id"`
	Name      string `json:"name" db:"name"`
	Type      string `json:"type" db:"type"`
	Version   string `json:"version" db:"version"`
	Ports     int    `json:"ports" db:"ports"`
	MaxMTU    int    `json:"max_mtu" db:"max_mtu"`

	// These are part of API request ONLY
	Moref string `json:"moref" db:"moref"` // only distributed switch has this

	// These are part of db record ONLY
	Id     string `db:"id"`
	EsxiId string `db:"esxi_id"`
}

type Poller struct {

	// These are part of BOTH API request AND db record
	VcenterId   string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Username    string `json:"user_name" db:"user_name" validate:"required"`
	Password    string `json:"password" db:"password" validate:"required"`
	Enabled     bool   `json:"enabled" db:"enabled" validate:"required"`
	IntervalMin int    `json:"interval_min" db:"interval_min" validate:"required"`
}

//CREATE TABLE vswitch
//(
//id VARCHAR(32) PRIMARY KEY,
//name VARCHAR(128),
//type VARCHAR(64),
//version VARCHAR(32) DEFAULT null,
//max_mtu SMALLINT UNSIGNED DEFAULT 0,
//ports SMALLINT UNSIGNED DEFAULT 0,
//esxi_id VARCHAR(32) DEFAULT null,
//vcenter_id VARCHAR(36),
//present TINYINT DEFAULT 1
//);

//Function Get-VSDistributedVswitch ( [string]$vc_uuid ){
//
//$objecttype = "DVS"
//
//&{Get-View -ViewType DistributedVirtualSwitch -Property Name,
//Summary.ProductInfo.Version,
//Config | %{
//$dvs = $_
//New-Object -TypeName PSobject -Property @{
//name = $dvs.Name
//moref = $dvs.MoRef.Value
//version = $dvs.Summary.ProductInfo.Version
//max_mtu = $dvs.Config.MaxMtu
//ports = $dvs.Config.NumPorts
//vcenter_id = $vc_uuid
//objecttype = $objecttype
//} ## end new-object
//} ## end foreach-object
//} | ConvertTo-Json
//}

//$objecttype = "SVS"
//
//&{Get-View -ViewType HostSystem -Property Name,
//Config.Network.Vswitch | %{
//$esxi = $_
//$esxi.Config.Network.Vswitch | %{
//$vswitch = $_
//New-Object -TypeName PSobject -Property @{
//name = $vswitch.Name
//ports = $vswitch.Spec.NumPorts
//max_mtu = $vswitch.Mtu
//esxi_moref = $esxi.MoRef.Value
//vcenter_id = $vc_uuid
//objecttype = $objecttype
//} ## end new-object
//} ## end foreach-object
//} ## end foreach-object
//} | ConvertTo-Json

//New-Object -TypeName PSobject -Property @{
//name = $res.Name
//moref = $res.MoRef.Value
//type = $type
//status = $res.OverallStatus
//vapp_state = $vapp_state
//parent_moref = $res.Parent.Value
//cluster_moref = $res.Owner.Value
//configured_memory_mb = $res.Summary.ConfiguredMemoryMB
//cpu_reservation =  $res.summary.Config.CpuAllocation.Reservation
//cpu_limit = $res.summary.Config.CpuAllocation.Limit
//mem_reservation =  $res.summary.Config.MemoryAllocation.Reservation
//mem_limit = $res.summary.Config.MemoryAllocation.Limit
//vcenter_id = $vc_uuid
//objecttype = $objecttype

//CREATE TABLE resourcepool
//(
//id VARCHAR(32) PRIMARY KEY,
//moref VARCHAR(16),
//full_path VARCHAR(512),
//name VARCHAR(128),
//type VARCHAR(64),
//status VARCHAR(64),
//vapp_state VARCHAR(64),
//vapp_in_path TINYINT DEFAULT 0,
//configured_memory_mb BIGINT UNSIGNED,
//cpu_reservation BIGINT UNSIGNED,
//cpu_limit BIGINT,
//mem_reservation BIGINT UNSIGNED,
//mem_limit BIGINT,
//parent VARCHAR(32),
//parent_moref VARCHAR(16),
//cluster_id VARCHAR(32),
//vcenter_id VARCHAR(36),
//present TINYINT DEFAULT 1
//);
