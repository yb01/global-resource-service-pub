package event

import "resource-management/pkg/common-lib/types"

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
	Type EventType
	node *types.Node
}

func NewNodeEvent(node *types.Node, eventType EventType) *NodeEvent {
	return &NodeEvent{
		Type: eventType,
		node: node,
	}
}

func (e *NodeEvent) GetNode() *types.Node {
	return e.node
}
