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

package runtime

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"time"
)

type Object interface {
	GetResourceVersionInt64() uint64
	GetGeoInfo() types.NodeGeoInfo
	GetEventType() EventType
	GetId() string
	GetLocation() *location.Location
	GetLastUpdatedTime() time.Time

	SetCheckpoint(int)
	GetCheckpoints() []time.Time

	// Used to remove wrapper
	GetEvent() Object
}
