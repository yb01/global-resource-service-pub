package event

import (
	"k8s.io/klog/v2"
	"sync"

	"global-resource-service/resource-management/pkg/common-lib/metrics"
)

type LatencyMetricsAllCheckpoints struct {
	Aggregator_Received  *metrics.LatencyMetrics
	Distributor_Received *metrics.LatencyMetrics
	Distributor_Sending  *metrics.LatencyMetrics
	Distributor_Sent     *metrics.LatencyMetrics
	Serializer_Encoded   *metrics.LatencyMetrics
	Serializer_Sent      *metrics.LatencyMetrics
}

var latencyMetricsAllCheckpoints *LatencyMetricsAllCheckpoints
var latencyMetricsLock sync.RWMutex

func init() {
	latencyMetricsAllCheckpoints = new(LatencyMetricsAllCheckpoints)
	latencyMetricsAllCheckpoints.Aggregator_Received = metrics.NewLatencyMetrics(string(metrics.Aggregator_Received))
	latencyMetricsAllCheckpoints.Distributor_Received = metrics.NewLatencyMetrics(string(metrics.Distributor_Received))
	latencyMetricsAllCheckpoints.Distributor_Sending = metrics.NewLatencyMetrics(string(metrics.Distributor_Sending))
	latencyMetricsAllCheckpoints.Distributor_Sent = metrics.NewLatencyMetrics(string(metrics.Distributor_Sent))
	latencyMetricsAllCheckpoints.Serializer_Encoded = metrics.NewLatencyMetrics(string(metrics.Serializer_Encoded))
	latencyMetricsAllCheckpoints.Serializer_Sent = metrics.NewLatencyMetrics(string(metrics.Serializer_Sent))
}

func AddLatencyMetricsAllCheckpoints(e *NodeEvent) {
	if e == nil {
		klog.Error("Nil event")
	}
	checkpointsPerEvent := e.GetCheckpoints()
	if checkpointsPerEvent == nil {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have checkpoint stamped", e.Type, e.Node.Id, e.Node.ResourceVersion)
	}
	lastUpdatedTime := e.Node.LastUpdatedTime

	agg_received_time, isOK1 := checkpointsPerEvent[metrics.Aggregator_Received]
	dis_received_time, isOK2 := checkpointsPerEvent[metrics.Distributor_Received]
	dis_sending_time, isOK3 := checkpointsPerEvent[metrics.Distributor_Sending]
	dis_sent_time, isOK4 := checkpointsPerEvent[metrics.Distributor_Sent]
	serializer_encoded_time, isOK5 := checkpointsPerEvent[metrics.Serializer_Encoded]
	serializer_sent_time, isOK6 := checkpointsPerEvent[metrics.Serializer_Sent]

	latencyMetricsLock.Lock()
	defer latencyMetricsLock.Unlock()
	if isOK1 {
		latencyMetricsAllCheckpoints.Aggregator_Received.AddLatencyMetrics(agg_received_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Aggregator_Received)
	}
	if isOK2 {
		latencyMetricsAllCheckpoints.Distributor_Received.AddLatencyMetrics(dis_received_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Distributor_Received)
	}
	if isOK3 {
		latencyMetricsAllCheckpoints.Distributor_Sending.AddLatencyMetrics(dis_sending_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Distributor_Sending)
	}
	if isOK4 {
		latencyMetricsAllCheckpoints.Distributor_Sent.AddLatencyMetrics(dis_sent_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Distributor_Sent)
	}
	if isOK5 {
		latencyMetricsAllCheckpoints.Serializer_Encoded.AddLatencyMetrics(serializer_encoded_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Serializer_Encoded)
	}
	if isOK6 {
		latencyMetricsAllCheckpoints.Serializer_Sent.AddLatencyMetrics(serializer_sent_time.Sub(lastUpdatedTime))
	} else {
		klog.Errorf("Event (%v, Id %s, RV %s) does not have %s stamped", e.Type, e.Node.Id, e.Node.ResourceVersion, metrics.Serializer_Sent)
	}
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
	klog.Infof(metrics_Message, metrics.Aggregator_Received, agg_received_summary.P50, agg_received_summary.P90, agg_received_summary.P99, agg_received_summary.TotalCount)
	klog.Infof(metrics_Message, metrics.Distributor_Received, dis_received_summary.P50, dis_received_summary.P90, dis_received_summary.P99, dis_received_summary.TotalCount)
	klog.Infof(metrics_Message, metrics.Distributor_Sending, dis_sending_summary.P50, dis_sending_summary.P90, dis_sending_summary.P99, dis_sending_summary.TotalCount)
	klog.Infof(metrics_Message, metrics.Distributor_Sent, dis_sent_summary.P50, dis_sent_summary.P90, dis_sent_summary.P99, dis_sent_summary.TotalCount)
	klog.Infof(metrics_Message, metrics.Serializer_Encoded, serializer_encoded_summary.P50, serializer_encoded_summary.P90, serializer_encoded_summary.P99, serializer_encoded_summary.TotalCount)
	klog.Infof(metrics_Message, metrics.Serializer_Sent, serializer_sent_summary.P50, serializer_sent_summary.P90, serializer_sent_summary.P99, serializer_sent_summary.TotalCount)
}
