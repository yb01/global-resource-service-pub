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
	"fmt"
	"k8s.io/klog/v2"
	"sync"

	common_lib "global-resource-service/resource-management/pkg/common-lib"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
)

type LatencyMetricsAllCheckpoints struct {
	Aggregator_Received  *LatencyMetrics
	Distributor_Received *LatencyMetrics
	Distributor_Sending  *LatencyMetrics
	Distributor_Sent     *LatencyMetrics
	Serializer_Encoded   *LatencyMetrics
	Serializer_Sent      *LatencyMetrics
}

var latencyMetricsAllCheckpoints *LatencyMetricsAllCheckpoints
var latencyMetricsLock sync.RWMutex

func init() {
	latencyMetricsAllCheckpoints = new(LatencyMetricsAllCheckpoints)
	latencyMetricsAllCheckpoints.Aggregator_Received = NewLatencyMetrics(int(Aggregator_Received))
	latencyMetricsAllCheckpoints.Distributor_Received = NewLatencyMetrics(int(Distributor_Received))
	latencyMetricsAllCheckpoints.Distributor_Sending = NewLatencyMetrics(int(Distributor_Sending))
	latencyMetricsAllCheckpoints.Distributor_Sent = NewLatencyMetrics(int(Distributor_Sent))
	latencyMetricsAllCheckpoints.Serializer_Encoded = NewLatencyMetrics(int(Serializer_Encoded))
	latencyMetricsAllCheckpoints.Serializer_Sent = NewLatencyMetrics(int(Serializer_Sent))
}

func AddLatencyMetricsAllCheckpoints(e runtime.Object) {
	if !common_lib.ResourceManagementMeasurement_Enabled {
		return
	}
	if e == nil {
		klog.Error("Nil event")
	}
	checkpointsPerEvent := e.GetCheckpoints()
	if checkpointsPerEvent == nil {
		klog.Errorf("Event (%v, Id %s, RV %v) does not have checkpoint stamped", e.GetEventType(), e.GetId(), e.GetResourceVersionInt64())
	}
	lastUpdatedTime := e.GetLastUpdatedTime()

	agg_received_time := checkpointsPerEvent[Aggregator_Received]
	dis_received_time := checkpointsPerEvent[Distributor_Received]
	dis_sending_time := checkpointsPerEvent[Distributor_Sending]
	dis_sent_time := checkpointsPerEvent[Distributor_Sent]
	serializer_encoded_time := checkpointsPerEvent[Serializer_Encoded]
	serializer_sent_time := checkpointsPerEvent[Serializer_Sent]

	errMsg := fmt.Sprintf("Event (%v, Id %s, RV %v)", e.GetEventType(), e.GetId(), e.GetResourceVersionInt64()) + " does not have %s stamped"
	latencyMetricsLock.Lock()
	if !agg_received_time.IsZero() {
		latencyMetricsAllCheckpoints.Aggregator_Received.AddLatencyMetrics(agg_received_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Aggregator_Received)
	}
	if !dis_received_time.IsZero() {
		latencyMetricsAllCheckpoints.Distributor_Received.AddLatencyMetrics(dis_received_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Distributor_Received)
	}
	if !dis_sending_time.IsZero() {
		latencyMetricsAllCheckpoints.Distributor_Sending.AddLatencyMetrics(dis_sending_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Distributor_Sending)
	}
	if !dis_sent_time.IsZero() {
		latencyMetricsAllCheckpoints.Distributor_Sent.AddLatencyMetrics(dis_sent_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Distributor_Sent)
	}
	if !serializer_encoded_time.IsZero() {
		latencyMetricsAllCheckpoints.Serializer_Encoded.AddLatencyMetrics(serializer_encoded_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Serializer_Encoded)
	}
	if !serializer_sent_time.IsZero() {
		latencyMetricsAllCheckpoints.Serializer_Sent.AddLatencyMetrics(serializer_sent_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf(errMsg, Serializer_Sent)
	}
	latencyMetricsLock.Unlock()
	klog.V(6).Infof("[Metrics][Detail] node %v RV %v: %s: %v, %s: %v, %s: %v, %s: %v, %s: %v, %s: %v",
		e.GetId(), e.GetResourceVersionInt64(),
		Aggregator_Received_Name, agg_received_time.Sub(lastUpdatedTime),
		Distributor_Received_Name, dis_received_time.Sub(lastUpdatedTime),
		Distributor_Sending_Name, dis_sending_time.Sub(lastUpdatedTime),
		Distributor_Sent_Name, dis_sent_time.Sub(lastUpdatedTime),
		Serializer_Encoded_Name, serializer_encoded_time.Sub(lastUpdatedTime),
		Serializer_Sent_Name, serializer_sent_time.Sub(lastUpdatedTime))
}

func PrintLatencyReport() {
	latencyMetricsLock.RLock()
	agg_received_summary := latencyMetricsAllCheckpoints.Aggregator_Received.GetSummary()
	dis_received_summary := latencyMetricsAllCheckpoints.Distributor_Received.GetSummary()
	dis_sending_summary := latencyMetricsAllCheckpoints.Distributor_Sending.GetSummary()
	dis_sent_summary := latencyMetricsAllCheckpoints.Distributor_Sent.GetSummary()
	serializer_encoded_summary := latencyMetricsAllCheckpoints.Serializer_Encoded.GetSummary()
	serializer_sent_summary := latencyMetricsAllCheckpoints.Serializer_Sent.GetSummary()

	latencyMetricsLock.RUnlock()
	metrics_Message := "[Metrics][%s] perc50 %v, perc90 %v, perc99 %v. Total count %v"
	klog.Infof(metrics_Message, Aggregator_Received_Name, agg_received_summary.P50, agg_received_summary.P90, agg_received_summary.P99, agg_received_summary.TotalCount)
	klog.Infof(metrics_Message, Distributor_Received_Name, dis_received_summary.P50, dis_received_summary.P90, dis_received_summary.P99, dis_received_summary.TotalCount)
	klog.Infof(metrics_Message, Distributor_Sending_Name, dis_sending_summary.P50, dis_sending_summary.P90, dis_sending_summary.P99, dis_sending_summary.TotalCount)
	klog.Infof(metrics_Message, Distributor_Sent_Name, dis_sent_summary.P50, dis_sent_summary.P90, dis_sent_summary.P99, dis_sent_summary.TotalCount)
	klog.Infof(metrics_Message, Serializer_Encoded_Name, serializer_encoded_summary.P50, serializer_encoded_summary.P90, serializer_encoded_summary.P99, serializer_encoded_summary.TotalCount)
	klog.Infof(metrics_Message, Serializer_Sent_Name, serializer_sent_summary.P50, serializer_sent_summary.P90, serializer_sent_summary.P99, serializer_sent_summary.TotalCount)
}
