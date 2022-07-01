package cache

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	nodeutil "global-resource-service/resource-management/pkg/distributor/node"
	"strconv"
	"strings"
	"testing"
)

var rvToGenerate = 10
var defaultLocBeijing_RP1 = location.NewLocation(location.Beijing, location.ResourcePartition1)

func Test_getEventIndexSinceResourceVersion_ByLoc(t *testing.T) {
	// initalize node event queue by loc
	qloc := newNodeQueueByLoc()
	for i := 1; i <= 100; i++ {
		qloc.enqueueEvent(generateManagedNodeEvent(defaultLocBeijing_RP1))
	}
	t.Logf("Current rv %v", rvToGenerate)

	// search rv 1
	index, err := qloc.getEventIndexSinceResourceVersion(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, -1, index)

	// search rv 10
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(10))
	assert.NotNil(t, err)
	assert.Equal(t, -1, index)

	// search rv 11
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(11))
	assert.Nil(t, err)
	assert.Equal(t, 1, index)

	// search rv 99
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(99))
	assert.Nil(t, err)
	assert.Equal(t, 89, index)

	// search rv 100
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(100))
	assert.Nil(t, err)
	assert.Equal(t, 90, index)

	// search rv 109
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(109))
	assert.Nil(t, err)
	assert.Equal(t, 99, index)

	// search rv 110
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(110))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)

	// search rv 111
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(111))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)

	// generate event to exceed length of queue
	for i := 1; i <= LengthOfNodeEventQueue; i++ {
		qloc.enqueueEvent(generateManagedNodeEvent(defaultLocBeijing_RP1))
	}
	t.Logf("Current rv %v", rvToGenerate)
	t.Logf("Current start pos %v, end pos %v", qloc.startPos, qloc.endPos)

	// search rv 11
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(11))
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "newer than requested resource version 11"))
	assert.Equal(t, -1, index)

	// search rv 10100
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(10100))
	assert.Nil(t, err)
	assert.Equal(t, 10090, index)

	// search rv 10110
	index, err = qloc.getEventIndexSinceResourceVersion(uint64(10110))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)
}

func generateManagedNodeEvent(loc *location.Location) *nodeutil.ManagedNodeEvent {
	rvToGenerate += 1
	node := createRandomNode(rvToGenerate, loc)
	nodeEvent := event.NewNodeEvent(node, event.Added)
	return nodeutil.NewManagedNodeEvent(nodeEvent, loc)
}

func createRandomNode(rv int, loc *location.Location) *types.LogicalNode {
	id := uuid.New()
	return &types.LogicalNode{
		Id:              id.String(),
		ResourceVersion: strconv.Itoa(rv),
		GeoInfo: types.NodeGeoInfo{
			Region:            types.RegionName(loc.GetRegion()),
			ResourcePartition: types.ResourcePartitionName(loc.GetResourcePartition()),
		},
	}
}
