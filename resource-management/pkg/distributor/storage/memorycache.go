package storage

import (
	"sync"

	"global-resource-service/resource-management/pkg/common-lib/interfaces/store"
	"global-resource-service/resource-management/pkg/common-lib/types"
)

type DistributorPersistHelper struct {
	persistNodeWaitGroup *sync.WaitGroup

	persistHelper store.StoreInterface
}

func NewDistributorPersistHelper(persistHelper store.StoreInterface) *DistributorPersistHelper {
	return &DistributorPersistHelper{
		persistNodeWaitGroup: new(sync.WaitGroup),
		persistHelper:        persistHelper,
	}
}

func (c *DistributorPersistHelper) SetPersistHelper(persistTool store.StoreInterface) {
	c.persistHelper = persistTool
}

func (c *DistributorPersistHelper) SetWaitCount(count int) {
	c.persistNodeWaitGroup.Add(count)
}

func (c *DistributorPersistHelper) PersistNodes(newNodes []*types.LogicalNode) {
	go func(persistHelper store.StoreInterface, nodes []*types.LogicalNode, wg *sync.WaitGroup) {
		retries := 0
		defer func(numberOfNodes int, wg *sync.WaitGroup) {
			for i := 0; i < numberOfNodes; i++ {
				wg.Done()
			}
		}(len(nodes), wg)

		for {
			result := persistHelper.PersistNodes(nodes)
			if result {

				return
			} else {
				// TODO - error processing
				if retries >= 5 {
					return
				}
			}
			retries++
		}
	}(c.persistHelper, newNodes, c.persistNodeWaitGroup)
}

// TODO - timeout
func (c *DistributorPersistHelper) WaitForAllNodesSaved() {
	c.persistNodeWaitGroup.Wait()
}

func (c *DistributorPersistHelper) PersistStoreConfigs(nodeStoreStatus *store.NodeStoreStatus) bool {
	// persist virtual nodes location and latest resource version map
	resultPersistRVs := c.persistStoreStatus(nodeStoreStatus)

	return resultPersistRVs
}

func (c *DistributorPersistHelper) PersistVirtualNodesAssignment(assignment *store.VirtualNodeAssignment) bool {
	return c.persistHelper.PersistVirtualNodesAssignments(assignment)
}

func (c *DistributorPersistHelper) persistStoreStatus(nodeStoreStatus *store.NodeStoreStatus) bool {
	return c.persistHelper.PersistNodeStoreStatus(nodeStoreStatus)
}
