package server

// Helper struct that defines the UI view
type UiView struct {
	Title        string
	AjaxEndpoint string

	// using a slice, will preserve the column order when ranging through it
	TableHeaders []tableColumnMap
}

// maps db columns to friendly names
// to be used by html table column titles and datatables json columns
type tableColumnMap struct {
	DbColumnName string
	FriendlyName string
}

// set the values for all UiView vars
var virtualMachineView = UiView{
	Title:        "Virtualmachines",
	AjaxEndpoint: "/api/dt/virtualmachines",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"folder", "Folder"},
		{"vcpu", "vCPU"},
		{"memory_mb", "Memory"},
		{"power_state", "PowerState"},
		{"guest_os", "Real GuestOS"},
		{"config_guest_os", "Config GuestOS"},
		{"config_version", "Version"},
		{"config_change_version", "ConfigChange"},
		{"guest_tools_version", "ToolsVersion"},
		{"guest_tools_running", "ToolRunning"},
		{"guest_hostname", "Hostname"},
		{"guest_ip", "IP"},
		{"cluster", "Cluster"},
		{"pool", "Pool"},
		{"datacenter", "Datacenter"},
		{"stat_cpu_usage", "CpuUsed"},
		{"stat_host_memory_usage", "HostMemUsed"},
		{"stat_guest_memory_usage", "GuestMemUsed"},
		{"stat_uptime_sec", "Uptime"},
		{"esxi_name", "ESXi"},
		{"esxi_current_evc", "ESXiEVC"},
		{"esxi_status", "ESXiStatus"},
		{"esxi_cpu_model", "ESXiCPU"},
		{"vdisks", "vDisks"},
		{"vnics", "vNICs"},
		{"vmx_path", "VMX"},
		{"vcenter_fqdn", "vCenter"},
		{"vcenter_short_name", "VC-ENV"},
	},
}

var esxiView = UiView{
	Title:        "ESXi",
	AjaxEndpoint: "/api/dt/esxi",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"max_evc", "MaxEVC"},
		{"current_evc", "EVC"},
		{"status", "Status"},
		{"power_state", "PowerState"},
		{"in_maintenance_mode", "Maintenance"},
		{"vendor", "Vendor"},
		{"model", "Model"},
		{"memory_bytes", "Memory"},
		{"cpu_model", "CPU"},
		{"cpu_mhz", "CpuMHZ"},
		{"cpu_sockets", "CpuSockets"},
		{"cpu_cores", "CpuCores"},
		{"cpu_threads", "CpuThreads"},
		{"nics", "NICs"},
		{"hbas", "HBAs"},
		{"version", "Version"},
		{"build", "Build"},
		{"stat_cpu_usage", "CpuUsed"},
		{"stat_memory_usage", "MemUsed"},
		{"stat_uptime_sec", "Uptime"},
		{"vms_powered_on", "VMsOn"},
		{"vcpus_powered_on", "vCPUs"},
		{"vmemory_mb_powered_on", "vRAM"},
		{"pnics", "pNICS"},
		{"cluster", "Cluster"},
		{"datacenter", "Datacenter"},
		{"vcenter_fqdn", "vCenter"},
		{"vcenter_short_name", "VC-ENV"},
	},
}

var portgroupView = UiView{
	Title:        "PortGroup",
	AjaxEndpoint: "/api/dt/portgroups",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"type", "Type"},
		{"vlan", "Vlan"},
		{"vlan_type", "VlanType"},
		{"vswitch_name", "vSwitch"},
		{"vswitch_type", "vSwitchType"},
		{"vswitch_max_mtu", "vSwitchMTU"},
		{"vnics", "vNics"},
		{"vcenter_fqdn", "vCenter"},
		{"vcenter_short_name", "VC-ENV"},
	},
}

var datastoreView = UiView{
	Title:        "Datastore",
	AjaxEndpoint: "/api/dt/datastores",
	TableHeaders: []tableColumnMap{
		{"name", "Name"},
		{"type", "Type"},
		{"status", "Status"},
		{"capacity_bytes", "Capacity"},
		{"free_bytes", "Free"},
	},
}
