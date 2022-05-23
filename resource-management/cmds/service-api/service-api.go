package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"time"

	"global-resource-service/resource-management/cmds/service-api/app"
)

func main() {
	flag.Usage = printUsage

	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}

	// keep a more frequent flush frequency as 1 second
	klog.StartFlushDaemon(time.Second * 1)

	defer klog.Flush()

	klog.Infof("Starting resource management service")

	if err := app.Run(); err != nil {
		klog.Errorf("error: %v\n", err)
	}

	klog.Infof("Exiting reesource management service")
}

// function to print the usage info for the resource management api server
func printUsage() {
	fmt.Println("Usage: ")

	// klog will use commandline log parameters with nil as named a few below:
	// --alsologtostderr=true  --logtostderr=false --log_file="/tmp/grs.log"
	fmt.Println("logging options: --alsologtostderr=true  --logtostderr=false --log_file=/tmp/grs.log ")
	os.Exit(0)
}
