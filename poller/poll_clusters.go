package poller

import (
	"context"
	"time"

	//"github.com/gbolo/go-util/lib/debugging"
	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetClusters() (clList []common.Cluster, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollClusters")

	// Create view for objects
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var clusters []mo.ClusterComputeResource
	err = v.Retrieve(
		ctx,
		[]string{"ClusterComputeResource"},
		[]string{"name", "configuration.dasConfig", "configuration.drsConfig", "overallStatus", "parent", "summary"},
		&clusters,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, cluster := range clusters {

		// cluster.Summary cannot be indexed :(
		summary := cluster.Summary.GetComputeResourceSummary()

		clStruct := common.Cluster{
			Name:             cluster.Name,
			Moref:            cluster.Self.Value,
			VcenterId:        v.Client().ServiceContent.About.InstanceUuid,
			Status:           string(cluster.OverallStatus),
			DatacenterMoref:  cluster.Parent.Value,
			TotalCpuThreads:  summary.NumCpuThreads,
			TotalCpuMhz:      summary.TotalCpu,
			TotalMemoryBytes: summary.TotalMemory,
			TotalVmotions:    int32(common.GetInt(cluster, "Summary", "NumVmotions")),
			NumHosts:         summary.NumHosts,
			DRSBehaviour:     string(cluster.Configuration.DrsConfig.DefaultVmBehavior),
			CurrentBalance:   int32(common.GetInt(cluster, "Summary", "CurrentBalance")),
			TargetBalance:    int32(common.GetInt(cluster, "Summary", "TargetBalance")),

			// defaults
			DRSEnabled: "false",
			HAEnabled:  "false",
		}

		if cluster.Configuration.DrsConfig.Enabled != nil {
			clStruct.DRSEnabled = common.BoolToString(*cluster.Configuration.DrsConfig.Enabled)
		}

		if cluster.Configuration.DasConfig.Enabled != nil {
			clStruct.HAEnabled = common.BoolToString(*cluster.Configuration.DasConfig.Enabled)
		}

		clList = append(clList, clStruct)

	}

	log.Infof("poller fetched summary of %d cluster(s)", len(clList))
	return

}
