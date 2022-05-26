package storage

import (
	"k8s.io/klog/v2"
	"time"

	"global-resource-service/resource-management/pkg/common-lib/interfaces/store"
	"global-resource-service/resource-management/pkg/common-lib/types"
)

const PersistDelayDefault = 20 * time.Nanosecond

type FakeStorageInterface struct {
	PersistDelayInNS int
}

func (fs *FakeStorageInterface) PersistNodes(nodesToPersist []*types.LogicalNode) bool {
	fs.simulateDelay(len(nodesToPersist))
	return true
}

func (fs *FakeStorageInterface) PersistNodeStoreStatus(nodeStoreStatus *store.NodeStoreStatus) bool {
	fs.simulateDelay(len(nodeStoreStatus.CurrentResourceVerions) + 3)
	return true
}

func (fs *FakeStorageInterface) PersistVirtualNodesAssignments(assignment *store.VirtualNodeAssignment) bool {
	fs.simulateDelay(len(assignment.VirtualNodes) + 1)
	return true
}

func (fs *FakeStorageInterface) simulateDelay(timesOfWrite int) {
	if fs.PersistDelayInNS > 0 {
		klog.V(3).Infof("Simulate disk persist operation delaying %v\n", time.Duration(timesOfWrite*fs.PersistDelayInNS)*time.Nanosecond)
		time.Sleep(time.Duration(timesOfWrite*fs.PersistDelayInNS) * time.Nanosecond)
	} else {
		klog.V(3).Infof("Simulate disk persist operation delaying %v\n", time.Duration(timesOfWrite)*PersistDelayDefault)
		time.Sleep(time.Duration(timesOfWrite) * PersistDelayDefault)
	}
}
