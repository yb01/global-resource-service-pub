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
	"time"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"global-resource-service/resource-management/pkg/common-lib/types/runtime"
)

// TODO - add more fields for minimal node record
type ManagedNodeEvent struct {
	nodeEvent *runtime.NodeEvent
	loc       *location.Location
}

func NewManagedNodeEvent(nodeEvent *runtime.NodeEvent, loc *location.Location) *ManagedNodeEvent {
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

func (n *ManagedNodeEvent) GetResourceVersionInt64() uint64 {
	return n.nodeEvent.Node.GetResourceVersionInt64()
}

func (n *ManagedNodeEvent) GetGeoInfo() types.NodeGeoInfo {
	return n.nodeEvent.Node.GeoInfo
}

func (n *ManagedNodeEvent) GetEventType() runtime.EventType {
	return n.nodeEvent.Type
}

func (n *ManagedNodeEvent) GetNodeEvent() *runtime.NodeEvent {
	return n.nodeEvent
}

func (n *ManagedNodeEvent) CopyNode() *types.LogicalNode {
	return n.nodeEvent.Node.Copy()
}

func (n *ManagedNodeEvent) SetCheckpoint(int) {
	klog.Error("Not implemented SetCheckpoint method")
}

func (e *ManagedNodeEvent) GetCheckpoints() []time.Time {
	klog.Error("Not implemented GetCheckpoints method")
	return nil
}

func (n *ManagedNodeEvent) GetLastUpdatedTime() time.Time {
	return n.nodeEvent.Node.LastUpdatedTime
}

func (n *ManagedNodeEvent) GetEvent() runtime.Object {
	return n.nodeEvent
}
