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
