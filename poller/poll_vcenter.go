package poller

import (
	"context"
	"time"

	"net/url"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
)

func (p *Poller) GetVcenter() (vcenter common.VCenter, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollVcenter")

	// Create view for objects
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	vcenter.Id = v.Client().ServiceContent.About.InstanceUuid
	url, err := url.Parse(p.Config.URL)
	if err == nil {
		vcenter.Host = url.Host
	}

	return

}
