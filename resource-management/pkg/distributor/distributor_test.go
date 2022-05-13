package distributor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"

	"resource-management/pkg/common-lib/types"
	"resource-management/pkg/common-lib/types/event"
	"resource-management/pkg/common-lib/types/location"
	"resource-management/pkg/distributor/cache"
	"resource-management/pkg/distributor/storage"
)

var existedNodeId = make(map[uuid.UUID]bool)
var rvToGenerate = 0

var singleTestLock = sync.Mutex{}

var defaultLocBeijing_RP1 = location.NewLocation(location.Beijing, location.ResourcePartition1)

func setUp() *ResourceDistributor {
	singleTestLock.Lock()
	return GetResourceDistributor()
}

func tearDown(distributor *ResourceDistributor) {
	defer singleTestLock.Unlock()

	// flush node stores
	distributor.defaultNodeStore = createNodeStore()

	// flush nodeEventQueueMap
	distributor.nodeEventQueueMap = make(map[string]*cache.NodeEventQueue)

	// flush clientToStores map
	distributor.clientToStores = make(map[string][]*storage.VirtualNodeStore)
}

func TestDistributorInit(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	assert.NotNil(t, distributor, "Distributor cannot be nil")

	// check default virtual node stores
	defaultNodeStores := distributor.defaultNodeStore.GetVirtualStores()
	assert.Equal(t, true, len(*defaultNodeStores) > 500, "Expecting virtual store number >= 500")

	lower := float64(0)
	for i := 0; i < len(*defaultNodeStores); i++ {
		store := (*defaultNodeStores)[i]
		assert.Equal(t, 0, store.GetHostNum(), "Initial host number should be 0")
		assert.Equal(t, "", store.GetAssignedClient(), "Virtual store should not be assigned to any client")
		lowerBound, upperBound := store.GetRange()
		assert.Equal(t, lower, lowerBound, "Expecting lower bound %f but got %f. store id %d, hash range (%f, %f]", lower, lowerBound, i, lowerBound, upperBound)
		assert.NotEqual(t, lowerBound, upperBound, "Expecting lower bound not equal to upper bound for virtual store %d. Got hash range (%f, %f]", i, lowerBound, upperBound)
		lower = upperBound
		if i == len(*defaultNodeStores)-1 {
			assert.Equal(t, location.RingRange, upperBound, "Expecting last virtual store upper bound equals %f but got %f", location.RingRange, upperBound)
		}

		loc := store.GetLocation()
		assert.NotNil(t, loc, "Location of store should not be empty")
		if defaultLocBeijing_RP1.Equal(loc) {
			fmt.Printf("virtual node store %d, location %s, hash range (%f, %f]\n", i, store.GetLocation(), lowerBound, upperBound)
		}
	}
}

func measureProcessEvent(t *testing.T, dis *ResourceDistributor, eventType string, events []*event.NodeEvent, previousNodeCount int) {
	start := time.Now()
	result, rvMap := dis.ProcessEvents(events)
	duration := time.Since(start)
	fmt.Printf("Processing %d %s events took %v. Composite RVs %v\n", len(events), eventType, duration, rvMap)

	assert.True(t, result, "Expecting successfull event processing but got error")
	assert.NotNil(t, rvMap, "Expecting non nill rv map")
	assert.Equal(t, len(events)+previousNodeCount, dis.defaultNodeStore.GetTotalHostNum(), "Expected host number %d does not match actual host number %d", len(events), dis.defaultNodeStore.GetTotalHostNum())

	// iterate over virtual node stores
	hostCount := 0
	for _, vNodeStore := range *dis.defaultNodeStore.GetVirtualStores() {
		hostCount += vNodeStore.GetHostNum()
		assert.NotNil(t, vNodeStore.GetLocation())
	}
	assert.Equal(t, len(events)+previousNodeCount, hostCount, "Expected host number %d does not match actual host number %d", len(events), hostCount)
}

func TestAddNodes(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	nodeCounts := []int{10, 100, 1000, 10000, 100000, 1000000}
	previousNodeCount := 0
	for i := 0; i < len(nodeCounts); i++ {
		eventsAdd := generateAddNodeEvent(nodeCounts[i])
		measureProcessEvent(t, distributor, "AddNode", eventsAdd, previousNodeCount)
		previousNodeCount += nodeCounts[i]
	}
}

func generateAddNodeEvent(eventNum int) []*event.NodeEvent {
	result := make([]*event.NodeEvent, eventNum)
	for i := 0; i < eventNum; i++ {
		rvToGenerate += 1
		node := createRandomNode(rvToGenerate)
		nodeEvent := event.NewNodeEvent(node, event.Added)
		result[i] = nodeEvent
	}
	return result
}

func generateUpdateNodeEvents(originalEvents []*event.NodeEvent) []*event.NodeEvent {
	result := make([]*event.NodeEvent, len(originalEvents))
	for i := 0; i < len(originalEvents); i++ {
		rvToGenerate += 1

		newEvent := event.NewNodeEvent(types.NewNode(originalEvents[i].GetNode().GetId(), strconv.Itoa(rvToGenerate), "", originalEvents[i].GetNode().GetLocation()),
			event.Modified)
		result[i] = newEvent
	}
	return result
}

func generatedUpdateNodeEventsFromNodeList(nodes []*types.Node) []*event.NodeEvent {
	result := make([]*event.NodeEvent, len(nodes))
	for i := 0; i < len(nodes); i++ {
		rvToGenerate += 1
		newEvent := event.NewNodeEvent(types.NewNode(nodes[i].GetId(), strconv.Itoa(rvToGenerate), "", nodes[i].GetLocation()),
			event.Modified)
		result[i] = newEvent
	}
	return result
}

func createRandomNode(rv int) *types.Node {
	id := uuid.New()
	return types.NewNode(id.String(), strconv.Itoa(rv), "", defaultLocBeijing_RP1)
}

func TestUpdateNodes(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	nodeCounts := []int{10, 100, 1000, 10000, 100000, 1000000}
	previousNodeCount := 0
	for i := 0; i < len(nodeCounts); i++ {
		addAndUpdateNodes(t, distributor, nodeCounts[i], previousNodeCount)
		previousNodeCount += nodeCounts[i]
	}
}

func addAndUpdateNodes(t *testing.T, distributor *ResourceDistributor, eventNum int, previousNodeCount int) {
	eventsAdd := generateAddNodeEvent(eventNum)
	measureProcessEvent(t, distributor, "AddNode", eventsAdd, previousNodeCount)
	// update nodes
	eventsUpdate := generateUpdateNodeEvents(eventsAdd)
	measureProcessEvent(t, distributor, "UpdateNode", eventsUpdate, previousNodeCount)
}

func TestRegisterClient_ErrorCases(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	result, rvMap := distributor.ProcessEvents(generateAddNodeEvent(10))
	assert.True(t, result)
	assert.NotNil(t, rvMap)
	assert.Equal(t, 10, distributor.defaultNodeStore.GetTotalHostNum())

	// not enough hosts
	clientId, result, err := distributor.RegisterClient(100)
	assert.False(t, result, "Expecting request fail due to not enough hosts")
	assert.NotNil(t, clientId, "Expecting not nil client id")
	assert.False(t, clientId == "", "Expecting non empty client id")
	assert.Equal(t, types.Error_HostRequestExceedLimit, err)

	// less than minimal request host number
	clientId, result, err = distributor.RegisterClient(MinimalRequestHostNum - 1)
	assert.False(t, result, "Expecting request fail due to less than minimal host request")
	assert.NotNil(t, clientId, "Expecting not nil client id")
	assert.False(t, clientId == "", "Expecting non empty client id")
	assert.Equal(t, types.Error_HostRequestLessThanMiniaml, err)
}

func TestRegisterClient_WithinLimit(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	result, rvMap := distributor.ProcessEvents(generateAddNodeEvent(10000))
	assert.True(t, result)
	assert.NotNil(t, rvMap)
	assert.Equal(t, 10000, distributor.defaultNodeStore.GetTotalHostNum())

	requestedHostNum := 500
	for i := 0; i < 10; i++ {
		start := time.Now()
		clientId, result, err := distributor.RegisterClient(requestedHostNum)
		duration := time.Since(start)

		assert.True(t, result, "Expecting register client successfully")
		assert.NotNil(t, clientId, "Expecting not nil client id")
		assert.False(t, clientId == "", "Expecting non empty client id")
		assert.Nil(t, err, "Expecting nil error")

		// check virtual node assignment
		virtualStoresAssignedToClient, isOK := distributor.clientToStores[clientId]
		assert.True(t, isOK, "Expecting get virtual stores assigned to client %s", clientId)
		assert.True(t, len(virtualStoresAssignedToClient) > 0, "Expecting get non empty virtual stores assigned to client %s", clientId)
		hostCount := 0
		for i := 0; i < len(virtualStoresAssignedToClient); i++ {
			vs := virtualStoresAssignedToClient[i]
			assert.Equal(t, clientId, vs.GetAssignedClient(), "Unexpected virtual store client id %s", clientId)
			lower, upper := vs.GetRange()
			fmt.Printf("Virtual node store (%f, %f] is assigned to client %s, host number %d\n", lower, upper, clientId, vs.GetHostNum())
			hostCount += vs.GetHostNum()
		}
		fmt.Printf("Total %d hosts are assigned to client %s\nTook %v to register the client.\n", hostCount, clientId, duration)
		assert.True(t, hostCount >= requestedHostNum, "Assigned host number %d is less than requested %d", hostCount, requestedHostNum)
	}
}

func TestRegistrationWorkflow(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	// initialize node store with 10K nodes
	eventsAdd := generateAddNodeEvent(10000)
	result, rvMap := distributor.ProcessEvents(eventsAdd)
	assert.True(t, result)
	assert.NotNil(t, rvMap)
	assert.Equal(t, 10000, distributor.defaultNodeStore.GetTotalHostNum())

	// update nodes
	eventsUpdate := generateUpdateNodeEvents(eventsAdd)
	result, rvMap = distributor.ProcessEvents(eventsUpdate)
	assert.True(t, result)
	assert.NotNil(t, rvMap)
	assert.Equal(t, 10000, distributor.defaultNodeStore.GetTotalHostNum())

	// register client
	requestedHostNum := 500
	clientId, result, err := distributor.RegisterClient(requestedHostNum)
	assert.True(t, result, "Expecting register client successfully")
	assert.NotNil(t, clientId, "Expecting not nil client id")
	assert.False(t, clientId == "", "Expecting non empty client id")
	assert.Nil(t, err, "Expecting nil error")

	// client list nodes
	nodes, latestRVs, err := distributor.ListNodesForClient(clientId)
	assert.Nil(t, err)
	assert.NotNil(t, latestRVs)
	assert.True(t, len(nodes) >= 500)
	fmt.Printf("Latest rvs: %v\n", latestRVs)
	// check each node event
	nodeIds := make(map[string]bool)
	for _, node := range nodes {
		assert.NotNil(t, node.GetLocation())
		assert.True(t, latestRVs[*node.GetLocation()] >= node.GetResourceVersion())
		if _, isOK := nodeIds[node.GetId()]; isOK {
			assert.Fail(t, "List nodes cannot have more than ")
		} else {
			nodeIds[node.GetId()] = true
		}
	}
	assert.Equal(t, len(nodes), len(nodeIds))

	// update nodes
	oldNodeRV := nodes[0].GetResourceVersion()
	updateNodeEvents := generatedUpdateNodeEventsFromNodeList(nodes)
	result2, rvMap2 := distributor.ProcessEvents(updateNodeEvents)
	assert.True(t, result2, "Expecting update nodes successfully")
	loc := nodes[0].GetLocation()
	assert.Equal(t, uint64(rvToGenerate), rvMap2[*loc])
	assert.Equal(t, oldNodeRV, nodes[0].GetResourceVersion(), "Expecting listed nodes are snapshoted and cannot be affected by update")

	// client watch node update
	watchCh := make(chan *event.NodeEvent)
	stopCh := make(chan struct{})
	err = distributor.Watch(clientId, latestRVs, watchCh, stopCh)
	if err != nil {
		assert.Fail(t, "Encountered error while building watch connection.", "Encountered error while building watch connection. Error %v", err)
		return
	}
	watchedEventCount := 0
	for e := range watchCh {
		assert.Equal(t, event.Modified, e.Type)
		assert.Equal(t, loc, e.GetNode().GetLocation())
		watchedEventCount++

		if watchedEventCount >= len(nodes) {
			break
		}
	}
	assert.Equal(t, len(nodes), watchedEventCount)
	fmt.Printf("Latest rvs after updates: %v\n", rvMap2)
}
