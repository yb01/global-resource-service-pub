package stats

import (
	"sync"
	"time"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/metrics"
	"global-resource-service/resource-management/pkg/common-lib/types"
)

const (
	LongWatchThreshold = time.Second
	RegionMapInitCap   = 10
	RpMapInitCap       = 20
)

type RegisterClientStats struct {
	RegisterClientDuration time.Duration
}

func NewRegisterClientStats() *RegisterClientStats {
	return &RegisterClientStats{}
}

func (rs *RegisterClientStats) PrintStats() {
	klog.Infof("[Metrics][Register]RegisterClientDuration: %v", rs.RegisterClientDuration)
}

type ListStats struct {
	ListDuration        time.Duration
	NumberOfNodesListed int
}

func NewListStats() *ListStats {
	return &ListStats{}
}

func (ls *ListStats) PrintStats() {
	klog.Infof("[Metrics][List]ListDuration: %v. Number of nodes listed: %v", ls.ListDuration, ls.NumberOfNodesListed)
}

type WatchStats struct {
	WatchDuration            time.Duration
	NumberOfProlongedItems   int
	NumberOfAddedNodes       int
	NumberOfUpdatedNodes     int
	NumberOfDeletedNodes     int
	NumberOfProlongedWatches int
	WatchDelayPerEvent       *metrics.LatencyMetrics
	WatchDelayLock           sync.RWMutex
}

func NewWatchStats() *WatchStats {
	return &WatchStats{WatchDelayPerEvent: metrics.NewLatencyMetrics(0)} // only one data point
}

func (ws *WatchStats) PrintStats() {
	klog.Infof("[Metrics][Watch]Watch session last: %v. Number of nodes Added :%v, Updated: %v, Deleted: %v. watch prolonged than %v: %v",
		ws.WatchDuration, ws.NumberOfAddedNodes, ws.NumberOfUpdatedNodes, ws.NumberOfDeletedNodes, LongWatchThreshold, ws.NumberOfProlongedItems)
	ws.WatchDelayLock.RLock()
	watchDelaySummary := ws.WatchDelayPerEvent.GetSummary()
	ws.WatchDelayLock.RUnlock()
	klog.Infof("[Metrics][Watch] perc50 %v, perc90 %v, perc99 %v. Total count %v",
		watchDelaySummary.P50, watchDelaySummary.P90, watchDelaySummary.P99, watchDelaySummary.TotalCount)
}

type groupByRp map[types.ResourcePartitionName]int

// GroupByRegionByRP groups nodes by region, RP for a given list of nodes
func GroupByRegionByRP(nodes []*types.LogicalNode) {
	groupByRegion := make(map[types.RegionName]groupByRp, 20)

	for _, n := range nodes {
		if k, found := groupByRegion[n.GeoInfo.Region]; found {
			if _, found2 := k[n.GeoInfo.ResourcePartition]; found2 {
				k[n.GeoInfo.ResourcePartition]++
			} else {
				k[n.GeoInfo.ResourcePartition] = 1
			}
		} else {
			if groupByRegion[n.GeoInfo.Region] == nil {
				groupByRegion[n.GeoInfo.Region] = make(groupByRp, RpMapInitCap)
			}
			groupByRegion[n.GeoInfo.Region][n.GeoInfo.ResourcePartition] = 1
		}
	}

	klog.Infof("Nodes group by region, by resource partition:")
	for region, rp := range groupByRegion {
		klog.Infof("Region: %v", region)
		for rpName, counts := range rp {
			klog.Infof("\t RP: %v, node counts: %v", rpName, counts)
		}
	}

	return
}

func GroupByRegion(nodes []*types.LogicalNode) {
	groupByRegion := make(map[types.RegionName]int, RegionMapInitCap)

	for _, n := range nodes {
		if _, found := groupByRegion[n.GeoInfo.Region]; found {
			groupByRegion[n.GeoInfo.Region]++
		} else {
			groupByRegion[n.GeoInfo.Region] = 1
		}
	}

	klog.Infof("Nodes group by region:")
	for region, counts := range groupByRegion {
		klog.Infof("\t Region: %v, node counts: %v", region, counts)
	}

	return
}
