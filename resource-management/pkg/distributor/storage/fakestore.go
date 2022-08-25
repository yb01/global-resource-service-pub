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
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"sync"
	"time"

	"global-resource-service/resource-management/pkg/common-lib/interfaces/store"
	"global-resource-service/resource-management/pkg/common-lib/types"
)

const PersistDelayDefault = 20 * time.Nanosecond

type FakeStorageInterface struct {
	PersistDelayInNS  int
	NodeIds           map[string]bool
	nodeIdCacheLock   sync.RWMutex
	isTestNodeIdMatch bool
}

func (fs *FakeStorageInterface) InitNodeIdCache() {
	fs.NodeIds = make(map[string]bool)
}

func (fs *FakeStorageInterface) GetNodeIdCount() int {
	return len(fs.NodeIds)
}

func (fs *FakeStorageInterface) SetTestNodeIdMatch(isMatch bool) {
	fs.isTestNodeIdMatch = isMatch
}

func (fs *FakeStorageInterface) PersistNodes(nodesToPersist []*types.LogicalNode) bool {
	fs.simulateDelay(len(nodesToPersist))

	//klog.Infof("TestNodeIdMatch = %v", fs.isTestNodeIdMatch)
	if fs.isTestNodeIdMatch {
		fs.nodeIdCacheLock.Lock()
		for i := 0; i < len(nodesToPersist); i++ {
			fs.NodeIds[nodesToPersist[i].Id] = true
		}
		fs.nodeIdCacheLock.Unlock()
	}
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

func (fs *FakeStorageInterface) PersistClient(clientId string, client *types.Client) error {
	fs.simulateDelay(1)
	return nil
}

func (fs *FakeStorageInterface) GetClient(clientId string) (*types.Client, error) {
	return nil, errors.New("not implemented")
}

func (fs *FakeStorageInterface) UpdateClient(clientId string, client *types.Client) error {
	return fmt.Errorf("not implemented")
}

func (fs *FakeStorageInterface) GetClients() ([]*types.Client, error) {
	return nil, fmt.Errorf("not implemented")
}
