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

package data

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

func TestGetRegionNodeModifiedEventsCRV(t *testing.T) {
	// create nodes
	rpNum := 10
	nodesPerRP := 50000
	start := time.Now()
	Init("Beijing", rpNum, nodesPerRP)
	// 2.827539846s
	duration := time.Since(start)
	assert.Equal(t, rpNum, len(RegionNodeEventsList))
	assert.Equal(t, nodesPerRP, len(RegionNodeEventsList[0]))
	t.Logf("Time to generate %d init events: %v", rpNum*nodesPerRP, duration)

	// List nodes
	start = time.Now()
	nodesEventList, count, rvs := ListNodes()
	duration = time.Since(start)
	assert.Equal(t, uint64(rpNum*nodesPerRP), count)
	assert.NotNil(t, nodesEventList)
	for i := 0; i < rpNum; i++ {
		loc := types.RvLocation{
			Region:    location.Beijing,
			Partition: location.ResourcePartition(i),
		}
		currentRV := rvs[loc]
		assert.Equal(t, uint64(nodesPerRP), currentRV)
	}

	// 500K nodes, list duration 179.605574ms
	t.Logf("List %v nodes, return RVS %v. duration %v", count, rvs, duration)

	// Watch node events
	watchCh := make(chan *event.NodeEvent)
	stopCh := make(chan struct{})
	err := Watch(rvs, watchCh, stopCh)
	if err != nil {
		assert.Fail(t, "Encountered error while building watch connection.", "Encountered error while building watch connection. Error %v", err)
		return
	}
	allWaitGroup := new(sync.WaitGroup)
	allWaitGroup.Add(1)
	updateEventCount := atEachMin10

	go func(expectedEventCount int, rvs types.TransitResourceVersionMap, watchCh chan *event.NodeEvent, wg *sync.WaitGroup) {
		eventCount := 0

		for e := range watchCh {
			assert.Equal(t, event.Modified, e.Type)
			loc := types.RvLocation{
				Region:    location.Beijing,
				Partition: location.ResourcePartition(e.Node.GeoInfo.ResourcePartition),
			}
			requestedRVForRP := rvs[loc]
			assert.True(t, requestedRVForRP < e.Node.GetResourceVersionInt64())

			eventCount++

			if eventCount >= expectedEventCount {
				wg.Done()
				close(watchCh)
				close(stopCh)
				return
			}
		}
	}(updateEventCount, rvs, watchCh, allWaitGroup)

	// generate update node events
	makeDataUpdate(atEachMin10)
	allWaitGroup.Wait()
	t.Logf("Watch %d events succeed!\n", updateEventCount)

	// watch from previous resource versions again
	watchCh = make(chan *event.NodeEvent)
	stopCh = make(chan struct{})
	err = Watch(rvs, watchCh, stopCh)
	if err != nil {
		assert.Fail(t, "Encountered error while building watch connection.", "Encountered error while building watch connection. Error %v", err)
		return
	}

	allWaitGroup.Add(1)
	go func(expectedEventCount int, rvs types.TransitResourceVersionMap, watchCh chan *event.NodeEvent, wg *sync.WaitGroup) {
		eventCount := 0

		for e := range watchCh {
			assert.Equal(t, event.Modified, e.Type)
			loc := types.RvLocation{
				Region:    location.Beijing,
				Partition: location.ResourcePartition(e.Node.GeoInfo.ResourcePartition),
			}
			requestedRVForRP := rvs[loc]
			assert.True(t, requestedRVForRP < e.Node.GetResourceVersionInt64())

			eventCount++

			if eventCount >= expectedEventCount {
				wg.Done()
				return
			}
		}
	}(updateEventCount, rvs, watchCh, allWaitGroup)

	allWaitGroup.Wait()
	t.Logf("Re-watch %d events succeed!\n", updateEventCount)
}
