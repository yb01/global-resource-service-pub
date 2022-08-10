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

package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"testing"
)

func TestResourceVersionMap_Marshall_UnMarshall(t *testing.T) {
	rvs := make(TransitResourceVersionMap)
	loc := RvLocation{Region: location.Beijing, Partition: location.ResourcePartition1}
	rvs[loc] = 100

	// marshall
	b, err := json.Marshal(rvs)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	// unmarshall
	var newRVMap TransitResourceVersionMap
	err = json.Unmarshal(b, &newRVMap)
	assert.Nil(t, err)
	assert.NotNil(t, newRVMap)
	assert.Equal(t, 1, len(newRVMap))
	assert.Equal(t, uint64(100), newRVMap[loc])
}
