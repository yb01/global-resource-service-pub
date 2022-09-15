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

package stats

import (
	"time"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/metrics"
)

type NodeQueryStats struct {
	NodeQueryInterval time.Duration
	NumberOfNodes     int
	NodeQueryLatency  *metrics.LatencyMetrics
}

func NewNodeQueryStats() *NodeQueryStats {
	return &NodeQueryStats{NodeQueryLatency: metrics.NewLatencyMetrics(0)} // only one data point
}

func (nqs *NodeQueryStats) PrintStats() {
	latencySummary := nqs.NodeQueryLatency.GetSummary()
	klog.Infof("[Metrics][Nodes]QueryInterval: %v, Number of nodes queried during each interval: %v, perc50: %v, perc90: %v, perc99: %v. ",
		nqs.NodeQueryInterval, nqs.NumberOfNodes, latencySummary.P50, latencySummary.P90, latencySummary.P99)
}
