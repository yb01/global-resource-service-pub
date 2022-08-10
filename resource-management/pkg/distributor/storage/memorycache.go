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
