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

	flag.StringVar(&c.RegionName, "region_name", "Beijing", "Region name, if not set, default to Beijing")
	flag.IntVar(&c.RpNum, "rp_num", 10, "The number of RPs per region, if not set, default to 10")
	flag.IntVar(&c.NodesPerRP, "nodes_per_rp", 25000, "The number of RPs per region, if not set, default to 25000")
	flag.StringVar(&c.MasterPort, "master_port", "9119", "Service port, if not set, default to 9119")
	flag.StringVar(&c.DataPattern, "data_pattern", "Outage", "Simulator data pattern, if not set, default to Outage Mode")
	flag.IntVar(&c.WaitTimeForMakeRpDown, "wait_time_for_make_rp_down", 5, "Wait time for make rp down, if not set, default to 5")

	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}

	// Keep a more frequent flush frequency as 1 second
	klog.StartFlushDaemon(time.Second * 1)

	defer klog.Flush()
	klog.Info("")
	klog.Infof("Region resource manager simulator config / region name:    (No.%v)", c.RegionName)
	klog.Infof("Region resource manager simulator config / rp number per region: (%v)", c.RpNum)
	klog.Infof("Region resource manager simulator config / node number per rp:  (%v)", c.NodesPerRP)
	klog.Infof("Region resource manager simulator config / simulator data pattern:  (%v)", c.DataPattern)
	klog.Infof("Region resource manager simulator config / wait time for make rp down:  (%v)", c.WaitTimeForMakeRpDown)

	klog.Info("")
	klog.Infof("Starting resource region manager simulator (%v)", c.RegionName)
	klog.Info("")

	// Initialize Added Event List and Modified Event List
	// Region node Added Event List - for initpull
	//data.Init(c.RegionName, c.DataPattern, c.RpNum, c.NodesPerRP)
	data.Init(c.RegionName, c.RpNum, c.NodesPerRP)

	// Generate update changes of Default Pattern
	// - simulate random RP down
	// OR
	// Generate update changes of Daily Pattern
	// - simulate 10 changes each minute
	data.MakeDataUpdate(c.DataPattern, c.WaitTimeForMakeRpDown)

	// Run simulater RSET API server
	if err := app.Run(c); err != nil {
		klog.Errorf("Error: %v\n", err)
	}

	klog.Infof("Exiting reesource management service")
}

// function to print the usage info for the resource management api server
func printUsage() {
	fmt.Println("\nUsage: Region Resource Manager Simulator")
	fmt.Println("\n       Per region config options: --region_name=<region name>  --rp_num=<number of rp>  --nodes_per_rp=<number of nodes> --master_port=<port> --data_pattern=<outage or daily> --wait_time_for_make_rp_down=<number of minutes>")
	fmt.Println()

	os.Exit(0)
}
