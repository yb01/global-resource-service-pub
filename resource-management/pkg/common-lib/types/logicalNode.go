package types

import (
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
)

const (
	PreserveNode_KeyPrefix = "MinNode"
)

// for now, simply define those as string
// RegionName and ResourcePartitionName are updated to int per initial performance test of distributor ProcessEvents
// Later the data type might be changed back to string due to further performance evaluation result
type RegionName int
type ResourcePartitionName int
type DataCenterName string
type AvailabilityZoneName string
type FaultDomainName string

type NodeGeoInfo struct {
	// Region and RsourcePartition are required
	Region            RegionName
	ResourcePartition ResourcePartitionName

	// Optional fields for fine-tuned resource management and application placements
	DataCenter       DataCenterName
	AvailabilityZone AvailabilityZoneName
	FaultDomain      FaultDomainName
}

type NodeTaints struct {
	// Do not allow new pods to schedule onto the node unless they tolerate the taint,
	// Enforced by the scheduler.
	NoSchedule bool
	// Evict any already-running pods that do not tolerate the taint
	NoExecute bool
}

// TODO: consider refine for GPU types, such as NVIDIA and AMD etc.
type NodeSpecialHardWareTypeInfo struct {
	HasGpu  bool
	HasFPGA bool
}

// struct definition from Arktos node_info.go
type NodeResource struct {
	MilliCPU         int64
	Memory           int64
	EphemeralStorage int64
	// We store allowedPodNumber (which is Node.Status.Allocatable.Pods().Value())
	// explicitly as int, to avoid conversions and improve performance.
	AllowedPodNumber int
	// ScalarResources such as GPU or FPGA etc.
	ScalarResources map[ResourceName]int64
}

// TODO: from the Node definition in resource cluster, to the logicalNode struct, to the scheduler node_info structure
//       the ResourceName need to be set and aligned
type ResourceName string

// LogicalNode is the abstraction of the node definition in the resource clusters
// LogicalNode is a minimum set of information the scheduler needs to place the workload to a node in the region-less platform
//
// Initial set of fields as shown below.
//
// TODO: add the annotation for serialization
//
type LogicalNode struct {
	// Node UUID from each resource partition cluster
	Id string

	// ResourceVersion is the RV from each resource partition cluster
	ResourceVersion string

	// GeoInfo defines the node location info such as region, DC, RP cluster etc. for application placement
	GeoInfo NodeGeoInfo

	// Taints defines scheduling or other control action for a node
	Taints NodeTaints

	// SpecialHardwareTypes defines if the node has special hardware such as GPU or FPGA etc
	SpecialHardwareTypes NodeSpecialHardWareTypeInfo

	// AllocatableReesource defines the resources on the node that can be used by schedulers
	AllocatableResource NodeResource

	// Conditions is a short version of the node condition array from Arktos, each bits in the byte defines one node condition
	Conditions byte

	// Reserved defines if the node is reserved at the resource partition cluster level
	// TBD Node reservation model for post 630
	Reserved bool

	// MachineType defines the type of category of the node, such as # of CPUs of the node, where the category can be
	// defined as highend, lowend, medium as an example
	// TBD for post 630
	MachineType NodeMachineType
}

func (n *LogicalNode) Copy() *LogicalNode {
	return &LogicalNode{
		Id:                   n.Id,
		ResourceVersion:      n.ResourceVersion,
		GeoInfo:              n.GeoInfo,
		Taints:               n.Taints,
		SpecialHardwareTypes: n.SpecialHardwareTypes,
		AllocatableResource:  n.AllocatableResource,
		Conditions:           n.Conditions,
		Reserved:             n.Reserved,
		MachineType:          n.MachineType,
	}
}

func (n *LogicalNode) GetResourceVersionInt64() uint64 {
	rv, err := strconv.ParseUint(n.ResourceVersion, 10, 64)
	if err != nil {
		klog.Errorf("Unable to convert resource version %s to uint64\n", n.ResourceVersion)
		return 0
	}
	return rv
}

func (n *LogicalNode) GetKey() string {
	if n != nil {
		return fmt.Sprintf("%s.%s.%v.%v", PreserveNode_KeyPrefix, n.Id, n.GeoInfo.Region, n.GeoInfo.ResourcePartition)
	}
	return ""
}

type NodeMachineType string
