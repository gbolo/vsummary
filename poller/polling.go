package poller

// GetPollResults will return pollResults along with all errors encountered during the polling
func (p *Poller) GetPollResults() (r pollResults, errors []error) {

	var err error
	r.Vcenter, err = p.GetVcenter()
	if err != nil {
		// if we can't get vcenter info, we might as well just quit here...
		appendIfError(&errors, err)
		return
	}

	// if we got past the vcenter poll, we can do the rest now
	r.Esxi, _, r.VSwitch, r.StdPortgroup, err = p.GetEsxi()
	appendIfError(&errors, err)
	r.Virtualmachine, r.VDisk, r.Vnic, err = p.GetVirtualMachines()
	appendIfError(&errors, err)
	r.Datacenter, err = p.GetDatacenters()
	appendIfError(&errors, err)
	r.Cluster, err = p.GetClusters()
	appendIfError(&errors, err)
	r.Datastore, err = p.GetDatastores()
	appendIfError(&errors, err)
	r.Dvs, err = p.GetDVS()
	appendIfError(&errors, err)
	r.DvsPortGroup, err = p.GetDVSPortgroups()
	appendIfError(&errors, err)
	r.ResourcePool, err = p.GetResourcepools()
	appendIfError(&errors, err)
	r.Folder, err = p.GetFolders()
	appendIfError(&errors, err)

	return
}
