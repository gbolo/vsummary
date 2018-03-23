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

type Cluster struct {

	// These are part of BOTH API request AND db record
	VcenterId        string `json:"vcenter_id" db:"vcenter_id"`
	Name             string `json:"name" db:"name"`
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
	Moref           string `json:"moref" db:"moref"`
	DatacenterMoref string `json:"datacenter_moref"`

	// These are part of db record ONLY
	Id           string `db:"id"`
	DatacenterId string `db:"datacenter_id"`
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
	VcenterHost string `json:"vcenter_host" db:"vcenter_host" validate:"required"`
	VcenterName string `json:"vcenter_name" db:"vcenter_name" validate:"required"`
	Username    string `json:"user_name" db:"user_name" validate:"required"`
	Password    string `json:"password" db:"password" validate:"required"`
	Enabled     bool   `json:"enabled" db:"enabled" validate:"required"`
	IntervalMin int    `json:"interval_min" db:"interval_min" validate:"required"`
}

type Esxi struct {

	// These are part of BOTH API request AND db record
	VcenterId         string `json:"vcenter_id" db:"vcenter_id"`
	Name              string `json:"name" db:"name"`
	MaxEvc            string `json:"max_evc" db:"max_evc"`
	CurrentEvc        string `json:"current_evc" db:"current_evc"`
	PowerState        string `json:"power_state" db:"power_state"`
	InMaintenanceMode string `json:"in_maintenance_mode" db:"in_maintenance_mode"`
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
	Status            string `json:"status" db:"status"`

	// These are part of API request ONLY
	Moref        string `json:"moref" db:"moref"`
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
