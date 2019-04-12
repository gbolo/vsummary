package common

type VirtualMachine struct {

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
	FolderMoref       string `json:"folder_moref""`
	VappMoref         string `json:"vapp_moref"`
	ResourcePoolMoref string `json:"resourcepool_moref"`

	// These are part of db record ONLY
	Id             string `db:"id"`
	EsxiId         string `db:"esxi_id"`
	FolderId       string `db:"folder_id"`
	VappId         string `db:"vapp_id"`
	ResourcePoolId string `db:"resourcepool_id"`
}

type Datacenter struct {

	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name      string `json:"name" db:"name" validate:"required"`
	Moref     string `json:"moref" db:"moref" validate:"required"`

	// These are part of API request ONLY
	EsxiFolderMoref string `json:"esxi_folder_moref"`
	VmFolderMoref   string `json:"vm_folder_moref"`

	// These are part of db record ONLY
	Id           string `db:"id"`
	EsxiFolderId string `db:"esxi_folder_id"`
	VmFolderId   string `db:"vm_folder_id"`
}

type Cluster struct {

	// These are part of BOTH API request AND db record
	VcenterId        string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name             string `json:"name" db:"name" validate:"required"`
	TotalCpuThreads  int16  `json:"total_cpu_threads" db:"total_cpu_threads"`
	TotalCpuMhz      int32  `json:"total_cpu_mhz" db:"total_cpu_mhz"`
	TotalMemoryBytes int64  `json:"total_memory_bytes" db:"total_memory_bytes"`
	TotalVmotions    int32  `json:"total_vmotions" db:"total_vmotions"`
	NumHosts         int32  `json:"num_hosts" db:"num_hosts"`
	DRSEnabled       string `json:"drs_enabled" db:"drs_enabled"`
	DRSBehaviour     string `json:"drs_behaviour" db:"drs_behaviour"`
	HAEnabled        string `json:"ha_enabled" db:"ha_enabled"`
	CurrentBalance   int32  `json:"current_balance" db:"current_balance"`
	TargetBalance    int32  `json:"target_balance" db:"target_balance"`
	Status           string `json:"status" db:"status"`

	// These are part of API request ONLY
	Moref           string `json:"moref" db:"moref" validate:"required"`
	DatacenterMoref string `json:"datacenter_moref" validate:"required"`

	// These are part of db record ONLY
	Id           string `db:"id"`
	DatacenterId string `db:"datacenter_id"`
}

type ResourcePool struct {

	// These are part of BOTH API request AND db record
	VcenterId          string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name               string `json:"name" db:"name" validate:"required"`
	Moref              string `json:"moref" db:"moref" validate:"required"`
	Type               string `json:"type" db:"type"`
	Status             string `json:"status" db:"status"`
	VappState          string `json:"vapp_state" db:"vapp_state"`
	ConfiguredMemoryMb int64  `json:"configured_memory_mb" db:"configured_memory_mb"`
	CpuReservation     int64  `json:"cpu_reservation" db:"cpu_reservation"`
	CpuLimit           int64  `json:"cpu_limit" db:"cpu_limit"`
	MemoryReservation  int64  `json:"mem_reservation" db:"mem_reservation"`
	MemoryLimit        int64  `json:"mem_limit" db:"mem_limit"`
	ParentMoref        string `json:"parent_moref" db:"parent_moref"`

	// These are part of API request ONLY
	ClusterMoref string `json:"cluster_moref"`

	// These are part of db record ONLY
	Id         string `db:"id"`
	Parent     string `db:"parent"`
	ClusterId  string `db:"cluster_id"`
	FullPath   string `db:"full_path"`
	VappInPath int    `db:"vapp_in_path"`
}

type Poller struct {

	// These are part of BOTH API request AND db record
	VcenterHost string `json:"vcenter_host" db:"vcenter_host" validate:"required" mapstructure:"hostname"`
	VcenterName string `json:"vcenter_name" db:"vcenter_name" validate:"required" mapstructure:"environment"`
	Username    string `json:"user_name" db:"user_name" validate:"required" mapstructure:"username"`
	Password    string `json:"password" db:"password" validate:"required"`
	Enabled     bool   `json:"enabled" db:"enabled" validate:"required"`
	IntervalMin int    `json:"interval_min" db:"interval_min" validate:"required"`

	// These are part of db record ONLY
	Id       string `db:"id"`
	Internal bool   `db:"internal"`
	LastPoll string `db:"last_poll"`

	// This is used by external poller only
	PlainTextPassword string `mapstructure:"password"`
}

type Esxi struct {

	// These are part of BOTH API request AND db record
	VcenterId         string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name              string `json:"name" db:"name" validate:"required"`
	MaxEvc            string `json:"max_evc" db:"max_evc"`
	CurrentEvc        string `json:"current_evc" db:"current_evc"`
	PowerState        string `json:"power_state" db:"power_state" validate:"required"`
	InMaintenanceMode string `json:"in_maintenance_mode" db:"in_maintenance_mode" validate:"required"`
	Vendor            string `json:"vendor" db:"vendor"`
	Model             string `json:"model" db:"model"`
	Uuid              string `json:"uuid" db:"uuid"`
	MemoryBytes       int64  `json:"memory_bytes" db:"memory_bytes"`
	CpuModel          string `json:"cpu_model" db:"cpu_model"`
	CpuMhz            int32  `json:"cpu_mhz" db:"cpu_mhz"`
	CpuSockets        int16  `json:"cpu_sockets" db:"cpu_sockets"`
	CpuThreads        int16  `json:"cpu_threads" db:"cpu_threads"`
	CpuCores          int16  `json:"cpu_cores" db:"cpu_cores"`
	Nics              int32  `json:"nics" db:"nics"`
	Hbas              int32  `json:"hbas" db:"hbas"`
	Version           string `json:"version" db:"version"`
	Build             string `json:"build" db:"build"`
	StatCpuUsage      int32  `json:"stat_cpu_usage" db:"stat_cpu_usage"`
	StatMemoryUsage   int32  `json:"stat_memory_usage" db:"stat_memory_usage"`
	StatUptimeSec     int32  `json:"stat_uptime_sec" db:"stat_uptime_sec"`
	Status            string `json:"status" db:"status" validate:"required"`

	// These are part of API request ONLY
	Moref        string `json:"moref" db:"moref" validate:"required"`
	ClusterMoref string `json:"cluster_moref"`

	// These are part of db record ONLY
	Id        string `db:"id"`
	ClusterId string `db:"cluster_id"`
}

type VCenter struct {

	// These are part of BOTH API request AND db record
	Id   string `json:"id" db:"id" validate:"required"`
	Host string `json:"host" db:"host" validate:"required"`
	Name string `json:"name" db:"name"`
}

type Folder struct {
	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name      string `json:"name" db:"name" validate:"required"`
	Type      string `json:"type" db:"type"`

	// These are part of API request ONLY
	Moref       string `json:"moref" validate:"required"`
	ParentMoref string `json:"parent_moref"`

	// These are part of db record ONLY
	Id                 string `db:"id"`
	FullPath           string `db:"full_path"`
	Parent             string `db:"parent"`
	ParentDatacenterId string `db:"parent_datacenter_id"`
}

type VDisk struct {
	// These are part of BOTH API request AND db record
	VcenterId       string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name            string `json:"name" db:"name" validate:"required"`
	CapacityBytes   int64  `json:"capacity_bytes" db:"capacity_bytes" validate:"required"`
	Path            string `json:"path" db:"path"`
	ThinProvisioned string `json:"thin_provisioned" db:"thin_provisioned"`
	Uuid            string `json:"uuid" db:"uuid"`
	DiskObjectId    string `json:"disk_object_id" db:"disk_object_id"`

	// These are part of API request ONLY
	CapacityKb          int64  `json:"capacity_kb"`
	DatastoreMoref      string `json:"datastore_moref"`
	VirtualmachineMoref string `json:"vm_moref"`
	EsxiMoref           string `json:"esxi_moref"`

	// These are part of db record ONLY
	Id               string `db:"id"`
	DatastoreId      string `db:"datastore_id"`
	VirtualMachineId string `db:"vm_id"`
	EsxiId           string `db:"esxi_id"`
}

type PNic struct {
	// These are part of BOTH API request AND db record
	VcenterId  string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name       string `json:"name" db:"name" validate:"required"`
	MacAddress string `json:"mac" db:"mac"`
	LinkSpeed  int32  `json:"link_speed" db:"link_speed"`
	Driver     string `json:"driver" db:"driver"`

	// These are part of API request ONLY
	EsxiMoref string `json:"esxi_moref" validate:"required"`

	// These are part of db record ONLY
	Id     string `db:"id"`
	EsxiId string `db:"esxi_id"`

	// TODO: unused for now
	VswicthId string `db:"vswitch_id"`
}

type VNic struct {
	// These are part of BOTH API request AND db record
	VcenterId  string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name       string `json:"name" db:"name" validate:"required"`
	MacAddress string `json:"mac" db:"mac" validate:"required"`
	Connected  string `json:"connected" db:"connected" validate:"required"`
	Status     string `json:"status" db:"status" validate:"required"`
	Type       string `json:"type" db:"type"`

	// These are part of API request ONLY
	Moref               string `json:"moref"`
	PortgroupMoref      string `json:"portgroup_moref"`
	PortgroupName       string `json:"portgroup_name"`
	VirtualmachineMoref string `json:"vm_moref" validate:"required"`
	EsxiMoref           string `json:"esxi_moref"`
	VswitchType         string `json:"vswitch_type"`
	VswitchName         string `json:"vswitch_name"`

	// These are part of db record ONLY
	Id               string `db:"id"`
	VirtualmachineId string `db:"vm_id"`
	PortgroupId      string `db:"portgroup_id"`
}

type VSwitch struct {
	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name      string `json:"name" db:"name" validate:"required"`
	MaxMtu    int32  `json:"max_mtu" db:"max_mtu"`
	Ports     int32  `json:"ports" db:"ports"`
	Version   string `json:"version" db:"version"`

	// These are part of API request ONLY
	Moref     string `json:"moref"`
	EsxiMoref string `json:"esxi_moref"`
	Type      string `json:"type"`

	// These are part of db record ONLY
	Id     string `db:"id"`
	EsxiId string `db:"esxi_id"`
}

type Portgroup struct {
	// These are part of BOTH API request AND db record
	VcenterId string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name      string `json:"name" db:"name" validate:"required"`
	Type      string `json:"type" db:"type"`
	VlanType  string `json:"vlan_type" db:"vlan_type"`
	Vlan      string `json:"vlan"`

	// These are part of API request ONLY
	Moref        string `json:"moref"`
	EsxiMoref    string `json:"esxi_moref"`
	VswitchMoref string `json:"vswitch_moref"`
	VswitchName  string `json:"vswitch_name"`

	// These are part of db record ONLY
	Id        string `db:"id"`
	VswitchId string `db:"vswitch_id"`
}

type Datastore struct {
	// These are part of BOTH API request AND db record
	VcenterId        string `json:"vcenter_id" db:"vcenter_id" validate:"required"`
	Name             string `json:"name" db:"name" validate:"required"`
	Moref            string `json:"moref" validate:"required"`
	Status           string `json:"status" db:"status" validate:"required"`
	CapacityBytes    int64  `json:"capacity_bytes" db:"capacity_bytes"`
	FreeBytes        int64  `json:"free_bytes" db:"free_bytes"`
	UncommittedBytes int64  `json:"uncommitted_bytes" db:"uncommitted_bytes"`
	Type             string `json:"type" db:"type"`

	// These are part of db record ONLY
	Id string `db:"id"`
}

// TODO: unused right now
type VmkNic struct {
	// These are part of BOTH API request AND db record
	VcenterId  string `json:"vcenter_id" db:"vcenter_id"`
	Name       string `json:"name" db:"name"`
	MacAddress string `json:"mac" db:"mac"`
	IP         string `json:"ip" db:"ip"`
	Netmask    string `json:"netmask" db:"netmask"`

	// These are part of API request ONLY
	Moref          string `json:"moref"`
	PortgroupMoref string `json:"portgroup_moref"`
	EsxiMoref      string `json:"esxi_moref"`

	// These are part of db record ONLY
	Id          string `db:"id"`
	PortgroupId string `db:"portgrouo_id"`
	EsxiId      string `db:"esxi_id"`
}
