package types

const (
	// Max and Min request per client request during registration or update resources
	MaxTotalMachinesPerRequest = 25000
	MinTotalMachinesPerRequest = 1000
)

// Client is detailed client info for a registered client to the resource management service
type Client struct {
	ClientId   string         `json:"client_id"`
	ClientInfo ClientInfoType `json:"client_info,omitempty"`
	// total granted resources for the client
	// it is initially set when client register to the service
	// it is updated on client to add/reduce resource via the {url}/resources/clientId/addResource(reduceResource) APIs
	Resource ResourceRequest `json:"client_granted_resource,omitempty"`
}

type ClientInfoType struct {
	// friendly name if client provided during registration
	ClientName string `json:"client_name,omitempty"`
	// which region the client is from, if client provided during registration
	Region string `json:"client_region,omitempty"`
}

// ResourcePerRegion is resource request for each region
// sum of the ResourcePerRegion is the total request for a given client
// post 630, revisit Resource request, and likely allow client to request in term of number of CPUs etc.
type ResourcePerRegion struct {
	// Name of the region
	RegionName string `json:"region_name"`

	// Machines requested per host machine type; machine type defined as CPU type etc.
	// flavors
	Machines map[NodeMachineType]int `json:"request_machines,omitempty"`

	// Machines requested per special hardware type, e.g., GPU / FPGA machines
	SpecialHardwareMachines map[string]int `json:"request_special_hardware_machines,omitempty"`
}

// ResourceRequest is used in the http request body for client to List resources
// default to request all region, 10K nodes total, no special hardware
// for 630, request with default only. machine flavors, special-hardware request will be supported post 630
type ResourceRequest struct {
	TotalMachines    int                 `json:"total_machines" default:"10000"`
	RequestInRegions []ResourcePerRegion `json:"resource_request,omitempty"`
}
