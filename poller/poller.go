package poller

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/op/go-logging"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

var log = logging.MustGetLogger("vsummary")

const apiVersion = "2"

type PollerConfig struct {
	URL         string
	UserName    string
	Password    string
	IntervalMin int

	// Don't validate TLS Cert
	Insecure bool
}

type Poller struct {
	VmwareClient *govmomi.Client
	Config       *PollerConfig
}

func NewPoller() *Poller {
	return &Poller{}
}

func (p *Poller) Configure(config *PollerConfig) {
	p.Config = config
}

func (p *Poller) Connect(ctx *context.Context) (err error) {

	if p.Config.URL == "" {
		err = errors.New("vsphere URL cannot be empty")
		return
	}

	vUrl, err := soap.ParseURL(p.Config.URL)
	if err != nil {
		return
	}

	vUrl.User = url.UserPassword(p.Config.UserName, p.Config.Password)

	p.VmwareClient, err = govmomi.NewClient(*ctx, vUrl, p.Config.Insecure)
	return

}

func (p *Poller) GetVMs() (vmList []common.Vm, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "poll")

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
	var vms []mo.VirtualMachine
	err = v.Retrieve(
		ctx,
		[]string{"VirtualMachine"},
		[]string{"summary", "config", "guest", "runtime", "parent", "resourcePool", "parentVApp"},
		&vms,
	)
	if err != nil {
		return
	}

	// Print summary per vm (see also: govc/vm/info.go)
	for _, vm := range vms {

		// create vm struct
		vmStruct := common.Vm{
			Name:                 vm.Summary.Config.Name,
			Moref:                vm.Summary.Vm.Value,
			VmxPath:              vm.Config.Files.VmPathName,
			Vcpu:                 vm.Config.Hardware.NumCPU,
			MemoryMb:             vm.Config.Hardware.MemoryMB,
			ConfigGuestOs:        vm.Config.GuestId,
			ConfigVersion:        vm.Config.Version,
			SmbiosUuid:           vm.Config.Firmware,
			InstanceUuid:         vm.Config.Uuid,
			ConfigChangeVersion:  vm.Config.ChangeVersion,
			GuestToolsVersion:    vm.Guest.ToolsVersion,
			GuestToolsRunning:    vm.Guest.ToolsRunningStatus,
			GuestHostname:        vm.Guest.HostName,
			GuestIp:              vm.Guest.IpAddress,
			GuestOs:              vm.Guest.GuestId,
			StatCpuUsage:         vm.Summary.QuickStats.OverallCpuUsage,
			StatHostMemoryUsage:  vm.Summary.QuickStats.HostMemoryUsage,
			StatGuestMemoryUsage: vm.Summary.QuickStats.GuestMemoryUsage,
			StatUptimeSec:        vm.Summary.QuickStats.UptimeSeconds,
			PowerState:           string(vm.Runtime.PowerState),
			EsxiMoref:            vm.Runtime.Host.Value,
			Template:             vm.Config.Template,
			VcenterId:            v.Client().ServiceContent.About.InstanceUuid,
		}

		// folder may not exist
		if vm.Parent != nil {
			vmStruct.FolderMoref = vm.Parent.Value
			//vmStruct.FolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", vmStruct.VcenterId, vmStruct.FolderMoref))
		}

		// vapps may not exist
		if vm.ParentVApp != nil {
			vmStruct.VappMoref = vm.ParentVApp.Value
			//vmStruct.VappId = common.GetMD5Hash(fmt.Sprintf("%s%s", vmStruct.VcenterId, vmStruct.VappMoref))
			vmStruct.FolderId = "vapp"
		} else {
			vmStruct.VappId = "none"
		}

		// resourcepool may not exist
		if vm.ResourcePool != nil {
			vmStruct.ResourcePoolMoref = vm.ResourcePool.Value
			//vmStruct.ResourcePoolId = common.GetMD5Hash(fmt.Sprintf("%s%s", vmStruct.VcenterId, vmStruct.ResourcePoolId))
		}

		// Fill in some required Ids
		//vmStruct.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", vmStruct.VcenterId, vmStruct.Moref))
		//vmStruct.EsxiId = common.GetMD5Hash(fmt.Sprintf("%s%s", vmStruct.VcenterId, vmStruct.EsxiMoref))

		vmList = append(vmList, vmStruct)

	}

	log.Infof("poller fetched summary of %d vms", len(vmList))
	return

}

func (p *Poller) SendVMs() (err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "request")

	// get Vms
	vms, err := p.GetVMs()
	if err != nil {
		log.Debugf("failed to retrieve VM list: %s", err)
		return
	}

	log.Infof("poller sending summary of %d vms", len(vms))

	jsonVms, err := json.Marshal(vms)
	if err != nil {
		log.Errorf("invalid json vm: %s", err)
		return
	}

	err = sendResults("/vm", jsonVms)
	if err != nil {
		return
	}

	return

}
