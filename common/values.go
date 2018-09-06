package common

import "fmt"

const (
	// this version number is used for api endpoints that other packages will consume
	ApiVersion     = "2"
	endpointFormat = "/api/v%s/%s"
)

var (
	// these variables are used by other packags to determine API endpoints
	EndpointVirtualMachine = getEndpoint("virtualmachine")
	EndpointESXi           = getEndpoint("esxi")
	EndpointDatacenter     = getEndpoint("datacenter")
	EndpointDatastore      = getEndpoint("datastore")
	EndpointVSwitch        = getEndpoint("vswitch")
	EndpointDVS            = EndpointVSwitch
	EndpointCluster        = getEndpoint("cluster")
	EndpointPortGroup      = getEndpoint("portgroup")
	EndpointVCenter        = getEndpoint("vcenter")
	EndpointResourcepool   = getEndpoint("resourcepool")
	EndpointVDisk          = getEndpoint("vdisk")
	EndpointVNIC           = getEndpoint("vnic")
	EndpointFolder         = getEndpoint("folder")
	EndpointPoller         = getEndpoint("poller")
)

// getEndpoint returns a properly formatted API endpoint
func getEndpoint(suffix string) string {
	return fmt.Sprintf(endpointFormat, ApiVersion, suffix)
}
