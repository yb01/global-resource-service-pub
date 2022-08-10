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

package event

import (
	"time"

	"global-resource-service/resource-management/pkg/common-lib/metrics"
	"global-resource-service/resource-management/pkg/common-lib/types"
)

// EventType defines the possible types of events.
type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Bookmark EventType = "BOOKMARK"
	Error    EventType = "ERROR"
)

type NodeEvent struct {
	Type        EventType
	Node        *types.LogicalNode
	checkpoints []time.Time
}

func NewNodeEvent(node *types.LogicalNode, eventType EventType) *NodeEvent {
	return &NodeEvent{
		Type:        eventType,
		Node:        node,
		checkpoints: make([]time.Time, metrics.Len_ResourceManagementCheckpoint),
	}
}

func (e *NodeEvent) SetCheckpoint(checkpoint metrics.ResourceManagementCheckpoint) {
	if !metrics.ResourceManagementMeasurement_Enabled {
		return
	}

	if e.checkpoints == nil {
		e.checkpoints = make([]time.Time, metrics.Len_ResourceManagementCheckpoint)
	}
	e.checkpoints[checkpoint] = time.Now().UTC()
}

func (e *NodeEvent) GetCheckpoints() []time.Time {
	return e.checkpoints
}
