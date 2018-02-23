package poller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetDatacenters() (dcList []common.Datacenter, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollDatacenters")

	// Create view for objects
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var dcs []mo.Datacenter
	err = v.Retrieve(
		ctx,
		[]string{"Datacenter"},
		[]string{"name", "hostFolder", "vmFolder"},
		&dcs,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, dc := range dcs {

		dcStruct := common.Datacenter{
			Name:            dc.Name,
			Moref:           dc.Self.Value,
			EsxiFolderMoref: dc.HostFolder.Value,
			VmFolderMoref:   dc.VmFolder.Value,
			VcenterId:       v.Client().ServiceContent.About.InstanceUuid,
		}

		dcList = append(dcList, dcStruct)

	}

	log.Infof("poller fetched summary of %d datacenters", len(dcList))
	return

}

func (p *Poller) SendDatacenters() (err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "requestDatacenters")

	// get objects
	datacenters, err := p.GetDatacenters()
	if err != nil {
		log.Debugf("failed to retrieve Datacenter list: %s", err)
		return
	}

	log.Infof("poller sending summary of %d datacenter(s)", len(datacenters))

	jsonObj, err := json.Marshal(datacenters)
	if err != nil {
		log.Errorf("invalid json datacenter: %s", err)
		return
	}

	err = sendResults("/datacenter", jsonObj)
	if err != nil {
		return
	}

	return

}
