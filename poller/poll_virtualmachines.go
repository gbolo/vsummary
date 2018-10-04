package poller

import (
	"context"
	"reflect"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	//"github.com/vmware/govmomi/vim25/types"
)

// this function returns VMs vDisks vNICs since they are all part of VirtualMachine managedObject
func (p *Poller) GetVirtualMachines() (VMs []common.VirtualMachine, vDisks []common.VDisk, vNICs []common.VNic, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollVirtualMachine")

	// Create view of VirtualMachine objects
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var moList []mo.VirtualMachine
	err = v.Retrieve(
		ctx,
		[]string{"VirtualMachine"},
		[]string{"summary", "config", "guest", "runtime", "parent", "resourcePool", "parentVApp"},
		&moList,
	)
	if err != nil {
		return
	}

	// Print summary per mo (see also: govc/mo/info.go)
	for _, mo := range moList {

		// create vm ---------------------------------------------------------------------------------------------------
		vm := common.VirtualMachine{
			Name:                 mo.Summary.Config.Name,
			Moref:                mo.Summary.Vm.Value,
			VmxPath:              mo.Config.Files.VmPathName,
			Vcpu:                 mo.Config.Hardware.NumCPU,
			MemoryMb:             mo.Config.Hardware.MemoryMB,
			ConfigGuestOs:        mo.Config.GuestId,
			ConfigVersion:        mo.Config.Version,
			SmbiosUuid:           mo.Config.Firmware,
			InstanceUuid:         mo.Config.Uuid,
			ConfigChangeVersion:  mo.Config.ChangeVersion,
			GuestToolsVersion:    mo.Guest.ToolsVersion,
			GuestToolsRunning:    mo.Guest.ToolsRunningStatus,
			GuestHostname:        mo.Guest.HostName,
			GuestIp:              mo.Guest.IpAddress,
			GuestOs:              mo.Guest.GuestId,
			StatCpuUsage:         mo.Summary.QuickStats.OverallCpuUsage,
			StatHostMemoryUsage:  mo.Summary.QuickStats.HostMemoryUsage,
			StatGuestMemoryUsage: mo.Summary.QuickStats.GuestMemoryUsage,
			StatUptimeSec:        mo.Summary.QuickStats.UptimeSeconds,
			PowerState:           string(mo.Runtime.PowerState),
			EsxiMoref:            mo.Runtime.Host.Value,
			Template:             mo.Config.Template,
			VcenterId:            v.Client().ServiceContent.About.InstanceUuid,
		}

		// folder may not exist
		if mo.Parent != nil {
			vm.FolderMoref = mo.Parent.Value
		}

		// vapps may not exist
		if mo.ParentVApp != nil {
			vm.VappMoref = mo.ParentVApp.Value
			vm.FolderId = "vapp"
		} else {
			vm.VappId = "none"
		}

		// resourcepool may not exist
		if mo.ResourcePool != nil {
			vm.ResourcePoolMoref = mo.ResourcePool.Value
		}

		VMs = append(VMs, vm)

		// loop through devices ----------------------------------------------------------------------------------------

		for _, device := range mo.Config.Hardware.Device {
			deviceType := reflect.TypeOf(device).String()

			// catch virtual disks
			if deviceType == "*types.VirtualDisk" {
				vdisk := common.VDisk{
					Name:                common.GetString(device, "DeviceInfo", "Label"),
					CapacityBytes:       common.GetInt(device, "CapacityInBytes"),
					CapacityKb:          common.GetInt(device, "CapacityInKB"),
					ThinProvisioned:     common.BoolToString(common.GetBool(device, "Backing", "ThinProvisioned")),
					DatastoreMoref:      common.GetString(device, "Backing", "Datastore", "Value"),
					Uuid:                common.GetString(device, "Backing", "Uuid"),
					DiskObjectId:        common.GetString(device, "DiskObjectId"),
					Path:                common.GetString(device, "Backing", "FileName"),
					EsxiMoref:           mo.Runtime.Host.Value,
					VcenterId:           v.Client().ServiceContent.About.InstanceUuid,
					VirtualmachineMoref: mo.Summary.Vm.Value,
					// TODO: add diskmode?
					//DiskMode: common.GetString(device, "Backing", "DiskMode"),
				}

				vDisks = append(vDisks, vdisk)
			}

			// catch virtual nics
			// TODO: should catch by github.com/vmware/govmomi/vim25/types.VirtualEthernetCard
			if deviceType == "*types.VirtualVmxnet3" ||
				deviceType == "*types.VirtualE1000" ||
				deviceType == "*types.VirtualE1000e" ||
				deviceType == "*types.VirtualPCNet32" {

				vnic := common.VNic{
					Name:                common.GetString(device, "DeviceInfo", "Label"),
					Type:                deviceType[7:],
					MacAddress:          common.GetString(device, "MacAddress"),
					Connected:           common.BoolToString(common.GetBool(device, "Connectable", "Connected")),
					Status:              common.GetString(device, "Connectable", "Status"),
					VirtualmachineMoref: mo.Summary.Vm.Value,
					EsxiMoref:           mo.Runtime.Host.Value,
					VcenterId:           v.Client().ServiceContent.About.InstanceUuid,
				}

				// if Backing.Port exists, then this is a DVS, or else its a vswitch
				if common.CheckIfKeyExists(device, "Backing", "Port") {
					vnic.VswitchType = "VmwareDistributedVirtualSwitch"
					vnic.PortgroupMoref = common.GetString(device, "Backing", "Port", "PortgroupKey")
				} else {
					vnic.VswitchType = "HostVirtualSwitch"
					vnic.PortgroupName = common.GetString(device, "Backing", "DeviceName")
				}

				vNICs = append(vNICs, vnic)
			}
		}
	}

	log.Infof("poller fetched summary of %d moList", len(VMs))
	return

}
