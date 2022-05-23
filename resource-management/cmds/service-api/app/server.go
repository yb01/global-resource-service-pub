package app

import (
	"net/http"
	"strconv"
	"time"

	"global-resource-service/resource-management/pkg/service-api/endpoints"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
)

// Run and create new service-api.  This should never exit.
func Run() error {
	klog.V(3).Infof("Starting the API server...")

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc(endpoints.ListResourcePath, endpoints.ResourceHandler).Methods("GET")
	r.HandleFunc(endpoints.WatchResourcePath, endpoints.ResourceHandler)
	r.HandleFunc(endpoints.UpdateResourcePath, endpoints.ResourceHandler)
	r.HandleFunc(endpoints.ReduceResourcePath, endpoints.ResourceHandler)

	server := &http.Server{
		Handler:      r,
		Addr:         "localhost:" + strconv.Itoa(endpoints.InsecureServiceAPIPort),
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	return server.ListenAndServe()
}
