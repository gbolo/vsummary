package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
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
		[]string{"name", "config", "config.defaultPortConfig"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		// TODO: this needs to be cleaned up
		object := common.Portgroup{
			Name:         mo.Name,
			Moref:        mo.Self.Value,
			Type:         "DVS",
			VswitchMoref: mo.Config.DistributedVirtualSwitch.Value,
			VcenterId:    v.Client().ServiceContent.About.InstanceUuid,
		}

		if common.CheckIfKeyExists(mo.Config.DefaultPortConfig, "Vlan") {

			pconfig := mo.Config.DefaultPortConfig
			switch vl := pconfig.(type) {
			case *types.VMwareDVSPortSetting:
				vlan := vl.Vlan
				switch vlanspec := vlan.(type) {

				case *types.VmwareDistributedVirtualSwitchVlanIdSpec:
					object.VlanType = "VmwareDistributedVirtualSwitchVlanIdSpec"
					object.Vlan = fmt.Sprint(common.GetInt(mo.Config.DefaultPortConfig, "Vlan", "VlanId"))

				case *types.VmwareDistributedVirtualSwitchTrunkVlanSpec:
					object.VlanType = "VmwareDistributedVirtualSwitchTrunkVlanSpec"
					for i, v := range vlanspec.VlanId {
						if v.Start == v.End {
							object.Vlan += fmt.Sprint(v.Start)
						} else {
							object.Vlan += fmt.Sprintf("%v - %v", v.Start, v.End)
						}
						// add a comma, if needed
						if i != len(vlanspec.VlanId)-1 && object.Vlan != "" {
							object.Vlan += ", "
						}
					}

				default:
					// TODO: support for spec: *types.VmwareDistributedVirtualSwitchPvlanSpec
					object.VlanType = "TypeNotImplemented"
					object.Vlan = "unknown"
				}
			}
		}
		list = append(list, object)
	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return
}
