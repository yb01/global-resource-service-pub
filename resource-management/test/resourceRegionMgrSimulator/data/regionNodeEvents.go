package data

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"

	simulatorTypes "global-resource-service/resource-management/test/resourceRegionMgrSimulator/types"
)

// The following varables are used to create Region Node Event List in multiply RPs of one region
//
// RegionNodeEventsList      - for initpull
// SnapshotNodeListEvents    - for subsequent pull and
//                                 post CRV to discard all old node events
var RegionNodeEventsList simulatorTypes.RegionNodeEvents

var SnapshotNodeListEvents simulatorTypes.RegionNodeEvents
var RegionId, RpNum, NodesPerRP int

// The constants are for repeatly generate new modified events
const at2thMin5k = 5000
const at5thMin25k = 25000
const at7thMin1k = 1000

// Initialize two events list
// RegionNodeEventsList - for initpull
//
func Init(regionId, rpNum, nodesPerRP int) {
	RegionNodeEventsList = generateAddedNodeEvents(regionId, rpNum, nodesPerRP)
	RegionId = regionId
	RpNum = rpNum
	NodesPerRP = nodesPerRP
}

// Generate region node update event changes to
// add them into RegionNodeEventsList
//

func MakeDataUpdate() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			// At 2th minute mark, generate 5k modified and added node events
			time.Sleep(2 * time.Minute)
			makeDataUpdate(at2thMin5k)
			klog.Info("At 2th minute mark, generating 5k modified and added node events is completed")

			// At 5th time mark, generate 25k modified node events
			time.Sleep(3 * time.Minute)
			makeDataUpdate(at5thMin25k)
			klog.Info("At 5th time mark, generating 25k modified node events is completed")

			// At 7th time mark, generate 1k modified events
			time.Sleep(2 * time.Minute)
			makeDataUpdate(at7thMin1k)
			klog.Info("At 7th time mark, generating 1k modified events is completed")
		}
	}()
}

///////////////////////////////////////////////
// The following functions are for handler.
//////////////////////////////////////////////

// Return region node added events in BATCH LENGTH from all RPs
// TO DO: paginate support
//
func GetRegionNodeAddedEvents(batchLength int) simulatorTypes.RegionNodeEvents {
	klog.Infof("Total (%v) Added events are to be pulled", len(RegionNodeEventsList))
	return RegionNodeEventsList

}

// Return region node modified events with CRVs in BATCH LENGTH from all RPs
// TO DO: paginate support
//
func GetRegionNodeModifiedEventsCRV(rvs types.ResourceVersionMap) simulatorTypes.RegionNodeEvents {
	SnapshotNodeListEvents = RegionNodeEventsList
	length := len(SnapshotNodeListEvents)
	var eventUpdate simulatorTypes.RegionNodeEvents

	for i := 0; i < length; i++ {
		region := SnapshotNodeListEvents[i].Node.GeoInfo.Region
		rp := SnapshotNodeListEvents[i].Node.GeoInfo.ResourcePartition
		loc := location.NewLocation(location.Region(region), location.ResourcePartition(rp))

		if SnapshotNodeListEvents[i].Node.GetResourceVersionInt64() > rvs[*loc] {
			eventUpdate = append(eventUpdate, SnapshotNodeListEvents[i])
		}
	}

	klog.Infof("Total (%v) Modified events are to be pulled", len(eventUpdate))
	return eventUpdate
}

////////////////////////////////////////
// The below are all helper functions
////////////////////////////////////////

// This function is used to initialize the region node added event
//
func generateAddedNodeEvents(regionId, rpNum, nodesPerRP int) simulatorTypes.RegionNodeEvents {
	regionName := location.Regions[regionId-1]
	cap := rpNum * nodesPerRP
	eventsAdd := make(simulatorTypes.RegionNodeEvents, cap)

	index := -1
	for j := 0; j < rpNum; j++ {
		rpName := location.ResourcePartitions[j]
		loc := location.NewLocation(regionName, rpName)

		// Initialize the resource version starting from 0 for each RP
		var rvToGenerateRPs = 0
		for i := 0; i < nodesPerRP; i++ {
			rvToGenerateRPs += 1
			index += 1

			node := createRandomNode(rvToGenerateRPs, loc)
			nodeEvent := event.NewNodeEvent(node, event.Added)
			eventsAdd[index] = nodeEvent
		}

	}
	return eventsAdd
}

// This function is used to create region node modified event by go routine
//
func makeDataUpdate(changesThreshold int) {
	// Calculate how many node changes per RP
	nodeChangesPerRP := changesThreshold / RpNum

	// Make data update for each RP
	for j := 0; j < RpNum; j++ {
		startIndex := j * NodesPerRP
		endIndex := (j + 1) * NodesPerRP

		// Search the nodes in the RP to get the highestRV
		var highestRVForRP uint64 = 0
		for k := startIndex; k < endIndex; k++ {
			currentResourceVersion := RegionNodeEventsList[k].Node.Copy().GetResourceVersionInt64()
			if highestRVForRP < currentResourceVersion {
				highestRVForRP = currentResourceVersion
			}
		}

		// Pick up 'nodeChangesPerRP' nodes and make changes and assign the node RV to highestRV + 1
		count := 0
		rvToGenerateRPs := highestRVForRP + 1
		for count < nodeChangesPerRP {
			// Randonly create data update per RP node events list
			i := int(rand.Intn(NodesPerRP)) + startIndex

			node := RegionNodeEventsList[i].Node.Copy()

			// special case: Consider 5000 changes per RP for 500 nodes per RP
			// each node has 10 changes within this cycle
			currentResourceVersion := RegionNodeEventsList[i].Node.Copy().GetResourceVersionInt64()
			if currentResourceVersion < rvToGenerateRPs {
				node.ResourceVersion = strconv.FormatUint(rvToGenerateRPs, 10)
			} else {
				node.ResourceVersion = strconv.FormatUint(currentResourceVersion+1, 10)
			}

			newEvent := event.NewNodeEvent(node, event.Modified)
			RegionNodeEventsList[i] = newEvent

			count++
		}
	}

	klog.Infof("Actually total (%v) new modified events are created.", changesThreshold)
}

// Create logical node with random UUID
//
func createRandomNode(rv int, loc *location.Location) *types.LogicalNode {
	id := uuid.New()

	return &types.LogicalNode{
		Id:              id.String(),
		ResourceVersion: strconv.Itoa(rv),
		GeoInfo: types.NodeGeoInfo{
			Region:            types.RegionName(loc.GetRegion()),
			ResourcePartition: types.ResourcePartitionName(loc.GetResourcePartition()),
			DataCenter:        types.DataCenterName(strconv.Itoa(int(rand.Intn(10000))) + "-DataCenter"),
			AvailabilityZone:  types.AvailabilityZoneName("AZ-" + strconv.Itoa(int(rand.Intn(6)))),
			FaultDomain:       types.FaultDomainName(id.String() + "-FaultDomain"),
		},
		Taints: types.NodeTaints{
			NoSchedule: false,
			NoExecute:  false,
		},
		SpecialHardwareTypes: types.NodeSpecialHardWareTypeInfo{
			HasGpu:  true,
			HasFPGA: true,
		},
		AllocatableResource: types.NodeResource{
			MilliCPU:         int64(rand.Intn(200) + 20),
			Memory:           int64(rand.Intn(2000)),
			EphemeralStorage: int64(rand.Intn(2000000)),
			AllowedPodNumber: int(rand.Intn(20000000)),
			ScalarResources: map[types.ResourceName]int64{
				"GPU":  int64(rand.Intn(200)),
				"FPGA": int64(rand.Intn(200)),
			},
		},
		Conditions:  111,
		Reserved:    false,
		MachineType: types.NodeMachineType(id.String() + "-highend"),
	}
}
