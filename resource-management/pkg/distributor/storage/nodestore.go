package storage

import (
	"fmt"
	"math"
	"sync"

	"resource-management/pkg/common-lib/hash"
	"resource-management/pkg/common-lib/types"
	"resource-management/pkg/common-lib/types/event"
	"resource-management/pkg/common-lib/types/location"
	"resource-management/pkg/distributor/cache"
)

const (
	VirtualStoreInitSize = 100
)

type VirtualNodeStore struct {
	mu         sync.RWMutex
	nodeByHash map[float64]*types.Node
	lowerbound float64
	upperbound float64

	// one virtual store can only have nodes from one resource partition
	location location.Location

	clientId   string
	eventQueue *cache.NodeEventQueue
}

func (vs *VirtualNodeStore) GetHostNum() int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return len(vs.nodeByHash)
}

func (vs *VirtualNodeStore) GetLocation() location.Location {
	return vs.location
}

func (vs *VirtualNodeStore) GetAssignedClient() string {
	return vs.clientId
}

func (vs *VirtualNodeStore) AssignToClient(clientId string, eventQueue *cache.NodeEventQueue) bool {
	if vs.clientId != "" {
		return false
	} else if clientId == "" {
		return false
	} else if eventQueue == nil {
		return false
	}
	vs.clientId = clientId
	vs.eventQueue = eventQueue

	return true
}

func (vs *VirtualNodeStore) Release() {
	vs.clientId = ""
}

func (vs *VirtualNodeStore) GetRange() (float64, float64) {
	return vs.lowerbound, vs.upperbound
}

func (vs *VirtualNodeStore) SnapShot() ([]*types.Node, types.ResourceVersionMap) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	nodesCopy := make([]*types.Node, len(vs.nodeByHash))
	index := 0
	rvs := make(types.ResourceVersionMap)
	for _, node := range vs.nodeByHash {
		nodesCopy[index] = node.Copy()
		newRV := node.GetResourceVersion()
		if lastRV, isOK := rvs[*node.GetLocation()]; isOK {
			if lastRV < newRV {
				rvs[*node.GetLocation()] = newRV
			}
		} else {
			rvs[*node.GetLocation()] = newRV
		}
		index++
	}

	return nodesCopy, rvs
}

func (vs *VirtualNodeStore) GetOneNode() *types.Node {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	for _, n := range vs.nodeByHash {
		return n.Copy()
	}
	return nil
}

type NodeStore struct {
	// granularity of the ring - degree for each virtual node managed arc
	granularOfRing float64

	// # of regions
	regionNum int

	// # of max resource partition in each region
	partitionMaxNum int

	// # of different resource slots - computation various
	resourceSlots int

	virtualNodeNum int
	// node stores by virtual nodes
	// Using map instead of array to avoid expanding cost
	vNodeStores *[]*VirtualNodeStore
	// mutex for virtual store number/size adjustment
	nsLock sync.RWMutex

	totalHostNum int
	hostNumLock  sync.RWMutex

	// Latest resource version map
	currentRVs types.ResourceVersionMap
	rvLock     sync.RWMutex
}

func NewNodeStore(vNodeNumPerRP int, regionNum int, partitionMaxNum int) *NodeStore {
	totalVirtualNodeNum := vNodeNumPerRP * regionNum * partitionMaxNum
	virtualNodeStores := make([]*VirtualNodeStore, totalVirtualNodeNum)

	ns := &NodeStore{
		virtualNodeNum:  totalVirtualNodeNum,
		vNodeStores:     &virtualNodeStores,
		granularOfRing:  location.RingRange / (float64(totalVirtualNodeNum)),
		regionNum:       regionNum,
		partitionMaxNum: partitionMaxNum,
		resourceSlots:   regionNum * partitionMaxNum,
		currentRVs:      make(types.ResourceVersionMap),
		totalHostNum:    0,
	}

	ns.generateVirtualNodeStores(vNodeNumPerRP)
	return ns
}

// TODO - verify whether the original value can be changed. If so, return a deepcopy
func (ns *NodeStore) GetCurrentResourceVersions() types.ResourceVersionMap {
	ns.rvLock.RLock()
	defer ns.rvLock.RUnlock()
	return ns.currentRVs.Copy()
}

func (ns *NodeStore) GetTotalHostNum() int {
	ns.hostNumLock.RLock()
	defer ns.hostNumLock.RUnlock()
	return ns.totalHostNum
}

func (ns *NodeStore) CheckFreeCapacity(requestedHostNum int) bool {
	ns.nsLock.Lock()
	defer ns.nsLock.Unlock()
	allocatableHostNum := 0
	for _, vs := range *ns.vNodeStores {
		allocatableHostNum += vs.GetHostNum()
		if allocatableHostNum >= requestedHostNum {
			return true
		}
	}

	return false
}

func (ns *NodeStore) GetVirtualStores() *[]*VirtualNodeStore {
	return ns.vNodeStores
}

func (ns *NodeStore) generateVirtualNodeStores(vNodeNumPerRP int) {
	ns.nsLock.Lock()
	defer ns.nsLock.Unlock()

	vNodeIndex := 0
	for k := 0; k < ns.regionNum; k++ {
		region := location.Regions[k]
		rpsInRegion := location.GetRPsForRegion(region)

		for m := 0; m < ns.partitionMaxNum; m++ {
			loc := location.NewLocation(region, rpsInRegion[m])
			lowerBound, upperBound := loc.GetArcRangeFromLocation()

			for i := 0; i < vNodeNumPerRP; i++ {

				(*ns.vNodeStores)[vNodeIndex] = &VirtualNodeStore{
					mu:         sync.RWMutex{},
					nodeByHash: make(map[float64]*types.Node, VirtualStoreInitSize),
					lowerbound: lowerBound,
					upperbound: lowerBound + ns.granularOfRing,
					location:   *loc,
				}
				lowerBound += ns.granularOfRing
				vNodeIndex++
			}

			// remove the impact of inaccuracy
			(*ns.vNodeStores)[vNodeIndex-1].upperbound = upperBound
		}
	}

	(*ns.vNodeStores)[ns.virtualNodeNum-1].upperbound = location.RingRange
}

func (ns *NodeStore) CreateNode(nodeEvent *event.NodeEvent) {
	hashValue, ringId := ns.getNodeHash(nodeEvent.GetNode())
	isNewNode := ns.addNodeToRing(hashValue, ringId, nodeEvent)
	if !isNewNode {
		ns.updateNodeInRing(hashValue, ringId, nodeEvent)
	}
}

func (ns *NodeStore) UpdateNode(nodeEvent *event.NodeEvent) {
	hashValue, ringId := ns.getNodeHash(nodeEvent.GetNode())
	ns.updateNodeInRing(hashValue, ringId, nodeEvent)
}

// TODO
func (ns NodeStore) DeleteNode(nodeEvent event.NodeEvent) {
}

func (ns *NodeStore) ProcessNodeEvents(nodeEvents []*event.NodeEvent) (bool, types.ResourceVersionMap) {
	for _, e := range nodeEvents {
		ns.processNodeEvent(e)
	}

	// persist disk

	// TODO - make a copy of currentRVs in case modification happen unexpectedly
	return true, ns.currentRVs
}

func (ns *NodeStore) processNodeEvent(nodeEvent *event.NodeEvent) bool {
	switch nodeEvent.Type {
	case event.Added:
		ns.CreateNode(nodeEvent)
	case event.Modified:
		ns.UpdateNode(nodeEvent)
	default:
		return false
	}

	// Update ResourceVersionMap
	newRV := nodeEvent.GetNode().GetResourceVersion()
	ns.rvLock.Lock()
	if lastRV, isOK := ns.currentRVs[*nodeEvent.GetNode().GetLocation()]; isOK {
		if lastRV < newRV {
			ns.currentRVs[*nodeEvent.GetNode().GetLocation()] = newRV
		}
	} else {
		ns.currentRVs[*nodeEvent.GetNode().GetLocation()] = newRV
	}
	ns.rvLock.Unlock()

	return true
}

// return location on the ring, and ring Id
// ring Id is reserved for multiple rings
func (ns *NodeStore) getNodeHash(node *types.Node) (float64, int) {
	// map node id to uint32
	initHashValue := hash.HashStrToUInt64(node.GetId())

	// map node id to hash ring: (0 - 1]
	var ringValue float64
	if initHashValue == 0 {
		ringValue = 1
	} else {
		ringValue = float64(initHashValue) / float64(math.MaxUint64)
	}

	// compact to ring slice where this location belongs to
	lower, upper := node.GetLocation().GetArcRangeFromLocation()

	// compact ringValue onto (lower, upper]
	return lower + ringValue*(upper-lower), 0
}

func (ns *NodeStore) addNodeToRing(hashValue float64, ringId int, nodeEvent *event.NodeEvent) (isNewNode bool) {
	virtualNodeIndex := int(math.Floor(hashValue / ns.granularOfRing))
	node := nodeEvent.GetNode()
	vNodeStore := (*ns.vNodeStores)[virtualNodeIndex]
	vNodeStore.mu.Lock()
	defer vNodeStore.mu.Unlock()

	if oldNode, isOK := vNodeStore.nodeByHash[hashValue]; isOK {
		if oldNode.GetId() != node.GetId() {
			fmt.Printf("Found existing node (uuid %s) with same hash value %f. New node (uuid %s)\n", oldNode.GetId(), hashValue, node.GetId())
			// TODO - put node into linked list
		} else {
			return false
		}
	}
	vNodeStore.nodeByHash[hashValue] = node

	ns.hostNumLock.Lock()
	ns.totalHostNum++
	ns.hostNumLock.Unlock()

	// add event to event queue
	if vNodeStore.eventQueue != nil {
		vNodeStore.eventQueue.EnqueueEvent(nodeEvent)
	}

	return true
}

func (ns *NodeStore) updateNodeInRing(hashValue float64, ringId int, nodeEvent *event.NodeEvent) {
	virtualNodeIndex := int(math.Floor(hashValue / ns.granularOfRing))
	node := nodeEvent.GetNode()
	vNodeStore := (*ns.vNodeStores)[virtualNodeIndex]
	vNodeStore.mu.Lock()
	if oldNode, isOK := vNodeStore.nodeByHash[hashValue]; isOK {
		// TODO - check uuid to make sure updating right node
		if oldNode.GetId() == node.GetId() {
			if oldNode.GetResourceVersion() < node.GetResourceVersion() {
				vNodeStore.nodeByHash[hashValue] = node
			} else {
				fmt.Printf("Discard node update events due to resource version is older: %d. Existing rv %d", node.GetResourceVersion(), oldNode.GetResourceVersion())
				vNodeStore.mu.Unlock()
				return
			}
		} else {
			// TODO - check linked list to get right
			fmt.Printf("Updating node got same hash value (%f) but different node id: (%s and %s)", hashValue,
				oldNode.GetId(), node.GetId())
		}

		// add event to event queue
		if vNodeStore.eventQueue != nil {
			vNodeStore.eventQueue.EnqueueEvent(nodeEvent)
		}

		vNodeStore.mu.Unlock()
	} else {
		// ?? - report error or not?
		vNodeStore.mu.Unlock()
		ns.addNodeToRing(hashValue, ringId, nodeEvent)
	}
}
