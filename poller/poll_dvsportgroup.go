package poller

import (
	"context"
	"time"

	//"fmt"
	//"reflect"
	//
	//"github.com/gbolo/go-util/lib/debugging"
	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetDVSPortgroups() (list []common.Portgroup, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollDatastores")

	// Create view for objects
	moType := "DistributedVirtualPortgroup"
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{moType}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var molist []mo.DistributedVirtualPortgroup
	err = v.Retrieve(
		ctx,
		[]string{moType},
		[]string{"name", "config"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		// TODO: this needs to be cleaned up
		list = append(list, common.Portgroup{
			Name:      mo.Name,
			Moref:     mo.Self.Value,
			Type:      "DVS",
			VcenterId: v.Client().ServiceContent.About.InstanceUuid,
		})

		//mo.Config.DefaultPortConfig
		//debugging.PrettyPrint(mo.Config.DefaultPortConfig.GetDVPortSetting())
		//fmt.Println("test-->", reflect.TypeOf(mo.Config.DefaultPortConfig.GetDVPortSetting()).String())
	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return

}
