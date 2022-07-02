package redis

import (
	"reflect"
	"testing"

	"global-resource-service/resource-management/pkg/common-lib/interfaces/store"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

var GR *Goredis

func init() {
	GR = NewRedisClient()
}

func TestNewRedisClient(t *testing.T) {
	GR = NewRedisClient()
}

// Simply Test Set String and Get String without the need of Marshal and Unmarshal
//
func TestSetGettString(t *testing.T) {
	success := GR.setString("testkey1NM", "testvalue1NM")
	if success == false {
		t.Error("Store testkey1NM into Redis in failure")
	}

	value := GR.getString("testkey1NM")
	if value == "" {
		t.Error("Read value of testkey1NM from redis in failure")
	}

	if value != "testvalue1NM" {
		t.Errorf("Unexpected value %s", string(value))
	}
}

// Simply Test Persist Nodes and Get Nodes
//
func TestPersistNodes(t *testing.T) {
	testCases := make([]*types.LogicalNode, 1)

	var testCase0 = &types.LogicalNode{
		Id:              "0001",
		ResourceVersion: "0002",
		GeoInfo: types.NodeGeoInfo{
			Region:            1000,
			ResourcePartition: 1000,
		},
		Conditions:  255,
		Reserved:    true,
		MachineType: "machineType1",
	}

	testCases[0] = testCase0

	success := GR.PersistNodes(testCases)

	if success == false {
		t.Error("Store Logical Nodes of testCases into Redis in failure")
	}

	logicalNodes := GR.GetNodes()

	if len(logicalNodes) == 0 {
		t.Error("Zero Logical Node is read from Redis")
	}

	if logicalNodes[0].Id != testCases[0].Id || logicalNodes[0].ResourceVersion != testCases[0].ResourceVersion {
		t.Error("logicalNodes[0] is : ", *logicalNodes[0])
		t.Error("testCases[0]    is : ", *testCases[0])
	}

	if logicalNodes[0].GeoInfo.Region != testCases[0].GeoInfo.Region || logicalNodes[0].GeoInfo.ResourcePartition != testCases[0].GeoInfo.ResourcePartition {
		t.Error("logicalNodes[0].GeoInfo.Region is : ", logicalNodes[0].GeoInfo.Region)
		t.Error("testCases[0].GeoInfo.Region    is : ", testCases[0].GeoInfo.Region)
		t.Error("logicalNodes[0].GeoInfo.ResourcePartition is : ", logicalNodes[0].GeoInfo.ResourcePartition)
		t.Error("testCases[0].GeoInfo.ResourcePartition    is : ", testCases[0].GeoInfo.ResourcePartition)
	}
}

// Simply Test Persist Node Store Status
//
func TestPersistNodeStoreStatus(t *testing.T) {
	var CRV = make(types.TransitResourceVersionMap, 1)
	testLocation := types.RvLocation{Region: location.Beijing, Partition: location.ResourcePartition1}
	CRV[testLocation] = 1000

	testCase0 := &store.NodeStoreStatus{
		RegionNum:              1000,
		PartitionMaxNum:        1000,
		VirtualNodeNumPerRP:    1000,
		CurrentResourceVerions: CRV,
	}

	success := GR.PersistNodeStoreStatus(testCase0)

	if success == false {
		t.Error("Store Node Store Status of testCase into Redis in failure")
	}

	var nodeStoreStatus = GR.GetNodeStoreStatus()

	if reflect.ValueOf(nodeStoreStatus).IsNil() {
		t.Error("Read Node Store Status from Redis in failure")
	}

	if nodeStoreStatus.RegionNum != testCase0.RegionNum {
		t.Error("nodeStoreStatus.RegionNum is : ", nodeStoreStatus.RegionNum)
		t.Error("testCases0.RegionNum      is : ", testCase0.RegionNum)
	}

	if nodeStoreStatus.PartitionMaxNum != testCase0.PartitionMaxNum {
		t.Error("nodeStoreStatus.PartitionMaxNum is : ", nodeStoreStatus.PartitionMaxNum)
		t.Error("testCases0.PartitionMaxNum      is : ", testCase0.PartitionMaxNum)
	}

	if nodeStoreStatus.VirtualNodeNumPerRP != testCase0.VirtualNodeNumPerRP {
		t.Error("nodeStoreStatus.VirtualNodeNumPerRP is : ", nodeStoreStatus.VirtualNodeNumPerRP)
		t.Error("testCases0.VirtualNodeNumPerRP      is : ", testCase0.VirtualNodeNumPerRP)
	}

	if nodeStoreStatus.CurrentResourceVerions[testLocation] != testCase0.CurrentResourceVerions[testLocation] {
		t.Error("nodeStoreStatus.[testLocation] is : ", testLocation, nodeStoreStatus.CurrentResourceVerions[testLocation])
		t.Error("testCases0.[testLocation]      is : ", testLocation, testCase0.CurrentResourceVerions[testLocation])
	}
}

// Simply Test Persist Virtual Nodes Assignments
//
func TestPersistVirtualNodesAssignments(t *testing.T) {
	vNodeConfigs := make([]*store.VirtualNodeConfig, 1)

	vNodeToSave := &store.VirtualNodeConfig{
		Lowerbound: 1000.00,
		Upperbound: 2000.00,
		Location:   *location.NewLocation(location.Beijing, location.ResourcePartition1),
	}

	vNodeConfigs[0] = vNodeToSave

	testCase0 := &store.VirtualNodeAssignment{
		ClientId:     "1000",
		VirtualNodes: vNodeConfigs,
	}

	success := GR.PersistVirtualNodesAssignments(testCase0)

	if success == false {
		t.Error("Store Node Store Status of testCase into Redis in failure")
	}

	var virtualNodeAssignment = GR.GetVirtualNodesAssignments()

	if reflect.ValueOf(virtualNodeAssignment).IsNil() {
		t.Error("Read Virtual Node Assignment from Redis in failure")
	}

	if virtualNodeAssignment.ClientId != testCase0.ClientId {
		t.Error("virtualNodeAssignment.ClientId is : ", virtualNodeAssignment.ClientId)
		t.Error("testCases0.ClientId      is : ", testCase0.ClientId)
	}

	if virtualNodeAssignment.VirtualNodes[0].Lowerbound != testCase0.VirtualNodes[0].Lowerbound {
		t.Error("virtualNodeAssignment.VirtualNodes[0].Lowerbound is : ", virtualNodeAssignment.VirtualNodes[0].Lowerbound)
		t.Error("testCases0.VirtualNodes[0].Lowerbound      is : ", testCase0.VirtualNodes[0].Lowerbound)
	}

	if virtualNodeAssignment.VirtualNodes[0].Upperbound != testCase0.VirtualNodes[0].Upperbound {
		t.Error("virtualNodeAssignment.VirtualNodes[0].Upperbound is : ", virtualNodeAssignment.VirtualNodes[0].Upperbound)
		t.Error("testCases0.VirtualNodes[0].Upperbound      is : ", testCase0.VirtualNodes[0].Upperbound)
	}

	if float64(virtualNodeAssignment.VirtualNodes[0].Location.GetRegion()) != float64(testCase0.VirtualNodes[0].Location.GetRegion()) {
		t.Error("virtualNodeAssignment.VirtualNodes[0].Location.GetRegion() is : ", virtualNodeAssignment.VirtualNodes[0].Location.GetRegion())
		t.Error("testCases0.VirtualNodes[0].Location.GetRegion()      is : ", testCase0.VirtualNodes[0].Location.GetRegion())
	}

	if float64(virtualNodeAssignment.VirtualNodes[0].Location.GetResourcePartition()) != float64(testCase0.VirtualNodes[0].Location.GetResourcePartition()) {
		t.Error("virtualNodeAssignment.VirtualNodes[0].Location.GetResourcePartition() is : ", virtualNodeAssignment.VirtualNodes[0].Location.GetResourcePartition())
		t.Error("testCases0.VirtualNodes[0].Location.GetResourcePartition()      is : ", testCase0.VirtualNodes[0].Location.GetResourcePartition())
	}
}
