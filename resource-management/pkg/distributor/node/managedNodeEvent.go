package node

import (
	"fmt"
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

func NewManagedNodeEvent1(id, rv string, eventType event.EventType, location *location.Location) *ManagedNodeEvent {
	nodeEvent := event.NewNodeEvent(
		&types.LogicalNode{Id: id, ResourceVersion: rv},
		eventType)

	return &ManagedNodeEvent{
		nodeEvent: nodeEvent,
		loc:       location,
	}
}

func (n *ManagedNodeEvent) GetId() string {
	return n.nodeEvent.Node.Id
}

func (n *ManagedNodeEvent) GetLocation() *location.Location {
	return n.loc
}

func (n *ManagedNodeEvent) GetResourceVersion() uint64 {
	rv, err := strconv.ParseUint(n.nodeEvent.Node.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Printf("Unable to convert resource version %s to uint64\n", n.nodeEvent.Node.ResourceVersion)
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
