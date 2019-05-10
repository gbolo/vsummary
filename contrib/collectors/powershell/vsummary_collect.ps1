#REQUIRES -Version 3.0

<#

This powershell script is under active development and is designed to be used in conjunction with:
vSummary - https://github.com/gbolo/vsummary


DESCRIPTION:
    The Function of this script is to retrieve data from vcenter;
    then send that data via http POST to a local/remote php server in JSON format.

TODO:
    - Merge get-view of same viewtype
    - Create a view into vmkernel interfaces
    - Create a view into VM Snapshots
    - Fix `#requires -PsSnapin VMware.VimAutomation.Core -Version 5` to work with 6
    - Alot more stuff I can't think of right now!

#>

function Send-VSVSummaryData ($json,$url) {
  # maybe add gzip and auth or api key?
  try {
    $request = Invoke-WebRequest -Uri $url -Body $json -ContentType "application/json" -Method Post -ErrorAction SilentlyContinue
  }
  catch [System.Net.WebException]{
    $request = $_.Exception.Response
    return 500
  }
  catch {
    Write-Error $_.Exception
    return 500
  }
  return $request.StatusCode
}


function Hash ($textToHash){
  $hasher = New-Object System.Security.Cryptography.SHA1Managed
  $toHash = [System.Text.Encoding]::UTF8.GetBytes($textToHash)
  $hashByteArray = $hasher.ComputeHash($toHash)
  foreach ($byte in $hashByteArray)
  {
    $res += $byte.ToString()
  }
  return $res;
}

function Get-VSVirtualMachine ([string]$vc_uuid) {

  & { Get-View -ViewType VirtualMachine -Property Name,
    Config.Files.VmPathName,
    Config.Hardware.NumCPU,
    Config.Hardware.MemoryMB,
    Config.GuestId,
    Config.Version,
    Config.Uuid,
    Config.instanceUuid,
    Config.changeVersion,
    Config.Template,
    Guest.ToolsVersion,
    Guest.ToolsRunningStatus,
    Guest.Hostname,
    Guest.IpAddress,
    Guest.GuestId,
    Parent,
    ParentVApp,
    ResourcePool,
    Summary.Quickstats.OverallCpuUsage,
    Summary.Quickstats.HostMemoryUsage,
    Summary.Quickstats.GuestMemoryUsage,
    Summary.Quickstats.UptimeSeconds,
    Runtime.PowerState,
    Runtime.Host | ForEach-Object {
      $vm = $_
      New-Object -TypeName PSobject -Property @{
        Name = $vm.Name
        moref = $vm.moref.Value
        vmx_path = $vm.Config.Files.VmPathName
        vcpu = $vm.Config.Hardware.NumCPU
        memory_mb = $vm.Config.Hardware.MemoryMB
        config_guest_os = $vm.Config.GuestId
        config_version = $vm.Config.Version
        smbios_uuid = $vm.Config.Uuid
        instance_uuid = $vm.Config.instanceUuid
        config_change_version = $vm.Config.changeVersion
        guest_tools_version = $vm.Guest.ToolsVersion
        guest_tools_running = $vm.Guest.ToolsRunningStatus
        guest_hostname = $vm.Guest.Hostname
        guest_ip = $vm.Guest.IpAddress
        guest_os = $vm.Guest.GuestId
        stat_cpu_usage = $vm.Summary.Quickstats.OverallCpuUsage
        stat_host_memory_usage = $vm.Summary.Quickstats.HostMemoryUsage
        stat_guest_memory_usage = $vm.Summary.Quickstats.GuestMemoryUsage
        stat_uptime_sec = $vm.Summary.Quickstats.UptimeSeconds
        power_state = [string]$vm.Runtime.PowerState
        esxi_moref = $vm.Runtime.Host.Value
        folder_moref = $vm.Parent.Value
        template = $vm.Config.template
        vapp_moref = $vm.ParentVApp.Value
        resourcepool_moref = $vm.ResourcePool.Value
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSResourcePool ([string]$vc_uuid) {

  & { Get-View -ViewType ResourcePool -Property Name,
    Owner,
    OverallStatus,
    Parent,
    Summary | ForEach-Object {
      $res = $_
      switch ($res.GetType().Name) {
        # ResourcePool
        "ResourcePool" {
          $type = "ResourcePool"
          $vapp_state = "n/a"

        }
        # vApp
        "VirtualApp" {
          $type = "VirtualApp"
          $vapp_state = $res.Summary.VAppState
        }
      }

      New-Object -TypeName PSobject -Property @{
        Name = $res.Name
        moref = $res.moref.Value
        type = $type
        status = [string]$res.OverallStatus
        vapp_state = $vapp_state
        parent_moref = $res.Parent.Value
        cluster_moref = $res.Owner.Value
        configured_memory_mb = $res.Summary.ConfiguredMemoryMB
        cpu_reservation = $res.Summary.Config.CpuAllocation.Reservation
        cpu_limit = $res.Summary.Config.CpuAllocation.Limit
        mem_reservation = $res.Summary.Config.MemoryAllocation.Reservation
        mem_limit = $res.Summary.Config.MemoryAllocation.Limit
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSVirtualNic ([string]$vc_uuid) {

  $dvs = @(Get-View -ViewType DistributedVirtualSwitch -Property Name,Uuid)

  $dvs | ForEach-Object { $_.UpdateViewData("Portgroup.Key","Portgroup.Name") }

  & { Get-View -ViewType VirtualMachine -Property Name,Config.Hardware.Device,Runtime.Host.Value | ForEach-Object {
      $vm = $_
      $vm_moref = $vm.moref.Value
      $esxi_moref = $vm.Runtime.Host.Value
      ## updated a bit of View data (to be used in the LinkedView properties later -- this is faster than using multiple Get-View calls for properties that are MoRefs themselves)
      ## source: https://communities.vmware.com/message/1887826
      $vm.UpdateViewData("Runtime.Host.ConfigManager.NetworkSystem.NetworkInfo.Vswitch","Runtime.Host.ConfigManager.NetworkSystem.NetworkInfo.ProxySwitch","Runtime.Host.ConfigManager.NetworkSystem.NetworkInfo.PortGroup")
      $vm.Config.Hardware.Device | Where-Object { $_ -is [VMware.Vim.VirtualEthernetCard] } | ForEach-Object {
        $vnic = $_
        $portgroup_name = $vswitch_type = $vswitch_name = $null

        $connected = $vnic.Connectable.Connected
        $status = $vnic.Connectable.status

        switch ($vnic.Backing.GetType().Name) {
          ## Standard vSwitch
          "VirtualEthernetCardNetworkBackingInfo" {
            $portgroup_moref = "null"
            $portgroup_name = $vnic.Backing.DeviceName
            $pg = $vm.Runtime.LinkedView.Host.ConfigManager.LinkedView.NetworkSystem.NetworkInfo.Portgroup | Where-Object { $_.Spec.Name -eq $vnic.Backing.DeviceName }
            $vswitch_name = $pg.Spec.VswitchName
            $vswitch_vm_obj = $vm.Runtime.LinkedView.Host.ConfigManager.LinkedView.NetworkSystem.NetworkInfo.Vswitch | Where-Object { $_.Key -eq $pg.Vswitch }
            $vswitch_type = if ($vswitch_vm_obj) { $vswitch_vm_obj.GetType().Name } else { "vSwitch type not found" }
            break;
          }
          ## DVS Switch
          "VirtualEthernetCardDistributedVirtualPortBackingInfo" {
            $dvs_vm_obj = $dvs | Where-Object { $_.Uuid -eq $vnic.Backing.Port.SwitchUuid }
            $pg = $dvs_vm_obj.LinkedView.Portgroup | Where-Object { $_.Key -eq $vnic.Backing.Port.PortgroupKey }
            $portgroup_moref = $pg.moref.Value
            $portgroup_name = $pg.Name
            $vswitch_name = $dvs_vm_obj.Name
            $vswitch_type = if ($dvs_vm_obj) { $dvs_vm_obj.GetType().Name } else { "dvSwitch type not found" }
            break;
          }
        }

        New-Object -TypeName PSobject -Property @{
          Name = $_.DeviceInfo.Label
          vm_moref = $vm_moref
          esxi_moref = $esxi_moref
          type = $_.GetType().Name
          mac = $_.MacAddress
          Connected = [string]$connected
          status = $status
          portgroup_name = $portgroup_name
          portgroup_moref = $portgroup_moref
          vswitch_type = $vswitch_type
          vswitch_name = $vswitch_name
          vcenter_id = $vc_uuid
        } ## end new-object
      } ## end foreach-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray



}

function Get-VSEsxi ([string]$vc_uuid) {

  & { Get-View -ViewType HostSystem -Property Name,
    Parent,
    Summary.MaxEVCModeKey,
    Summary.CurrentEVCModeKey,
    Summary.OverallStatus,
    Summary.Runtime.PowerState,
    Summary.Runtime.InMaintenanceMode,
    Summary.Hardware.Vendor,
    Summary.Hardware.Model,
    Summary.Hardware.Uuid,
    Summary.Hardware.MemorySize,
    Summary.Hardware.CpuModel,
    Summary.Hardware.CpuMhz,
    Summary.Hardware.NumCpuPkgs,
    Summary.Hardware.NumCpuCores,
    Summary.Hardware.NumCpuThreads,
    Summary.Hardware.NumNics,
    Summary.Hardware.NumHBAs,
    Summary.Config.Product.Version,
    Summary.Config.Product.Build,
    Summary.Quickstats.OverallCpuUsage,
    Summary.Quickstats.OverallMemoryUsage,
    Summary.Quickstats.Uptime | ForEach-Object {
      $esxi = $_
      New-Object -TypeName PSobject -Property @{
        Name = $esxi.Name
        moref = $esxi.moref.Value
        max_evc = $esxi.Summary.MaxEVCModeKey
        current_evc = $esxi.Summary.CurrentEVCModeKey
        status = [string]$esxi.Summary.OverallStatus
        power_state = [string]$esxi.Summary.Runtime.PowerState
        in_maintenance_mode = [string][boolean]$esxi.Summary.Runtime.InMaintenanceMode
        vendor = $esxi.Summary.Hardware.vendor
        model = $esxi.Summary.Hardware.model
        Uuid = $esxi.Summary.Hardware.Uuid
        memory_bytes = $esxi.Summary.Hardware.MemorySize
        cpu_model = $esxi.Summary.Hardware.CpuModel
        cpu_mhz = $esxi.Summary.Hardware.CpuMhz
        cpu_sockets = $esxi.Summary.Hardware.NumCpuPkgs
        cpu_cores = $esxi.Summary.Hardware.NumCpuCores
        cpu_threads = $esxi.Summary.Hardware.NumCpuThreads
        nics = $esxi.Summary.Hardware.NumNics
        hbas = $esxi.Summary.Hardware.NumHBAs
        Version = $esxi.Summary.Config.Product.Version
        build = $esxi.Summary.Config.Product.build
        stat_cpu_usage = $esxi.Summary.Quickstats.OverallCpuUsage
        stat_memory_usage = $esxi.Summary.Quickstats.OverallMemoryUsage
        stat_uptime_sec = $esxi.Summary.Quickstats.Uptime
        cluster_moref = $esxi.Parent.Value
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSPhysicalNic ([string]$vc_uuid) {

  & { Get-View -ViewType HostSystem -Property Name,
    Config.Network.Pnic | ForEach-Object {
      $esxi = $_
      $esxi.Config.Network.Pnic | ForEach-Object {
        $pnic = $_
        New-Object -TypeName PSobject -Property @{
          Name = $pnic.Device
          mac = $pnic.mac
          driver = $pnic.driver
          link_speed = $pnic.LinkSpeed.SpeedMB
          esxi_moref = $esxi.moref.Value
          vcenter_id = $vc_uuid
        } ## end new-object
      } ## end foreach-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}


function Get-VSStandardVswitch ([string]$vc_uuid) {

  $switch_type = "vSwitch"

  & { Get-View -ViewType HostSystem -Property Name,
    Config.Network.Vswitch | ForEach-Object {
      $esxi = $_
      $esxi.Config.Network.Vswitch | ForEach-Object {
        $vswitch = $_
        New-Object -TypeName PSobject -Property @{
          Name = $vswitch.Name
          ports = $vswitch.Spec.NumPorts
          max_mtu = $vswitch.Mtu
          esxi_moref = $esxi.moref.Value
          vcenter_id = $vc_uuid
          type = $switch_type
        } ## end new-object
      } ## end foreach-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSCluster ([string]$vc_uuid) {

  & { Get-View -ViewType ClusterComputeResource -Property Name,
    OverallStatus,
    Parent,
    Configuration.DasConfig,
    Configuration.DrsConfig,
    Summary | ForEach-Object {
      $cluster = $_
      New-Object -TypeName PSobject -Property @{
        Name = $cluster.Name
        moref = $cluster.moref.Value
        datacenter_moref = $cluster.Parent.Value
        total_cpu_threads = $cluster.Summary.NumCpuThreads
        total_cpu_mhz = $cluster.Summary.TotalCpu
        total_memory_bytes = $cluster.Summary.TotalMemory
        total_vmotions = $cluster.Summary.NumVmotions
        num_hosts = $cluster.Summary.NumHosts
        current_balance = $cluster.Summary.CurrentBalance
        target_balance = $cluster.Summary.TargetBalance
        drs_enabled = [string]$cluster.Configuration.DrsConfig.Enabled
        drs_behaviour = [string]$cluster.Configuration.DrsConfig.DefaultVmBehavior
        ha_enabled = [string]$cluster.Configuration.DasConfig.Enabled
        status = [string]$cluster.OverallStatus
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSDistributedVswitch ([string]$vc_uuid) {

  $switch_type = "DVS"

  & { Get-View -ViewType DistributedVirtualSwitch -Property Name,
    Summary.ProductInfo.Version,
    Config | ForEach-Object {
      $dvs = $_
      New-Object -TypeName PSobject -Property @{
        Name = $dvs.Name
        moref = $dvs.moref.Value
        Version = $dvs.Summary.ProductInfo.Version
        max_mtu = $dvs.Config.MaxMtu
        ports = $dvs.Config.NumPorts
        vcenter_id = $vc_uuid
        type = $switch_type
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSDatacenter ([string]$vc_uuid) {

  & { Get-View -ViewType Datacenter -Property Name,
    HostFolder,
    VmFolder | ForEach-Object {
      $dc = $_
      New-Object -TypeName PSobject -Property @{
        Name = $dc.Name
        moref = $dc.moref.Value
        vm_folder_moref = $dc.VmFolder.Value
        esxi_folder_moref = $dc.HostFolder.Value
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}


function Get-VSFolder ([string]$vc_uuid) {

  & { Get-View -ViewType Folder -Property Name,
    Parent,
    ChildType | ForEach-Object {
      $folder = $_
      New-Object -TypeName PSobject -Property @{
        Name = $folder.Name
        moref = $folder.moref.Value
        type = [string]$folder.ChildType
        parent_moref = $folder.Parent.Value
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}


function Get-VSDistributedPortGroup ([string]$vc_uuid) {

  $switch_type = "DVS"

  & { Get-View -ViewType DistributedVirtualPortgroup -Property Name,
    Config.DefaultPortConfig,
    Config.DistributedVirtualSwitch | ForEach-Object {
      $pg = $_
      $vlan_type = $pg.Config.DefaultPortConfig.Vlan.GetType().Name

      # single vlan id
      if ($vlan_type -eq "VmwareDistributedVirtualSwitchVlanIdSpec") {
        $vlan = [string]$pg.Config.DefaultPortConfig.Vlan.VlanId
        $vlan_start = "na"
        $vlan_end = "na"
      } elseif ($vlan_type -eq "VmwareDistributedVirtualSwitchTrunkVlanSpec") {
        $vlan = "na"
        $vlan_start = [string]$pg.Config.DefaultPortConfig.Vlan.VlanId.Start
        $vlan_end = [string]$pg.Config.DefaultPortConfig.Vlan.VlanId.End
      } else {
        $vlan = "TypeNotImplemented"
        $vlan_start = "na"
        $vlan_end = "na"
      }
      #SUPPORT IS NEEDED FOR VLAN TRUNKING AND OTHER TYPES
      # VmwareDistributedVirtualSwitchTrunkVlanSpec
      # $DVPG.Config.DefaultPortConfig.Vlan.VlanId.Start
      # $DVPG.Config.DefaultPortConfig.Vlan.VlanId.End
      New-Object -TypeName PSobject -Property @{
        Name = $pg.Name
        moref = $pg.moref.Value
        vlan_type = $vlan_type
        Vlan = $vlan
        vlan_start = $vlan_start
        vlan_end = $vlan_end
        vswitch_moref = $pg.Config.DistributedVirtualSwitch.Value
        vcenter_id = $vc_uuid
        type = $switch_type
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Get-VSStandardPortGroup ([string]$vc_uuid) {

  $switch_type = "vSwitch"

  & { Get-View -ViewType HostSystem -Property Name,
    Config.Network.Portgroup | ForEach-Object {
      $esxi = $_
      $esxi.Config.Network.Portgroup.Spec | ForEach-Object {
        $pg = $_
        New-Object -TypeName PSobject -Property @{
          Name = $pg.Name
          vswitch_name = $pg.VswitchName
          Vlan = [string]$pg.VlanId
          esxi_moref = $esxi.moref.Value
          vcenter_id = $vc_uuid
          type = $switch_type
        } ## end new-object
      } ## end foreach-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}


function Get-VSDatastore ([string]$vc_uuid) {

  & { Get-View -ViewType Datastore -Property Name,
    OverallStatus,
    Summary.Capacity,
    Summary.FreeSpace,
    Summary.Type,
    Summary.Uncommitted | ForEach-Object {
      $ds = $_
      New-Object -TypeName PSobject -Property @{
        Name = $ds.Name
        moref = $ds.moref.Value
        status = [string]$ds.OverallStatus
        capacity_bytes = $ds.Summary.Capacity
        free_bytes = $ds.Summary.FreeSpace
        uncommitted_bytes = $ds.Summary.Uncommitted
        type = $ds.Summary.type
        vcenter_id = $vc_uuid
      } ## end new-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}



function Get-VSVirtualDisk ([string]$vc_uuid) {

  & { Get-View -ViewType VirtualMachine -Property Name,
    Config.Hardware.Device,
    Config.instanceUuid,
    Runtime.Host | ForEach-Object {
      $vm = $_
      $vm.Config.Hardware.Device | Where-Object { $_ -is [VMware.Vim.VirtualDisk] } | ForEach-Object {
        $vdisk = $_

        ## Collect both capacity_bytes and capacityInKB since vm version vmx-07 and lower will not have capacity_bytes
        ## https://www.vmware.com/support/developer/converter-sdk/conv55_apireference/vim.vm.device.VirtualDisk.html
        New-Object -TypeName PSobject -Property @{
          Name = $vdisk.DeviceInfo.Label
          capacity_bytes = $vdisk.CapacityInBytes
          capacity_kb = $vdisk.capacityInKB
          path = $vdisk.Backing.Filename
          thin_provisioned = [string]$vdisk.Backing.ThinProvisioned
          datastore_moref = $vdisk.Backing.Datastore.Value
          Uuid = $vdisk.Backing.Uuid
          disk_object_id = $vdisk.diskObjectId
          vm_moref = $vm.moref.Value
          esxi_moref = $vm.Runtime.Host.Value
          vcenter_id = $vc_uuid
        } ## end new-object
      } ## end foreach-object
    } ## end foreach-object
  } | ConvertTo-Json -AsArray
}

function Invoke-VSFunctions ([string]$vc_uuid,[string]$url) {

  ## Create an OrderedDictionary
  $hshChecksToRun = [ordered]@{
    esxi = & Get-VSEsxi $vc_uuid
    #    PNIC = & Get-VSPhysicalNic $vc_uuid
    datastore = & Get-VSDatastore $vc_uuid
    virtualmachine = & Get-VSVirtualMachine $vc_uuid
    vswitch = (& Get-VSStandardVswitch $vc_uuid | ConvertFrom-Json) + (& Get-VSDistributedVswitch $vc_uuid | ConvertFrom-Json) | ConvertTo-Json -AsArray
    portgroup = (& Get-VSStandardPortGroup $vc_uuid | ConvertFrom-Json) + (& Get-VSDistributedPortGroup $vc_uuid | ConvertFrom-Json) | ConvertTo-Json -AsArray
    vnic = & Get-VSVirtualNic $vc_uuid
    vdisk = & Get-VSVirtualDisk $vc_uuid
    resourcepool = & Get-VSResourcePool $vc_uuid
    datacenter = & Get-VSDatacenter $vc_uuid
    # Folder check needs to be done after datacenter check
    folder = & Get-VSFolder $vc_uuid
    cluster = & Get-VSCluster $vc_uuid
  }

  ## Run Checks in OrderedDictionary
  $hshChecksToRun.Keys | ForEach-Object {
    $endpoint = $_
    $jsonData = $hshChecksToRun[$endpoint]
    # Check if $json is empty
    if (!$jsonData) {
      $vc_obj = New-Object -TypeName PSobject -Property @{
        vcenter_id = $vc_uuid
        objecttype = $endpoint
        Empty = 'DELETE'
      }
      $jsonData = $vc_obj | ConvertTo-Json -AsArray
    }
    $status = Send-VSVSummaryData $jsonData "$url/api/v2/$endpoint"
    Write-Verbose -Verbose "$endpoint check http status code: $status"
  }

}

# Load/Source the required configuration file
. .\vsummary_config.ps1

# Don't validate TLS certs
$pcliConfig = Get-PowerCLIConfiguration -Scope User
if ($pcliConfig.InvalidCertificateAction -ne "Ignore"){
  Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Scope User -confirm:$false
}

foreach ($vc in $vcenters.Keys){
  $vc_shortname = $vc
  $vc_fqdn = $vcenters.Item($vc).fqdn
  $vc_port = 443
  if ($vcenters.Item($vc).port) {
    $vc_port = $vcenters.Item($vc).port
  }
  $vc_user = $vcenters.Item($vc).readonly_user
  $vc_pass = $vcenters.Item($vc).password

  if ($global:DefaultVIServers.Count -gt 0) {
    Disconnect-VIServer -Server * -Force -Confirm:$false -WarningAction SilentlyContinue -ErrorAction SilentlyContinue | Out-Null
  }

  Write-Host "---"
  if ($c) { Remove-Variable c }
  $c = Connect-VIServer $vc_fqdn -port $vc_port -user $vc_user -password $vc_pass

  if ($c) {
    Write-Host "Connected to vCenter: $vc_fqdn"

    $vc_uuid = $c.instanceUuid
    $vc_obj = New-Object -TypeName PSobject -Property @{
      Name = $vc_shortname
      id = $vc_uuid
      Host = $vc_fqdn
    }
    $json = $vc_obj | ConvertTo-Json

    # SEND VCENTER INFO
    $status = Send-VSVSummaryData $json "$vsummary_url/api/v2/vcenter"
    Write-Verbose -Verbose "vcenter check http status code: $status"

    # SEND ALL CHECKS
    Invoke-VSFunctions $vc_uuid $vsummary_url

  } else {
    Write-Host "!! Could not connect to $vc_fqdn"
  }
}

if ($global:DefaultVIServers.Count -gt 0) {
  Disconnect-VIServer -Server * -Force -Confirm:$false
}
