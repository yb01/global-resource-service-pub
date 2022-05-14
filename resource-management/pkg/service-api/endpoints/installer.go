package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"global-resource-service/resource-management/pkg/distributor"
)

//TODO: will move construction of the distributor to main function once each components has basic structures in

var dist = distributor.ResourceDistributor{}

func init() {
	dist = *distributor.GetResourceDistributor()
}

func ResourceHandler(resp http.ResponseWriter, req *http.Request) {
	log.Printf("handle /resource. URL path: %s", req.URL.Path)

	switch req.Method {
	case http.MethodGet:
		ctx := req.Context()
		clientId := ctx.Value("clientid").(string)

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
