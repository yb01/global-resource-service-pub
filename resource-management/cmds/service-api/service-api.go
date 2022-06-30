package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"time"

	"global-resource-service/resource-management/cmds/service-api/app"
)

func main() {
	flag.Usage = printUsage

	// get the commandline arguments
	c := &app.Config{}

	var urls string
	flag.StringVar(&c.MasterIp, "master_ip", "localhost", "Service IP address, if not set, default to localhost")
	flag.StringVar(&c.MasterPort, "master_port", "8080", "Service port, if not set, default to 8080")
	flag.StringVar(&urls, "resource_urls", "", "Resource urls of the resource manager services in each region")

	if !flag.Parsed() {
		klog.InitFlags(nil)
		flag.Parse()
	}

	if urls == "" {
		klog.Errorf("resource_urls is missing or with invalid values")
		os.Exit(1)
	}

	c.ResourceUrls = strings.Split(urls, ",")

	// keep a more frequent flush frequency as 1 second
	klog.StartFlushDaemon(time.Second * 1)

	defer klog.Flush()

	klog.Infof("Service config: %v", c)

	klog.Infof("Starting resource management service")

	if err := app.Run(c); err != nil {
		klog.Errorf("error: %v\n", err)
	}

	klog.Infof("Exiting resource management service")
}

// function to print the usage info for the resource management api server
func printUsage() {
	fmt.Println("Usage: ")

	// klog will use commandline log parameters with nil as named a few below:
	// --alsologtostderr=true  --logtostderr=false --log_file="/tmp/grs.log"
	fmt.Println("logging options: --alsologtostderr=true  --logtostderr=false --log_file=/tmp/grs.log ")
	fmt.Println("service config options: --master_ip=<master address>  --master_port=<port> --resource_urls=<url1,url2,...>")
	fmt.Println("Explanation: <master address> could be public ip address or public dns name of the server")

	os.Exit(0)
}
