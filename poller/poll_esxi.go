package poller

import (
	"context"
	"reflect"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetEsxi() (esxiList []common.Esxi, pNics []common.PNic, vSwitches []common.VSwitch, vSwitchPortgroups []common.Portgroup, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollEsxi")

	// Create view for objects
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	// https://godoc.org/github.com/vmware/govmomi/vim25/mo#HostSystem
	var esxis []mo.HostSystem
	err = v.Retrieve(
		ctx,
		[]string{"HostSystem"},
		[]string{"name", "parent", "summary", "config"},
		&esxis,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, esxi := range esxis {

		clStruct := common.Esxi{
			Name:              esxi.Name,
			Moref:             esxi.Self.Value,
			VcenterId:         v.Client().ServiceContent.About.InstanceUuid,
			Status:            string(esxi.Summary.OverallStatus),
			InMaintenanceMode: common.BoolToString(esxi.Summary.Runtime.InMaintenanceMode),
			MaxEvc:            esxi.Summary.MaxEVCModeKey,
			CurrentEvc:        esxi.Summary.CurrentEVCModeKey,
			PowerState:        string(esxi.Summary.Runtime.PowerState),
			Vendor:            esxi.Summary.Hardware.Vendor,
			Model:             esxi.Summary.Hardware.Model,
			Uuid:              esxi.Summary.Hardware.Uuid,
			MemoryBytes:       esxi.Summary.Hardware.MemorySize,
			CpuModel:          esxi.Summary.Hardware.CpuModel,
			CpuMhz:            esxi.Summary.Hardware.CpuMhz,
			CpuThreads:        esxi.Summary.Hardware.NumCpuThreads,
			CpuSockets:        esxi.Summary.Hardware.NumCpuPkgs,
			CpuCores:          esxi.Summary.Hardware.NumCpuCores,
			Nics:              esxi.Summary.Hardware.NumNics,
			Hbas:              esxi.Summary.Hardware.NumHBAs,
			Version:           esxi.Summary.Config.Product.Version,
			Build:             esxi.Summary.Config.Product.Build,
			StatCpuUsage:      esxi.Summary.QuickStats.OverallCpuUsage,
			StatMemoryUsage:   esxi.Summary.QuickStats.OverallMemoryUsage,
			StatUptimeSec:     esxi.Summary.QuickStats.Uptime,
			ClusterMoref:      esxi.Parent.Value,
		}

		esxiList = append(esxiList, clStruct)

		// avoid nil pointers!
		if esxi.Config != nil && esxi.Config.Network.Pnic != nil {

			// Get Physical Network Cards
			for _, pnic := range esxi.Config.Network.Pnic {

				if reflect.TypeOf(pnic).String() == "types.PhysicalNic" {

					newPnic := common.PNic{
						Name:       pnic.Device,
						MacAddress: pnic.Mac,
						Driver:     pnic.Driver,
						EsxiMoref:  esxi.Self.Value,
						VcenterId:  v.Client().ServiceContent.About.InstanceUuid,
					}

					// this may not be reported. avoid nil pointers!
					if pnic.LinkSpeed != nil {
						newPnic.LinkSpeed = pnic.LinkSpeed.SpeedMb
					}

					pNics = append(pNics, newPnic)
				}
			}
		}

		// avoid nil pointers!
		if esxi.Config != nil && esxi.Config.Network.Vswitch != nil {

			// Get Standard vSwitches
			for _, svswitch := range esxi.Config.Network.Vswitch {

				if reflect.TypeOf(svswitch).String() == "types.HostVirtualSwitch" {

					vSwitches = append(vSwitches, common.VSwitch{
						Type:      "SVS",
						Name:      svswitch.Name,
						Ports:     svswitch.Spec.NumPorts,
						MaxMtu:    svswitch.Mtu,
						EsxiMoref: esxi.Self.Value,
						VcenterId: v.Client().ServiceContent.About.InstanceUuid,
					})
				}
			}
		}

		// avoid nil pointers!
		if esxi.Config != nil && esxi.Config.Network.Portgroup != nil {

			// Get Standard vSwitch Portgroups
			for _, spg := range esxi.Config.Network.Portgroup {

				if reflect.TypeOf(spg).String() == "types.HostPortGroup" {

					vSwitchPortgroups = append(vSwitchPortgroups, common.Portgroup{
						Type:        "vSwitch",
						Name:        spg.Spec.Name,
						VswitchName: spg.Spec.VswitchName,
						Vlan:        spg.Spec.VlanId,
						EsxiMoref:   esxi.Self.Value,
						VcenterId:   v.Client().ServiceContent.About.InstanceUuid,
					})
				}
			}
		}

	}

	log.Infof(
		"poller fetched summary of %d esxi hosts, %d pNICS, %d standard vswitches, %d standard portgroups",
		len(esxiList),
		len(pNics),
		len(vSwitches),
		len(vSwitchPortgroups),
	)

	return

}
