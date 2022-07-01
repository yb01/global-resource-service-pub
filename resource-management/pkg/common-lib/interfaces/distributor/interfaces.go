package distributor

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

type Interface interface {
	RegisterClient(*types.Client) error

	ListNodesForClient(clientId string) ([]*types.LogicalNode, types.TransitResourceVersionMap, error)
	Watch(clientId string, rvs types.TransitResourceVersionMap, watchChan chan *event.NodeEvent, stopCh chan struct{}) error
	ProcessEvents(events []*event.NodeEvent) (bool, types.TransitResourceVersionMap)
}
