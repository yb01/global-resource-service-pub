package stats

import (
	"time"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/types"
)

const (
	LongWatchThreshold = time.Duration(1000 * time.Millisecond)
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
	klog.Infof("RegisterClientDuration: %v", rs.RegisterClientDuration)
}

type ListStats struct {
	ListDuration        time.Duration
	NumberOfNodesListed int
}

func NewListStats() *ListStats {
	return &ListStats{}
}

func (ls *ListStats) PrintStats() {
	klog.Infof("ListDuration: %v. Number of nodes listed: %v", ls.ListDuration, ls.NumberOfNodesListed)
}

type WatchStats struct {
	WatchDuration            time.Duration
	NumberOfProlongedItems   int
	NumberOfAddedNodes       int
	NumberOfUpdatedNodes     int
	NumberOfDeletedNodes     int
	NumberOfProlongedWatches int
}

func NewWatchStats() *WatchStats {
	return &WatchStats{}
}

func (ws *WatchStats) PrintStats() {
	klog.Infof("Watch session last: %v", ws.WatchDuration)
	klog.Infof("Number of nodes Added: %v", ws.NumberOfAddedNodes)
	klog.Infof("Number of nodes Updated: %v", ws.NumberOfUpdatedNodes)
	klog.Infof("Number of nodes Deleted: %v", ws.NumberOfDeletedNodes)
	klog.Infof("Number of nodes watch prolonged than %v: %v", LongWatchThreshold, ws.NumberOfProlongedItems)
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
