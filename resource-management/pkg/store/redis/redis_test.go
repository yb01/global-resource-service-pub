package redis

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"testing"
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
