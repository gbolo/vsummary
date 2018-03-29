package poller

import (
	"context"
	"time"

	//"github.com/gbolo/go-util/lib/debugging"
	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetEsxi() (clList []common.Esxi, err error) {

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
		[]string{"name", "parent", "summary"},
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

		clList = append(clList, clStruct)

	}

	log.Infof("poller fetched summary of %d esxi(s)", len(clList))
	return

}
