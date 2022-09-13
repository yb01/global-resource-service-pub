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
	"net/http"
	"time"

	distributor "global-resource-service/resource-management/pkg/common-lib/interfaces/distributor"
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/event"
)

type Aggregator struct {
	urls           []string
	EventProcessor distributor.Interface
}

// To be client of Resource Region Manager
// RRM: Resource Region Manager
//
type ClientOfRRM struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RRM: Resource Region Manager
//
type ResponseFromRRM struct {
	RegionNodeEvents [][]*event.NodeEvent
	RvMap            types.TransitResourceVersionMap
	Length           uint64
}

// RRM: Resource Region Manager
//
type PullDataFromRRM struct {
	BatchLength uint64
	DefaultCRV  uint64
	CRV         types.TransitResourceVersionMap
}

const (
	DefaultBatchLength  = 20000
	httpPrefix          = "http://"
	defaultPullInterval = 10 * time.Millisecond // 10ms as default pull interval
)

// Initialize aggregator
//
func NewAggregator(urls []string, EventProcessor distributor.Interface) *Aggregator {
	return &Aggregator{
		urls:           urls,
		EventProcessor: EventProcessor,
	}
}
