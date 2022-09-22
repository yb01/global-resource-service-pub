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

package runtime

import (
	"time"

	common_lib "global-resource-service/resource-management/pkg/common-lib"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
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
		checkpoints: make([]time.Time, common_lib.Len_ResourceManagementCheckpoint),
	}
}

func (e *NodeEvent) SetCheckpoint(checkpoint int) {
	if !common_lib.ResourceManagementMeasurement_Enabled {
		return
	}

	if e.checkpoints == nil {
		e.checkpoints = make([]time.Time, common_lib.Len_ResourceManagementCheckpoint)
	}
	e.checkpoints[checkpoint] = time.Now().UTC()
}

func (e *NodeEvent) GetCheckpoints() []time.Time {
	return e.checkpoints
}

func (e NodeEvent) GetResourceVersionInt64() uint64 {
	return e.Node.GetResourceVersionInt64()
}

func (e NodeEvent) GetGeoInfo() types.NodeGeoInfo {
	return e.Node.GeoInfo
}

func (e NodeEvent) GetId() string {
	return e.Node.Id
}

func (e NodeEvent) GetEventType() EventType {
	return e.Type
}

func (n *NodeEvent) GetLocation() *location.Location {
	return location.NewLocation(location.Region(n.Node.GeoInfo.Region), location.ResourcePartition(n.Node.GeoInfo.ResourcePartition))
}

func (n *NodeEvent) GetLastUpdatedTime() time.Time {
	return n.Node.LastUpdatedTime
}

func (n *NodeEvent) GetEvent() Object {
	return n
}
