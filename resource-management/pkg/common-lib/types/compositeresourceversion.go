package types

import (
	"flag"
	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

// in case of UT or test uses the klog along with the t.log
func init() {
	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}
}

type CompositeResourceVersion struct {
	RegionId            string
	ResourcePartitionId string
	ResourceVersion     uint64
}

// Map from (regionId, ResourcePartitionId) to resourceVersion
type ResourceVersionMap map[location.Location]uint64

func (rvs *ResourceVersionMap) Copy() ResourceVersionMap {
	dupRVs := make(ResourceVersionMap, len(*rvs))
	for loc, rv := range *rvs {
		dupRVs[loc] = rv
	}

	return dupRVs
}
