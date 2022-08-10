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

package node

import (
	"k8s.io/klog/v2"
	"strconv"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

// TODO - add more fields for minimal node record
type ManagedNodeEvent struct {
	nodeEvent *event.NodeEvent
	loc       *location.Location
}

func NewManagedNodeEvent(nodeEvent *event.NodeEvent, loc *location.Location) *ManagedNodeEvent {
	return &ManagedNodeEvent{
		nodeEvent: nodeEvent,
		loc:       loc,
	}
}

func (n *ManagedNodeEvent) GetId() string {
	return n.nodeEvent.Node.Id
}

func (n *ManagedNodeEvent) GetLocation() *location.Location {
	return n.loc
}

func (n *ManagedNodeEvent) GetRvLocation() *types.RvLocation {

	return &types.RvLocation{Region: n.loc.GetRegion(), Partition: n.loc.GetResourcePartition()}
}

func (n *ManagedNodeEvent) GetResourceVersion() uint64 {
	rv, err := strconv.ParseUint(n.nodeEvent.Node.ResourceVersion, 10, 64)
	if err != nil {
		klog.Errorf("Unable to convert resource version %s to uint64\n", n.nodeEvent.Node.ResourceVersion)
		return 0
	}
	return rv
}

func (n *ManagedNodeEvent) GetEventType() event.EventType {
	return n.nodeEvent.Type
}

func (n *ManagedNodeEvent) GetNodeEvent() *event.NodeEvent {
	return n.nodeEvent
}

func (n *ManagedNodeEvent) CopyNode() *types.LogicalNode {
	return n.nodeEvent.Node.Copy()
}
