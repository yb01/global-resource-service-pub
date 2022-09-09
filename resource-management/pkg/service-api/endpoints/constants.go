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

package endpoints

// URL path
const (
	ClientAdminitrationPath = "/clients"

	//RegionlessResourcePath is the default api service url
	RegionlessResourcePath = "/resource"

	// TODO revisit and evaluate API paths.
	NodeStatusPath = "/nodes"

	ListWatchResourcePath = RegionlessResourcePath + "/{clientid}"
	UpdateResourcePath    = RegionlessResourcePath + "/{clientid}" + "/addResource"
	ReduceResourcePath    = RegionlessResourcePath + "/{clientid}" + "/reduceResource"
	// InsecureServiceAPIPort is the default port for Service-api when running insecure mode.
	// TODO: Can be overridden by a flag at startup.
	InsecureServiceAPIPort = "8080"
	// SecureServiceAPIPort is the default port for Service-api when running secure mode.
	// TODO: Can be overridden by a flag at startup.
	SecureServiceAPIPort = "443"

	WatchChannelSize         = 100
	WatchParameter           = "watch"
	WatchParameterTrue       = "true"
	ListLimitParameter       = "limit"
	DefaultResponseTrunkSize = 500
)
