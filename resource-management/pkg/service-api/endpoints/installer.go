package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/pkg/distributor"
	apiTypes "global-resource-service/resource-management/pkg/service-api/types"
)

//TODO: will move construction of the distributor to main function once each components has basic structures in

var dist = &distributor.ResourceDistributor{}

func init() {
	dist = distributor.GetResourceDistributor()
}

func ResourceHandler(resp http.ResponseWriter, req *http.Request) {
	log.Printf("handle /resource. URL path: %s", req.URL.Path)

	switch req.Method {
	case http.MethodGet:
		ctx := req.Context()
		clientId := ctx.Value("clientid").(string)

		if req.URL.Query().Get(WatchParameter) == WatchParameterTrue {
			serverWatch(resp, req, clientId)
			return
		}

		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "text/plain")

		nodes, _, err := dist.ListNodesForClient(clientId)

		ret, err := json.Marshal(nodes)
		log.Printf("node ret: %s", ret)
		if err != nil {
			log.Printf("error read get node list. error %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.Write(ret)
	case http.MethodPut:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

// simple watch routine
// TODO: add timeout support
// TODO: with serialization options
// TODO: error code and string definition
//
func serverWatch(resp http.ResponseWriter, req *http.Request, clientId string) {
	log.Printf("Serving watch for client: %s", clientId)

	// For 630 distributor impl, watchChannel and stopChannel passed into the Watch routine from API layer
	watchCh := make(chan *event.NodeEvent, WatchChannelSize)
	stopCh := make(chan struct{})

	// Signal the distributor to stop the watch for this client on exit
	defer stopWatch(stopCh)

	// read request body and get the crv
	crvMap, err := getResourceVersionsMap(req)
	if err != nil {
		log.Printf("uUable to get the resource versions. Error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	// start the watcher
	err = dist.Watch(clientId, crvMap, watchCh, stopCh)
	if err != nil {
		log.Printf("uUable to start the watch at store. Error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	done := req.Context().Done()
	flusher, ok := resp.(http.Flusher)
	if !ok {
		log.Printf("unable to start watch - can't get http.Flusher")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	// begin the stream
	resp.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	resp.Header().Set("Transfer-Encoding", "chunked")
	resp.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-done:
			return
		case record, ok := <-watchCh:
			if !ok {
				// End of results.
				return
			}

			if err := json.NewEncoder(resp).Encode(*record.GetNode()); err != nil {
				log.Printf("encoding record failed. error %v", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			if len(watchCh) == 0 {
				flusher.Flush()
			}
		}
	}
}

// Helper functions
func stopWatch(stopCh chan struct{}) {
	stopCh <- struct{}{}
}

func getResourceVersionsMap(req *http.Request) (types.ResourceVersionMap, error) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	wr := apiTypes.WatchRequest{}

	err = json.Unmarshal(body, wr)
	if err != nil {
		return nil, err
	}

	return wr.ResourceVersions, nil
}
