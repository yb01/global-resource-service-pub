package interfaces

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

// TODO:
//     Coordinate with Distributor with the interfaces defined here.
//     Distributor will need to make a minor code changes to switch to the new place
type InterfacesOfDistributor interface {
	RegisterClient(requestedHostNum int) (string, bool, error)
	ListNodesForClient(clientId string) ([]*types.Node, types.ResourceVersionMap, error)
	Watch(clientId string, rvs types.ResourceVersionMap, watchChan chan *event.NodeEvent, stopCh chan struct{}) error
	ProcessEvents(events []*event.NodeEvent) (bool, types.ResourceVersionMap)
}
