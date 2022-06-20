package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/klog/v2"

	"global-resource-service/resource-management/test/resourceRegionMgrSimulator/app"
	"global-resource-service/resource-management/test/resourceRegionMgrSimulator/data"
)

func main() {
	flag.Usage = printUsage

	// Get the commandline arguments
	c := &app.RegionConfig{}

	flag.IntVar(&c.RegionId, "region_id", 9, "Region Id, if not set, default to 9")
	flag.IntVar(&c.RpNum, "rp_num", 10, "The number of RPs per region, if not set, default to 10")
	flag.IntVar(&c.NodesPerRP, "nodes_per_rp", 3000, "The number of RPs per region, if not set, default to 10")
	flag.StringVar(&c.MasterPort, "master_port", "8080", "Service port, if not set, default to 8080")

	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}

	// Keep a more frequent flush frequency as 1 second
	klog.StartFlushDaemon(time.Second * 1)

	defer klog.Flush()
	klog.Info("")
	klog.Infof("Region resource manager simulator config / region identifier:    (No.%v)", c.RegionId)
	klog.Infof("Region resource manager simulator config / rp number per region: (%v)", c.RpNum)
	klog.Infof("Region resource manager simulator config / node number per rp:  (%v)", c.NodesPerRP)
	klog.Info("")
	klog.Infof("Starting resource region manager simulator (No.%v)", c.RegionId)
	klog.Info("")

	// Initialize Added Event List and Modified Event List
	// Region node Added Event List - for initpull
	data.Init(c.RegionId, c.RpNum, c.NodesPerRP)

	// Generate update changes
	// at ~2 min time, 5k changes
	// at ~5 min mark, 25k changes
	// at ~7 minutes: 1k changes.
	// -- repeat
	data.MakeDataUpdate()

	// Run simulater RSET API server
	if err := app.Run(c); err != nil {
		klog.Errorf("Error: %v\n", err)
	}

	klog.Infof("Exiting reesource management service")
}

// function to print the usage info for the resource management api server
func printUsage() {
	fmt.Println("\nUsage: Region Resource Manager Simulator")
	fmt.Println("\n       Per region config options: --region_id=<region Identifier>  --rp_num=<number of rp>  --nodes_per_rp=<number of nodes> --master_port=<port>")
	fmt.Println()

	os.Exit(0)
}
