package poller

import (
	"context"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetResourcepools() (list []common.ResourcePool, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollResourcepools")

	// Create view for objects
	moType := "ResourcePool"
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{moType}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var molist []mo.ResourcePool
	err = v.Retrieve(
		ctx,
		[]string{moType},
		[]string{"name", "owner", "parent", "runtime", "summary"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		object := common.ResourcePool{
			Type:               "ResourcePool",
			Name:               mo.Name,
			Moref:              mo.Self.Value,
			VcenterId:          v.Client().ServiceContent.About.InstanceUuid,
			Status:             string(mo.Runtime.OverallStatus),
			VappState:          "None",
			ParentMoref:        mo.Parent.Value,
			ConfiguredMemoryMb: common.GetInt(mo, "Summary", "ConfiguredMemoryMB"),
			CpuReservation:     common.GetInt(mo, "Summary", "Config", "CpuAllocation", "Reservation"),
			CpuLimit:           common.GetInt(mo, "Summary", "Config", "CpuAllocation", "Limit"),
			MemoryReservation:  common.GetInt(mo, "Summary", "Config", "MemoryAllocation", "Reservation"),
			MemoryLimit:        common.GetInt(mo, "Summary", "Config", "MemoryAllocation", "Limit"),
		}

		if mo.Owner.Type == "ClusterComputeResource" {
			object.ClusterMoref = mo.Owner.Value
		}

		list = append(list, object)

	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return

}
