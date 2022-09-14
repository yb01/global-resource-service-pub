/*
Copyright 2022 Authors of Global Resource Service.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aggregrator

import (
	"k8s.io/klog/v2"
	"sync"
	"time"

	utilruntime "global-resource-service/resource-management/pkg/clientSdk/util/runtime"
	common_lib "global-resource-service/resource-management/pkg/common-lib"
	"global-resource-service/resource-management/pkg/common-lib/metrics"
	"global-resource-service/resource-management/pkg/common-lib/types"
	event "global-resource-service/resource-management/pkg/common-lib/types/runtime"
)

// LwRun implements Run interface of Aggregator
func (a *Aggregator) Run() (err error) {
	numberOfRegions := len(a.urls)

	klog.V(3).Infof("Running for loop to connect to to resource region manager...")

	for i := 0; i < numberOfRegions; i++ {
		go func(i int) {
			klog.V(3).Infof("Starting goroutine for region: %v", a.urls[i])
			var err error
			defer func() {
				klog.V(3).Infof("Exiting goroutine for region: %v", a.urls[i])
				if err != nil {
					klog.Errorf("Error for goroutine: %v", err)
				}
			}()

			var crv types.TransitResourceVersionMap
			var regionNodeEvents [][]*event.NodeEvent
			var length uint64
			var eventProcess bool

			// create client to resource region manager
			c := NewRrmsClient(Config{ServiceUrl: a.urls[i], RequestTimeout: 30 * time.Minute})

			klog.V(3).Infof("Starting loop list-watching nodes from region: %v", a.urls[i])

			// TODO: add limit in config
			// TODO: add pagination feature in region resource mgr service
			// hack, list all 1m node in one call without pagination
			regionNodeEvents, crv, length, err = a.listNodes(c, ListOptions{Limit: 1000000}, a.urls[i])

			// Exit the goroutine with failed List operation
			if err != nil {
				return
			}

			if length != 0 {
				klog.V(4).Infof("Total (%v) region node events are listed successfully in (%v) RPs", length, len(regionNodeEvents))
			} else {
				// TODO: handel empty list
			}

			// Convert 2D array to 1D array
			// TODO: add dynamically watch at RP granularity level feature in GRS
			//       by spawning threads for watch each RP, we can further optimize concurrency and also avoid below conversions
			minRecordNodeEvents := make([]*event.NodeEvent, 0, length)
			for j := 0; j < len(regionNodeEvents); j++ {
				minRecordNodeEvents = append(minRecordNodeEvents, regionNodeEvents[j]...)
			}

			start := time.Now()
			eventProcess, crv = a.EventProcessor.ProcessEvents(minRecordNodeEvents)
			end := time.Now()
			klog.V(6).Infof("Event Processor Processed nodes results : %v. duration: %v", eventProcess, end.Sub(start))

			// start watch node changes
			a.watchNodes(c, crv, a.urls[i])

		}(i)
	}

	klog.V(3).Infof("Finished for loop to connect to to resource region manager...")
	return nil
}

func (a *Aggregator) listNodes(client RrmsInterface, listOpts ListOptions, url string) (nodeList [][]*event.NodeEvent,
	crv types.TransitResourceVersionMap, length uint64, err error) {
	var start, end time.Time

	retryList := 0
	for {
		klog.Infof("List resources from region manager %v", url)
		start = time.Now().UTC()
		nodeList, crv, length, err = client.List(listOpts)
		end = time.Now().UTC()
		if err != nil {
			klog.Warningf("failed list resource from region manager. error %v. retry in one second", err)
			//TODO: chain to the clientSDK and with necessary changes to bubble up the http error code
			//      and take action on the different type of errors
			// for now, the common cost in test env is the region manager started after the GRS service, and hence the wait here
			time.Sleep(1 * time.Second)
			if retryList == 60 {
				klog.Errorf("failed list resource from region manager after retries. error %v", err)
				return nil, nil, 0, err
			}
			retryList++
			continue
		} else {
			break
		}
	}
	klog.V(3).Infof("Got [%v] RPs, [%v] nodes from region manager, list duration: %v", len(nodeList), length, end.Sub(start))

	if common_lib.ResourceManagementMeasurement_Enabled {
		for i := 0; i < len(nodeList); i++ {
			for j := 0; j < len(nodeList[i]); j++ {
				if nodeList[i][j] != nil {
					nodeList[i][j].SetCheckpoint(int(metrics.Aggregator_Received))
				}
			}
		}
	}

	return nodeList, crv, length, nil
}

func (a *Aggregator) watchNodes(client RrmsInterface, crv types.TransitResourceVersionMap, url string) {
	var start, end time.Time

	klog.V(3).Infof("Watch resources update from region manager %v", url)
	start = time.Now().UTC()
	watcher, err := client.Watch(crv)
	if err != nil {
		klog.Errorf("failed list resource from region manager. error %v", err)
	}

	watchCh := watcher.ResultChan()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer utilruntime.HandleCrash()
		// retrieve updates from watcher
		for {
			select {
			case record, ok := <-watchCh:
				if !ok {
					// End of results.
					klog.Infof("End of results")
					return
				}

				klog.V(9).Infof("Got node event from region manager, nodeId: %v", record.Node.Id)

				// TODO: refine this go routine to sub functions
				// TODO: persists the CRV returned from processEvents() for re-watch
				go func() {
					a.processNode(&record)
				}()
			}
		}
	}()
	wg.Wait()
	end = time.Now().UTC()
	klog.V(3).Infof("Watch session last: %v", end.Sub(start))
	return
}

// TODO: lock this function if the distributor cannot handel concurrent node processing
func (a *Aggregator) processNode(node *event.NodeEvent) {
	node.SetCheckpoint(int(metrics.Aggregator_Received))
	a.EventProcessor.ProcessEvents([]*event.NodeEvent{node})
}
