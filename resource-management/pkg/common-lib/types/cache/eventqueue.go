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

	"global-resource-service/resource-management/pkg/common-lib/metrics"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
)

const LengthOfEventQueue = 10000

type EventQueue struct {
	circularEventQueue []runtime.Object
	// circular event queue start position and end position
	startPos int
	endPos   int

	// mutex for event queue operation
	eqLock sync.RWMutex
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		circularEventQueue: make([]runtime.Object, LengthOfEventQueue),
		startPos:           0,
		endPos:             0,
		eqLock:             sync.RWMutex{},
	}
}

func (q *EventQueue) EnqueueEvent(e runtime.Object) {
	q.eqLock.Lock()
	defer q.eqLock.Unlock()

	if q.endPos == q.startPos+LengthOfEventQueue {
		// cache is full - remove the oldest element
		q.startPos++
	}

	q.circularEventQueue[q.endPos%LengthOfEventQueue] = e
	q.endPos++
}

func (q *EventQueue) GetEventsFromIndex(startIndex int) ([]runtime.Object, error) {
	q.eqLock.RLock()
	defer q.eqLock.RUnlock()

	if q.startPos == q.endPos || q.startPos > startIndex || startIndex > q.endPos { // queue is empty or out of range
		return nil, errors.New(fmt.Sprintf("Event queue start pos %d, end pos %d, invalid start index %d", q.startPos, q.endPos, startIndex))
	}

	length := q.endPos - startIndex
	result := make([]runtime.Object, length)
	for i := 0; i < length; i++ {
		result[i] = q.circularEventQueue[(startIndex+i)%LengthOfEventQueue].GetEvent()
	}

	return result, nil
}

func (q *EventQueue) GetEventIndexSinceResourceVersion(resourceVersion uint64) (int, error) {
	q.eqLock.RLock()
	defer q.eqLock.RUnlock()
	if q.endPos-q.startPos == 0 {
		return -1, errors.New(fmt.Sprintf("Empty event queue"))
	}
	e := q.circularEventQueue[q.startPos%LengthOfEventQueue]

	oldestRV := e.GetResourceVersionInt64()
	if oldestRV > resourceVersion {
		return -1, errors.New(fmt.Sprintf("Resource Partition %v events oldest resource Version %d is newer than requested resource version %d",
			e.GetGeoInfo().ResourcePartition, oldestRV, resourceVersion))
	}

	index := sort.Search(q.endPos-q.startPos, func(i int) bool {
		return q.circularEventQueue[(q.startPos+i)%LengthOfEventQueue].GetResourceVersionInt64() > resourceVersion
	})
	index += q.startPos
	if index == q.endPos {
		return -1, types.Error_EndOfEventQueue
	}
	if index > q.endPos || index < q.startPos {
		return -1, errors.New(fmt.Sprintf("Event queue start pos %d, end pos %d, found invalid start index %d", q.startPos, q.endPos, index))
	}
	return index, nil
}

func (q *EventQueue) GetStartPos() int {
	return q.startPos
}

func (q *EventQueue) GetEndPos() int {
	return q.endPos
}

type EventQueuesByLocation struct {
	watchChan chan runtime.Object

	// used to lock enqueue operation during snapshot
	enqueueLock sync.RWMutex

	eventQueueByLoc map[location.Location]*EventQueue
	locationLock    sync.RWMutex
}

func NewEventQueuesByLocation() *EventQueuesByLocation {
	return &EventQueuesByLocation{
		eventQueueByLoc: make(map[location.Location]*EventQueue),
	}
}

func (eq *EventQueuesByLocation) AcquireSnapshotRLock() {
	eq.enqueueLock.RLock()
}

func (eq *EventQueuesByLocation) ReleaseSnapshotRLock() {
	eq.enqueueLock.RUnlock()
}

func (eq *EventQueuesByLocation) EnqueueEvent(e runtime.Object) {
	eq.enqueueLock.Lock()
	defer eq.enqueueLock.Unlock()
	if eq.watchChan != nil {
		go func() {
			eq.watchChan <- e.GetEvent()
		}()
	}

	eq.locationLock.Lock()
	defer eq.locationLock.Unlock()
	queueByLoc, isOK := eq.eventQueueByLoc[*e.GetLocation()]
	if !isOK {
		queueByLoc = NewEventQueue()
		eq.eventQueueByLoc[*e.GetLocation()] = queueByLoc
	}
	queueByLoc.EnqueueEvent(e)
}

func (eq *EventQueuesByLocation) Watch(rvs types.InternalResourceVersionMap, clientWatchChan chan runtime.Object, stopCh chan struct{}) error {
	if eq.watchChan != nil {
		return errors.New("currently only support one watcher per object event queue")
	}

	// get events already in queues
	events, err := eq.getAllEventsSinceResourceVersion(rvs)
	if err != nil {
		return err
	}

	eq.watchChan = make(chan runtime.Object, 1000)
	// writing event to channel
	go func(downstreamCh chan runtime.Object, initEvents []runtime.Object, stopCh chan struct{}, upstreamCh chan runtime.Object) {
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
				klog.V(9).Infof("Sending event with object id %v", event.GetId())
				event.SetCheckpoint(int(metrics.Distributor_Sending))
				downstreamCh <- event
				event.SetCheckpoint(int(metrics.Distributor_Sent))
				klog.V(9).Infof("Event with object id %v sent", event.GetId())
			}
		}

	}(clientWatchChan, events, stopCh, eq.watchChan)

	return nil
}

func (eq *EventQueuesByLocation) getAllEventsSinceResourceVersion(rvs types.InternalResourceVersionMap) ([]runtime.Object, error) {
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

	eventsToReturn := make([]runtime.Object, 0)
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
			eventsToReturn = append(eventsToReturn, events...)
		}
	}

	return eventsToReturn, nil
}
