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
	"sort"
	"sync"

	"global-resource-service/resource-management/pkg/common-lib/types"
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
		result[i] = q.circularEventQueue[(startIndex+i)%LengthOfEventQueue]
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
		return -1, errors.New(fmt.Sprintf("Resource Partition %s events oldest resource Version %d is newer than requested resource version %d",
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
