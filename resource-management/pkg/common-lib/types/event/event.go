package event

import (
	"k8s.io/klog/v2"
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
	checkpoints map[metrics.ResourceManagementCheckpoint]time.Time
}

func NewNodeEvent(node *types.LogicalNode, eventType EventType) *NodeEvent {
	return &NodeEvent{
		Type: eventType,
		Node: node,
	}
}

func (e *NodeEvent) SetCheckpoint(checkpoint metrics.ResourceManagementCheckpoint) {
	if !metrics.ResourceManagementMeasurement_Enabled {
		return
	}
	if e.checkpoints == nil {
		e.checkpoints = make(map[metrics.ResourceManagementCheckpoint]time.Time, 5)
	}
	if _, isOK := e.checkpoints[checkpoint]; !isOK {
		e.checkpoints[checkpoint] = time.Now().UTC()
	} else {
		klog.Errorf("Checkpoint %v already set for event %s, node id %s, rv %s", checkpoint, e.Type, e.Node.Id, e.Node.ResourceVersion)
	}
}

func (e *NodeEvent) GetCheckpoints() map[metrics.ResourceManagementCheckpoint]time.Time {
	return e.checkpoints
}
