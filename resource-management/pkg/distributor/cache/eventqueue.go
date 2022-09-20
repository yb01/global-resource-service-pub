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

package cache

import (
	"errors"
	"k8s.io/klog/v2"
	"sync"

	"global-resource-service/resource-management/pkg/common-lib/metrics"
	"global-resource-service/resource-management/pkg/common-lib/types"
	objectcache "global-resource-service/resource-management/pkg/common-lib/types/cache"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
	"global-resource-service/resource-management/pkg/distributor/node"
)

type NodeEventQueue struct {
	// corresponding client id
	clientId  string
	watchChan chan *event.NodeEvent

	// used to lock enqueue operation during snapshot
	enqueueLock sync.RWMutex

	eventQueueByLoc map[location.Location]*objectcache.EventQueue
	locationLock    sync.RWMutex
}

func NewNodeEventQueue(clientId string) *NodeEventQueue {
	queue := &NodeEventQueue{
		clientId:        clientId,
		eventQueueByLoc: make(map[location.Location]*objectcache.EventQueue),
	}

	return queue
}

func (eq *NodeEventQueue) AcquireSnapshotRLock() {
	eq.enqueueLock.RLock()
}

func (eq *NodeEventQueue) ReleaseSnapshotRLock() {
	eq.enqueueLock.RUnlock()
}

func (eq *NodeEventQueue) EnqueueEvent(e *node.ManagedNodeEvent) {
	eq.enqueueLock.Lock()
	defer eq.enqueueLock.Unlock()
	if eq.watchChan != nil {
		go func() {
			eq.watchChan <- e.GetNodeEvent()
		}()
	}

	eq.locationLock.Lock()
	defer eq.locationLock.Unlock()
	queueByLoc, isOK := eq.eventQueueByLoc[*e.GetLocation()]
	if !isOK {
		queueByLoc = objectcache.NewEventQueue()
		eq.eventQueueByLoc[*e.GetLocation()] = queueByLoc
	}
	queueByLoc.EnqueueEvent(e)
}

func (eq *NodeEventQueue) Watch(rvs types.InternalResourceVersionMap, clientWatchChan chan runtime.Object, stopCh chan struct{}) error {
	if eq.watchChan != nil {
		return errors.New("currently only support one watcher per node event queue")
	}

	// get events already in queues
	events, err := eq.getAllEventsSinceResourceVersion(rvs)
	if err != nil {
		return err
	}

	eq.watchChan = make(chan *event.NodeEvent)
	// writing event to channel
	go func(downstreamCh chan runtime.Object, initEvents []runtime.Object, stopCh chan struct{}, upstreamCh chan *event.NodeEvent) {
		if downstreamCh == nil {
			return
		}
		// send init events
		for i := 0; i < len(initEvents); i++ {
			downstreamCh <- initEvents[i]
		}

		// continue to watch
		for {
			select {
			case <-stopCh:
				eq.watchChan = nil
				klog.V(3).Infof("Watch stopped due to client request")
				return
			case event, ok := <-upstreamCh:
				if !ok {
					break
				}
				klog.V(9).Infof("Sending event with node id %v", event.Node.Id)
				event.SetCheckpoint(metrics.Distributor_Sending)
				downstreamCh <- event
				event.SetCheckpoint(metrics.Distributor_Sent)
				klog.V(9).Infof("Event with node id %v sent", event.Node.Id)
			}
		}

	}(clientWatchChan, events, stopCh, eq.watchChan)

	return nil
}

func (eq *NodeEventQueue) getAllEventsSinceResourceVersion(rvs types.InternalResourceVersionMap) ([]runtime.Object, error) {
	locStartPostitions := make(map[location.Location]int)

	for loc, rv := range rvs {
		qByLoc, isOK := eq.eventQueueByLoc[loc]
		if isOK {
			startIndex, err := qByLoc.GetEventIndexSinceResourceVersion(rv)
			if err != nil {
				if err == types.Error_EndOfEventQueue {
					return nil, nil
				}
				return nil, err
			}
			locStartPostitions[loc] = startIndex
		}
	}

	nodeEvents := make([]runtime.Object, 0)
	for loc, qByLoc := range eq.eventQueueByLoc {
		startIndex, isOK := locStartPostitions[loc]
		var events []runtime.Object
		var err error
		if isOK {
			events, err = qByLoc.GetEventsFromIndex(startIndex)
		} else {
			events, err = qByLoc.GetEventsFromIndex(qByLoc.GetStartPos())
		}
		if err != nil {
			return nil, err
		}

		if len(events) > 0 {
			nodeEvents = append(nodeEvents, events...)
		}
	}

	return nodeEvents, nil
}
