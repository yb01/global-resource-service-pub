package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"

	"global-resource-service/resource-management/test/resourceRegionMgrSimulator/handlers"
)

type RegionConfig struct {
	RegionId   int
	RpNum      int
	NodesPerRP int
	MasterPort string
}

func Run(c *RegionConfig) error {

	// Create the handlers
	rh := handlers.NewRegionNodeEventsHander()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter().StrictSlash(true)

	// handlers for GET API
	//
	getRouter := sm.Methods(http.MethodGet).Subrouter()

	// For initial pull all mini node added events in all RPs of one specified region
	getRouter.HandleFunc(handlers.InitPullPath, rh.SimulatorHandler)

	// For subsequent pull all mini node modified events in all RPs of one specified region
	getRouter.HandleFunc(handlers.SubsequentPullPath, rh.SimulatorHandler)

	// handlers for POST API
	//
	postRouter := sm.Methods(http.MethodPost).Subrouter()

	// with CRV, discard all mini node modified events
	// which resource version is older than CRV in all RPs of specified region
	postRouter.HandleFunc(handlers.PostCRVPath, rh.SimulatorHandler)

	var bindAddress = ":" + c.MasterPort

	// Define HTTP Server
	s := &http.Server{
		Addr:         bindAddress,
		Handler:      sm,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		klog.Infof("\nStarting resource region manager simulator server on port (%v)", c.MasterPort)

		err := s.ListenAndServe()
		if err != nil {
			klog.Errorf("The HTTP server is not gracefully shutdown : ", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	klog.Infof("\n\nReceived terminate, graceful shutdown\n\n", sig)

	tc, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		klog.Errorf("The HTTP server is not gracefully shutdown : ", err)
	}
	s.Shutdown(tc)
	return nil
}
