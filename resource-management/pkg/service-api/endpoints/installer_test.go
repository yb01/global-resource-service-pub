package endpoints

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/distributor"
	"global-resource-service/resource-management/pkg/distributor/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//TODO: will have mock interface impl once the interface is moved out to the common lib
var existedNodeId = make(map[uuid.UUID]bool)
var rvToGenerate = 0

var singleTestLock = sync.Mutex{}

func setUp() *distributor.ResourceDistributor {
	singleTestLock.Lock()
	return distributor.GetResourceDistributor()
}

func tearDown(resourceDistributor *distributor.ResourceDistributor) {
	defer singleTestLock.Unlock()
}

func createRandomNode(rv int) *types.LogicalNode {
	id := uuid.New()

	return &types.LogicalNode{
		Id:              id.String(),
		ResourceVersion: strconv.Itoa(rv),
		GeoInfo: types.NodeGeoInfo{
			Region:            0,
			ResourcePartition: 0,
		},
	}
}

func generateAddNodeEvent(eventNum int) []*event.NodeEvent {
	result := make([]*event.NodeEvent, eventNum)
	for i := 0; i < eventNum; i++ {
		rvToGenerate += 1
		node := createRandomNode(rvToGenerate)
		nodeEvent := event.NewNodeEvent(node, event.Added)
		result[i] = nodeEvent
	}
	return result
}

func TestHttpGet(t *testing.T) {
	distributor := setUp()
	defer tearDown(distributor)

	fakeStorage := &storage.FakeStorageInterface{
		PersistDelayInNS: 20,
	}
	distributor.SetPersistHelper(fakeStorage)
	installer := NewInstaller(distributor)

	// initialize node store with 10K nodes
	eventsAdd := generateAddNodeEvent(10000)
	distributor.ProcessEvents(eventsAdd)

	//register client
	requestedHostNum := 500
	clientId, _, err := distributor.RegisterClient(requestedHostNum)

	// client list nodes
	expectedNodes, _, err := distributor.ListNodesForClient(clientId)

	if err != nil {
		t.Fatal(err)
	}

	resourcePath := RegionlessResourcePath + "/" + clientId
	req, err := http.NewRequest(http.MethodGet, resourcePath, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "clientid", clientId)

	installer.ResourceHandler(recorder, req.WithContext(ctx))

	actualNodes := make([]types.LogicalNode, 150)

	json.Unmarshal(recorder.Body.Bytes(), &actualNodes)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, len(expectedNodes), len(actualNodes))

	// Node list is not ordered, so have to do a one by one comparison
	for _, n := range expectedNodes {
		if findNodeInList(n, actualNodes) == false {
			t.Logf("expectd node Id [%v] not found", n.Id)
			t.Fatal("Nodes are not equal")
		}
	}

	return
}

func findNodeInList(n *types.LogicalNode, l []types.LogicalNode) bool {
	for i := 0; i < len(l); i++ {
		if n.Id == l[i].Id {
			return true
		}
	}

	return false
}
