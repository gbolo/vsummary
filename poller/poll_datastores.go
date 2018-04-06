package poller

import (
	"context"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetDatastores() (list []common.Datastore, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollDatastores")

	// Create view for objects
	moType := "Datastore"
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{moType}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var molist []mo.Datastore
	err = v.Retrieve(
		ctx,
		[]string{moType},
		[]string{"name", "summary", "overallStatus"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		object := common.Datastore{
			Name:             mo.Name,
			Moref:            mo.Self.Value,
			VcenterId:        v.Client().ServiceContent.About.InstanceUuid,
			Status:           string(mo.OverallStatus),
			CapacityBytes:    mo.Summary.Capacity,
			FreeBytes:        mo.Summary.FreeSpace,
			UncommittedBytes: mo.Summary.Uncommitted,
			Type:             mo.Summary.Type,
		}

		list = append(list, object)

	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return

}
