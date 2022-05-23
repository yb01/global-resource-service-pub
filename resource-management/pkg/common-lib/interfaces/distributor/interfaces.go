package distributor

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

type Interface interface {
	RegisterClient(requestedHostNum int) (string, bool, error)
	ListNodesForClient(clientId string) ([]*types.LogicalNode, types.ResourceVersionMap, error)
	Watch(clientId string, rvs types.ResourceVersionMap, watchChan chan *event.NodeEvent, stopCh chan struct{}) error
	ProcessEvents(events []*event.NodeEvent) (bool, types.ResourceVersionMap)
}
