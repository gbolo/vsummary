package poller

import (
	"context"
	"strings"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetFolders() (list []common.Folder, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollFolders")

	// Create view for objects
	moType := "Folder"
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{moType}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var molist []mo.Folder
	err = v.Retrieve(
		ctx,
		[]string{moType},
		[]string{"name", "parent", "childType"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		object := common.Folder{
			Name:        mo.Name,
			Moref:       mo.Self.Value,
			VcenterId:   v.Client().ServiceContent.About.InstanceUuid,
			ParentMoref: mo.Parent.Value,
			Type:        strings.Join(mo.ChildType, " "),
		}

		list = append(list, object)

	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return

}
