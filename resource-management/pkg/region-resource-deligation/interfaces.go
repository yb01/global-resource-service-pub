package region_resource_deligation

import "global-resource-service/resource-management/pkg/common-lib/types"

// RPClientConfigs is a list of clientConfig for each RP in the region
type rPClientConfigs struct {
}

// regionResourceDeligationConfig is the configuration for the regionResourceDeligationConfig API server to run
type regionResourceDeligationConfig struct {
}

type ArktosNode struct {
}

// RP node info retrieval functions
type RPNodeRetriever interface {
	ListNodes() ([]*ArktosNode, error)
	Watch(lastRv int64, watchChan chan *ArktosNode, stopCh chan struct{}) error
}

// Node - logical node translation functions
type LogicalNodeTransformer interface {
	TransformArktosNodeToLocalNode(node ArktosNode) (types.LogicalNode, error)
}

type RpNodeCache struct {
	ArktosNodes  []ArktosNode
	LogicalNodes []types.LogicalNode
}

type ListRpNodeCache []RpNodeCache

// node cache related functions
type RpNodeCacheOperation interface {
	Add(node ArktosNode) error
	Delete(node ArktosNode) error
	Update(node ArktosNode) error
}

// node events related functions
type NodeEvents interface {
}

// handler and rest api functions
