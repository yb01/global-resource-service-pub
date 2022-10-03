/*
Copyright 2022 Authors of Global Resource Service.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"strconv"
	"strings"
	"testing"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/cache"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
	nodeutil "global-resource-service/resource-management/pkg/distributor/node"
)

var rvToGenerate = 10
var defaultLocBeijing_RP1 = location.NewLocation(location.Beijing, location.ResourcePartition1)

func Test_getEventIndexSinceResourceVersion_ByLoc(t *testing.T) {
	// initalize node event queue by loc
	qloc := cache.NewEventQueue()
	for i := 1; i <= 100; i++ {
		qloc.EnqueueEvent(generateManagedNodeEvent(defaultLocBeijing_RP1))
	}
	t.Logf("Current rv %v", rvToGenerate)

	// search rv 1
	index, err := qloc.GetEventIndexSinceResourceVersion(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, -1, index)

	// search rv 10
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(10))
	assert.NotNil(t, err)
	assert.Equal(t, -1, index)

	// search rv 11
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(11))
	assert.Nil(t, err)
	assert.Equal(t, 1, index)

	// search rv 99
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(99))
	assert.Nil(t, err)
	assert.Equal(t, 89, index)

	// search rv 100
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(100))
	assert.Nil(t, err)
	assert.Equal(t, 90, index)

	// search rv 109
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(109))
	assert.Nil(t, err)
	assert.Equal(t, 99, index)

	// search rv 110
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(110))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)

	// search rv 111
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(111))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)

	// generate event to exceed length of queue
	for i := 1; i <= cache.LengthOfEventQueue; i++ {
		qloc.EnqueueEvent(generateManagedNodeEvent(defaultLocBeijing_RP1))
	}
	t.Logf("Current rv %v", rvToGenerate)
	t.Logf("Current start pos %v, end pos %v", qloc.GetStartPos(), qloc.GetEndPos())

	// search rv 11
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(11))
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "newer than requested resource version 11"))
	assert.Equal(t, -1, index)

	// search rv 10100
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(10100))
	assert.Nil(t, err)
	assert.Equal(t, 10090, index)

	// search rv 10110
	index, err = qloc.GetEventIndexSinceResourceVersion(uint64(10110))
	assert.NotNil(t, err)
	assert.Equal(t, types.Error_EndOfEventQueue, err)
	assert.Equal(t, -1, index)
}

func generateManagedNodeEvent(loc *location.Location) *nodeutil.ManagedNodeEvent {
	rvToGenerate += 1
	node := createRandomNode(rvToGenerate, loc)
	nodeEvent := runtime.NewNodeEvent(node, runtime.Added)
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
