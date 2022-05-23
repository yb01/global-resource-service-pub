package event

import "global-resource-service/resource-management/pkg/common-lib/types"

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
	Node *types.LogicalNode
}

func NewNodeEvent(node *types.LogicalNode, eventType EventType) *NodeEvent {
	return &NodeEvent{
		Type: eventType,
		Node: node,
	}
}
