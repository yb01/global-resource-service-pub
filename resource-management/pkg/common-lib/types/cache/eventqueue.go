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
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
	"sort"
	"sync"
)

const LengthOfNodeEventQueue = 10000

type EventQueue struct {
	circularEventQueue []*runtime.Object
	// circular event queue start position and end position
	startPos int
	endPos   int

	// mutex for event queue operation
	eqLock sync.RWMutex
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		circularEventQueue: make([]*event.NodeEvent, LengthOfNodeEventQueue),
		startPos:           0,
		endPos:             0,
		eqLock:             sync.RWMutex{},
	}
}

func (q *EventQueue) enqueueEvent(e *event.NodeEvent) {
	q.eqLock.Lock()
	defer q.eqLock.Unlock()

	if q.endPos == q.startPos+LengthOfNodeEventQueue {
		// cache is full - remove the oldest element
		q.startPos++
	}

	q.circularEventQueue[q.endPos%LengthOfNodeEventQueue] = e
	q.endPos++
}

func (q *EventQueue) getEventsFromIndex(startIndex int) ([]*event.NodeEvent, error) {
	q.eqLock.RLock()
	defer q.eqLock.RUnlock()

	if q.startPos == q.endPos || q.startPos > startIndex || startIndex > q.endPos { // queue is empty or out of range
		return nil, errors.New(fmt.Sprintf("Event queue start pos %d, end pos %d, invalid start index %d", q.startPos, q.endPos, startIndex))
	}

	length := q.endPos - startIndex
	result := make([]*event.NodeEvent, length)
	for i := 0; i < length; i++ {
		result[i] = q.circularEventQueue[(startIndex+i)%LengthOfNodeEventQueue]
	}

	return result, nil
}

func (q *EventQueue) getEventIndexSinceResourceVersion(resourceVersion uint64) (int, error) {
	q.eqLock.RLock()
	defer q.eqLock.RUnlock()
	if q.endPos-q.startPos == 0 {
		return -1, errors.New(fmt.Sprintf("Empty event queue"))
	}
	nodeEvent := q.circularEventQueue[q.startPos%LengthOfNodeEventQueue].Node

	oldestRV := nodeEvent.GetResourceVersionInt64()
	if oldestRV > resourceVersion {
		return -1, errors.New(fmt.Sprintf("Resource Partition %s events oldest resource Version %d is newer than requested resource version %d",
			nodeEvent.GeoInfo.ResourcePartition, oldestRV, resourceVersion))
	}

	index := sort.Search(q.endPos-q.startPos, func(i int) bool {
		return q.circularEventQueue[(q.startPos+i)%LengthOfNodeEventQueue].Node.GetResourceVersionInt64() > resourceVersion
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
