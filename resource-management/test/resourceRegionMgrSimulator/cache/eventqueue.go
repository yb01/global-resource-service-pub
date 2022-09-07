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
	"fmt"
	"k8s.io/klog/v2"
	"sort"
	"sync"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

const LengthOfNodeEventQueue = 10000

type nodeEventQueueByLoc struct {
	circularEventQueue []*event.NodeEvent
	// circular event queue start position and end position
	startPos int
	endPos   int

	// mutex for event queue operation
	eqLock sync.RWMutex
}

func newNodeQueueByLoc() *nodeEventQueueByLoc {
	return &nodeEventQueueByLoc{
		circularEventQueue: make([]*event.NodeEvent, LengthOfNodeEventQueue),
		startPos:           0,
		endPos:             0,
		eqLock:             sync.RWMutex{},
	}
}

func (qloc *nodeEventQueueByLoc) enqueueEvent(e *event.NodeEvent) {
	qloc.eqLock.Lock()
	defer qloc.eqLock.Unlock()

	if qloc.endPos == qloc.startPos+LengthOfNodeEventQueue {
		// cache is full - remove the oldest element
		qloc.startPos++
	}

	qloc.circularEventQueue[qloc.endPos%LengthOfNodeEventQueue] = e
	qloc.endPos++
}

func (qloc *nodeEventQueueByLoc) getEventsFromIndex(startIndex int) ([]*event.NodeEvent, error) {
	qloc.eqLock.RLock()
	defer qloc.eqLock.RUnlock()

	if qloc.startPos == qloc.endPos || qloc.startPos > startIndex || startIndex > qloc.endPos { // queue is empty or out of range
		return nil, errors.New(fmt.Sprintf("Event queue start pos %d, end pos %d, invalid start index %d", qloc.startPos, qloc.endPos, startIndex))
	}

	length := qloc.endPos - startIndex
	result := make([]*event.NodeEvent, length)
	for i := 0; i < length; i++ {
		result[i] = qloc.circularEventQueue[(startIndex+i)%LengthOfNodeEventQueue]
	}

	return result, nil
}

func (qloc *nodeEventQueueByLoc) getEventIndexSinceResourceVersion(resourceVersion uint64) (int, error) {
	qloc.eqLock.RLock()
	defer qloc.eqLock.RUnlock()
	if qloc.endPos-qloc.startPos == 0 {
		return -1, errors.New(fmt.Sprintf("Empty event queue"))
	}
	nodeEvent := qloc.circularEventQueue[qloc.startPos%LengthOfNodeEventQueue].Node

	oldestRV := nodeEvent.GetResourceVersionInt64()
	if oldestRV > resourceVersion {
		return -1, errors.New(fmt.Sprintf("Loc %s events oldest resource Version %d is newer than requested resource version %d",
			nodeEvent.GeoInfo.Region, nodeEvent.GeoInfo.ResourcePartition, oldestRV, resourceVersion))
	}

	index := sort.Search(qloc.endPos-qloc.startPos, func(i int) bool {
		return qloc.circularEventQueue[(qloc.startPos+i)%LengthOfNodeEventQueue].Node.GetResourceVersionInt64() > resourceVersion
	})
	index += qloc.startPos
	if index == qloc.endPos {
		return -1, types.Error_EndOfEventQueue
	}
	if index > qloc.endPos || index < qloc.startPos {
		return -1, errors.New(fmt.Sprintf("Event queue start pos %d, end pos %d, found invalid start index %d", qloc.startPos, qloc.endPos, index))
	}
	return index, nil
}

type NodeEventQueue struct {
	watchChan chan *event.NodeEvent

	// used to lock enqueue operation during snapshot
	enqueueLock sync.RWMutex

	eventQueueByLoc map[location.Location]*nodeEventQueueByLoc
	locationLock    sync.RWMutex
}

func NewNodeEventQueue(clientId string) *NodeEventQueue {
	queue := &NodeEventQueue{
		eventQueueByLoc: make(map[location.Location]*nodeEventQueueByLoc),
	}

	return queue
}

func (eq *NodeEventQueue) EnqueueEvent(e *event.NodeEvent) {
	eq.enqueueLock.Lock()
	defer eq.enqueueLock.Unlock()
	if eq.watchChan != nil {
		go func() {
			eq.watchChan <- e
		}()
	}

	loc := location.NewLocation(location.Region(e.Node.GeoInfo.Region), location.ResourcePartition(e.Node.GeoInfo.ResourcePartition))

	eq.locationLock.Lock()
	defer eq.locationLock.Unlock()

	queueByLoc, isOK := eq.eventQueueByLoc[*loc]
	if !isOK {
		queueByLoc = newNodeQueueByLoc()
		eq.eventQueueByLoc[*loc] = queueByLoc
	}
	queueByLoc.enqueueEvent(e)
}

func (eq *NodeEventQueue) Watch(rvs types.InternalResourceVersionMap, clientWatchChan chan *event.NodeEvent, stopCh chan struct{}) error {
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
	go func(downstreamCh chan *event.NodeEvent, initEvents []*event.NodeEvent, stopCh chan struct{}, upstreamCh chan *event.NodeEvent) {
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
				downstreamCh <- event
				klog.V(9).Infof("Event with node id %v sent", event.Node.Id)
			}
		}

	}(clientWatchChan, events, stopCh, eq.watchChan)

	return nil
}

func (eq *NodeEventQueue) getAllEventsSinceResourceVersion(rvs types.InternalResourceVersionMap) ([]*event.NodeEvent, error) {
	locStartPostitions := make(map[location.Location]int)

	for loc, rv := range rvs {
		qByLoc, isOK := eq.eventQueueByLoc[loc]
		if isOK {
			startIndex, err := qByLoc.getEventIndexSinceResourceVersion(rv)
			if err != nil {
				if err == types.Error_EndOfEventQueue {
					return nil, nil
				}
				return nil, err
			}
			locStartPostitions[loc] = startIndex
		}
	}

	nodeEvents := make([]*event.NodeEvent, 0)
	for loc, qByLoc := range eq.eventQueueByLoc {
		startIndex, isOK := locStartPostitions[loc]
		var events []*event.NodeEvent
		var err error
		if isOK {
			events, err = qByLoc.getEventsFromIndex(startIndex)
		} else {
			events, err = qByLoc.getEventsFromIndex(qByLoc.startPos)
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
