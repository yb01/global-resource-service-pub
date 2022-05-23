package endpoints

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/distributor"

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

	// initialize node store with 10K nodes
	eventsAdd := generateAddNodeEvent(10000)
	distributor.ProcessEvents(eventsAdd)

	//register client
	requestedHostNum := 500
	clientId, _, err := distributor.RegisterClient(requestedHostNum)

	// client list nodes
	//nodes, _, err := distributor.ListNodesForClient(clientId)
	resourcePath := RegionlessResourcePath + "/" + clientId
	req, err := http.NewRequest(http.MethodGet, resourcePath, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "clientid", clientId)

	ResourceHandler(recorder, req.WithContext(ctx))
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, []byte(clientId), recorder.Body.Bytes())
}
