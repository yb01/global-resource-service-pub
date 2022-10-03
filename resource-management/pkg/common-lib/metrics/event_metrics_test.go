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

package metrics

import (
	"github.com/google/uuid"
	"runtime"
	"strconv"
	"testing"
	"time"

	common_lib "global-resource-service/resource-management/pkg/common-lib"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	runtime2 "global-resource-service/resource-management/pkg/common-lib/types/runtime"
)

var defaultLocBeijing_RP1 = location.NewLocation(location.Beijing, location.ResourcePartition1)
var rvToGenerate = 0

func Test_PrintLatencyReport(t *testing.T) {
	ne := createNodeEvent()

	time.Sleep(100 * time.Millisecond)
	ne.SetCheckpoint(int(Aggregator_Received))
	ne.SetCheckpoint(int(Distributor_Received))
	ne.SetCheckpoint(int(Distributor_Sending))
	ne.SetCheckpoint(int(Distributor_Sent))
	ne.SetCheckpoint(int(Serializer_Encoded))
	ne.SetCheckpoint(int(Serializer_Sent))
	AddLatencyMetricsAllCheckpoints(ne)
	PrintLatencyReport()
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
		LastUpdatedTime: time.Now().UTC(),
	}
}

func Test_MemoryUsageOfLatencyReport(t *testing.T) {
	count := 1000000
	// Get memory usage for 1M node events
	common_lib.ResourceManagementMeasurement_Enabled = false
	nodes := make([]*runtime2.NodeEvent, count)
	for i := 0; i < count; i++ {
		nodes[i] = createNodeEvent()
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("Alloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)

	// Enable metrics
	common_lib.ResourceManagementMeasurement_Enabled = true
	for i := 0; i < count; i++ {
		nodes[i].SetCheckpoint(int(Aggregator_Received))
		nodes[i].SetCheckpoint(int(Distributor_Received))
		nodes[i].SetCheckpoint(int(Distributor_Sending))
		nodes[i].SetCheckpoint(int(Distributor_Sent))
		nodes[i].SetCheckpoint(int(Serializer_Encoded))
		nodes[i].SetCheckpoint(int(Serializer_Sent))
		AddLatencyMetricsAllCheckpoints(nodes[i])
	}
	PrintLatencyReport()

	runtime.ReadMemStats(&m)
	t.Logf("Alloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)
}

func createNodeEvent() *runtime2.NodeEvent {
	n := createRandomNode(rvToGenerate+1, defaultLocBeijing_RP1)
	ne := runtime2.NewNodeEvent(n, runtime2.Added)
	return ne
}
