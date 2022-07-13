package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"sync"
	"time"

	"global-resource-service/resource-management/pkg/clientSdk/rmsclient"
	"global-resource-service/resource-management/pkg/clientSdk/tools/cache"
	utilruntime "global-resource-service/resource-management/pkg/clientSdk/util/runtime"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
	"global-resource-service/resource-management/test/e2e/stats"
)

const (
	register = "register"
	list     = "list"
	watch    = "watch"
)

type testConfig struct {
	testDuration time.Duration
	action       string
	repeats      int
}

func main() {
	flag.Usage = printUsage

	testCfg := testConfig{}
	cfg := rmsclient.Config{}
	listOpts := rmsclient.ListOptions{}
	var regions string

	flag.StringVar(&cfg.ServiceUrl, "service_url", "localhost:8080", "Service IP address, if not set, default to localhost")
	flag.DurationVar(&cfg.RequestTimeout, "request_timeout", 30*time.Minute, "Timeout for client requests and responses")
	flag.StringVar(&cfg.ClientFriendlyName, "friendly_name", "testclient", "Client friendly name other that the assigned Id")
	flag.StringVar(&cfg.ClientRegion, "client_region", "Beijing", "Client identify where it is located")
	flag.IntVar(&cfg.InitialRequestTotalMachines, "request_machines", 2500, "Initial request of number of machines")
	flag.StringVar(&regions, "request_regions", "Beijing", "list of regions, in comma separated string, to allocate the machines for the client")
	flag.DurationVar(&testCfg.testDuration, "test_duration", 10*time.Minute, "Test duration, measured by number minutes of watch of node changes. default 10 minutes")
	flag.StringVar(&testCfg.action, "action", "", "action to perform, can be register list or watch, default to register-list-watch")
	flag.IntVar(&testCfg.repeats, "repeats", 1, "number of repeats of the action, default to 1")
	flag.IntVar(&listOpts.Limit, "limit", 25000, "limit for list nodes, default to 25000")

	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}
	klog.StartFlushDaemon(time.Second * 1)
	defer klog.Flush()

	cfg.InitialRequestRegions = strings.Split(regions, ",")
	client := rmsclient.NewRmsClient(cfg)

	store := createStore()

	registerStats := stats.NewRegisterClientStats()
	listStats := stats.NewListStats()
	watchStats := stats.NewWatchStats()

	switch testCfg.action {
	case register:
		for i := 0; i < testCfg.repeats; i++ {
			clientId := registerClient(client, registerStats)
			client.Id = clientId
			printTestStats(registerStats, listStats, watchStats)
		}
	case list:
		clientId := registerClient(client, registerStats)
		client.Id = clientId
		for i := 0; i < testCfg.repeats; i++ {
			_ = listNodes(client, client.Id, store, listStats, listOpts)
			printTestStats(registerStats, listStats, watchStats)
		}
	case watch:
		clientId := registerClient(client, registerStats)
		client.Id = clientId
		crv := listNodes(client, client.Id, store, listStats, listOpts)
		watchNodes(client, client.Id, crv, store, watchStats)
		printTestStats(registerStats, listStats, watchStats)
	default:
		clientId := registerClient(client, registerStats)
		client.Id = clientId
		crv := listNodes(client, client.Id, store, listStats, listOpts)
		watchNodes(client, client.Id, crv, store, watchStats)
		printTestStats(registerStats, listStats, watchStats)
	}

}

// function to print the usage info for the resource management api server
func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("--service_url=127.0.0.1:8080 --request_machines=10000 --request_regions=Beijing,Shanghai ")

	os.Exit(0)
}

func createStore() cache.Store {
	keyFunc := func(obj interface{}) (string, error) {
		node := obj.(types.LogicalNode)
		if len(node.Id) == 0 {
			return "", fmt.Errorf("empty node")
		}
		return node.Id, nil
	}
	return cache.NewStore(keyFunc)

}

func registerClient(client rmsclient.RmsInterface, registerStats *stats.RegisterClientStats) string {
	var start, end time.Time

	klog.Infof("Register client to service  ...")
	start = time.Now().UTC()
	registrationResp, err := client.Register()
	end = time.Now().UTC()
	if err != nil {
		klog.Errorf("failed register client to service. error %v", err)
		os.Exit(1)
	}
	klog.V(6).Infof("Got client registration from service: %v", registrationResp)

	registerStats.RegisterClientDuration = end.Sub(start)
	return registrationResp.ClientId
}

func listNodes(client rmsclient.RmsInterface, clientId string, store cache.Store, listStats *stats.ListStats, listOpts rmsclient.ListOptions) types.TransitResourceVersionMap {
	var start, end time.Time

	klog.Infof("List resources from service ...")
	start = time.Now().UTC()
	nodeList, crv, err := client.List(clientId, listOpts)
	end = time.Now().UTC()
	if err != nil {
		klog.Errorf("failed list resource from service. error %v", err)
		os.Exit(1)
	}
	klog.V(3).Infof("Got [%v] nodes from service", len(nodeList))
	klog.V(3).Infof("Got [%v] resource versions service", crv)

	stats.GroupByRegion(nodeList)
	stats.GroupByRegionByRP(nodeList)

	for _, node := range nodeList {
		store.Add(*node)
	}

	listStats.ListDuration = end.Sub(start)
	listStats.NumberOfNodesListed = len(nodeList)
	return crv
}

func watchNodes(client rmsclient.RmsInterface, clientId string, crv types.TransitResourceVersionMap, store cache.Store,
	watchStats *stats.WatchStats) {
	var start, end time.Time

	klog.Infof("Watch resources update from service ...")
	start = time.Now().UTC()
	watcher, err := client.Watch(clientId, crv)
	if err != nil {
		klog.Errorf("failed list resource from service. error %v", err)
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
				logIfProlonged(&record, time.Now().UTC(), watchStats)
				switch record.Type {
				case event.Added:
					store.Add(*record.Node)
					watchStats.NumberOfAddedNodes++
				case event.Modified:
					store.Update(*record.Node)
					watchStats.NumberOfUpdatedNodes++
				case event.Deleted:
					store.Delete(*record.Node)
					watchStats.NumberOfDeletedNodes++

				default:
					klog.Error("not supported event type")
				}

			}
		}
	}()
	wg.Wait()
	end = time.Now().UTC()
	watchStats.WatchDuration = end.Sub(start)
	return
}

func logIfProlonged(record *event.NodeEvent, t time.Time, ws *stats.WatchStats) {
	d := t.Sub(record.Node.LastUpdatedTime)
	if d > stats.LongWatchThreshold {
		klog.Infof("Prolonged watch node from server: %v, %v with time (%v)", record, *record.Node, d)
		ws.NumberOfProlongedItems++
	}
}

func printTestStats(rs *stats.RegisterClientStats, ls *stats.ListStats, ws *stats.WatchStats) {
	rs.PrintStats()
	ls.PrintStats()
	ws.PrintStats()
}
