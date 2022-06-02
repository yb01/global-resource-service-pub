package redis

import (
	"context"
	"encoding/json"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/pkg/common-lib/interfaces/store"
	"global-resource-service/resource-management/pkg/common-lib/types"

	"github.com/go-redis/redis/v8"
)

type Goredis struct {
	client *redis.Client
	ctx    context.Context
}

// Initialize Redis Client
//
func NewRedisClient() *Goredis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	return &Goredis{
		client: client,
		ctx:    ctx,
	}
}

// To Test Persist Simple String
//
func (gr *Goredis) setString(myKey, myValue string) bool {
	if len(myKey) == 0 || len(myValue) == 0 {
		klog.Errorf("The Key or Value is blank")
		return false
	}

	err := gr.client.Set(gr.ctx, myKey, []byte(myValue), 0).Err()

	if err != nil {
		klog.Errorf("Error to persist String key and Value to Redis Store:", err)
		return false
	}

	return true
}

// To Test Get Simple String
//
func (gr *Goredis) getString(myKey string) string {
	var myValue string

	if len(myKey) == 0 {
		klog.Errorf("The Key is blank")
		return ""
	}

	value, err := gr.client.Get(gr.ctx, myKey).Bytes()

	if err != nil {
		klog.Errorf("Error to get String Key from Redis Store:", err)
		return ""
	}

	if err != redis.Nil {
		myValue = string(value)
	}

	return myValue
}

// Use Redis data type - Set to store Logical Nodes
// One key has one record
//
// Note: Need re-visit these codes to see whether using function pointer is much better
//
// TODO: Error handling for loop persistence failure in the middle
//
func (gr *Goredis) PersistNodes(logicalNodes []*types.LogicalNode) bool {
	if logicalNodes == nil {
		klog.Errorf("The array of Logical Nodes is nil")
		return false
	}

	for _, logicalNode := range logicalNodes {
		logicalNodeKey := logicalNode.GetKey()
		logicalNodeBytes, err := json.Marshal(logicalNode)

		if err != nil {
			klog.Errorf("Error from JSON Marshal for Logical Nodes:", err)
			return false
		}

		err = gr.client.Set(gr.ctx, logicalNodeKey, logicalNodeBytes, 0).Err()

		if err != nil {
			klog.Errorf("Error to persist Logical Nodes to Redis Store:", err)
			return false
		}
	}

	return true
}

// Use Redis data type - String to store Node Store Status
//
func (gr *Goredis) PersistNodeStoreStatus(nodeStoreStatus *store.NodeStoreStatus) bool {
	nodeStoreStatusBytes, err := json.Marshal(nodeStoreStatus)

	if err != nil {
		klog.Errorf("Error from JSON Marshal for Node Store Status:", err)
		return false
	}

	err = gr.client.Set(gr.ctx, nodeStoreStatus.GetKey(), nodeStoreStatusBytes, 0).Err()

	if err != nil {
		klog.Errorf("Error to persist Node Store Status to Redis Store:", err)
		return false
	}

	return true
}

// Use Redis data type - String to store Virtual Node Assignment
//
func (gr *Goredis) PersistVirtualNodesAssignments(virtualNodeAssignment *store.VirtualNodeAssignment) bool {
	virtualNodeAssignmentBytes, err := json.Marshal(virtualNodeAssignment)

	if err != nil {
		klog.Errorf("Error from JSON Marshal for Virtual Node Assignment:", err)
		return false
	}

	err = gr.client.Set(gr.ctx, virtualNodeAssignment.GetKey(), virtualNodeAssignmentBytes, 0).Err()

	if err != nil {
		klog.Errorf("Error to persist Virtual Node Assignment to Redis Store:", err)
		return false
	}

	return true
}

// Get all Logical Nodes based on PreserveNode_KeyPrefix = "MinNode"
//
// Note: Need re-visit these codes to see whether using function pointer is much better
//
func (gr *Goredis) GetNodes() []*types.LogicalNode {
	keys := gr.client.Keys(gr.ctx, types.PreserveNode_KeyPrefix+"*").Val()

	logicalNodes := make([]*types.LogicalNode, len(keys))

	var logicalNode *types.LogicalNode

	for i, logicalNodeKey := range keys {
		value, err := gr.client.Get(gr.ctx, logicalNodeKey).Bytes()

		if err != nil {
			klog.Errorf("Error to get LogicalNode from Redis Store:", err)
			return nil
		}

		if err != redis.Nil {
			err = json.Unmarshal(value, &logicalNode)

			if err != nil {
				klog.Errorf("Error from JSON Unmarshal for LogicalNode:", err)
				return nil
			}

			logicalNodes[i] = logicalNode

		}
	}

	return logicalNodes
}

// Get Node Store Status
//
func (gr *Goredis) GetNodeStoreStatus() *store.NodeStoreStatus {
	var nodeStoreStatus *store.NodeStoreStatus

	value, err := gr.client.Get(gr.ctx, nodeStoreStatus.GetKey()).Bytes()

	if err != nil {
		klog.Errorf("Error to get NodeStoreStatus from Redis Store:", err)
		return nil
	}

	if err != redis.Nil {
		err = json.Unmarshal(value, &nodeStoreStatus)

		if err != nil {
			klog.Errorf("Error from JSON Unmarshal for NodeStoreStatus:", err)
			return nil
		}
	}

	return nodeStoreStatus
}

// Get Virtual Nodes Assignments
//
func (gr *Goredis) GetVirtualNodesAssignments() *store.VirtualNodeAssignment {
	var virtualNodeAssignment *store.VirtualNodeAssignment

	value, err := gr.client.Get(gr.ctx, virtualNodeAssignment.GetKey()).Bytes()

	if err != nil {
		klog.Errorf("Error to get VirtualNodeAssignment from Redis Store:", err)
		return nil
	}

	if err != redis.Nil {
		err = json.Unmarshal(value, &virtualNodeAssignment)

		if err != nil {
			klog.Errorf("Error from JSON Unmarshal for VirtualNodeAssignment:", err)
			return nil
		}
	}

	return virtualNodeAssignment
}
