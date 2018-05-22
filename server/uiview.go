package server

// Helper struct that defines the UI view
type UiView struct {
	Title        string
	AjaxEndpoint string
	Table        map[string]string

	// json column names for datatatbles
	DtColumns []string
}

// set the values for all UiView vars
var virtualMachineView = UiView{
	Title:        "Virtualmachines",
	AjaxEndpoint: "/api/dt/virtualmachines",
	Table: map[string]string{
		"name":                    "Name",
		"folder":                  "Folder",
		"vcpu":                    "vCPU",
		"memory_mb":               "Memory",
		"power_state":             "PowerState",
		"guest_os":                "Real GuestOS",
		"config_guest_os":         "Config GuestOS",
		"config_version":          "Version",
		"config_change_version":   "ConfigChange",
		"guest_tools_version":     "ToolsVersion",
		"guest_tools_running":     "ToolRunning",
		"guest_hostname":          "Hostname",
		"guest_ip":                "IP",
		"cluster":                 "Cluster",
		"pool":                    "Pool",
		"datacenter":              "Datacenter",
		"stat_cpu_usage":          "CpuUsed",
		"stat_host_memory_usage":  "HostMemUsed",
		"stat_guest_memory_usage": "GuestMemUsed",
		"stat_uptime_sec":         "Uptime",
		"esxi_name":               "ESXi",
		"esxi_current_evc":        "ESXiEVC",
		"esxi_status":             "ESXiStatus",
		"esxi_cpu_model":          "ESXiCPU",
		"vdisks":                  "vDisks",
		"vnics":                   "vNICs",
		"vmx_path":                "VMX",
		"vcenter_fqdn":            "vCenter",
		"vcenter_short_name":      "VC-ENV",
	},
}

func init() {

	// set the DtColumns for all uiviews
	setDtColumns(&virtualMachineView)
}
