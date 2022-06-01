package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
	"testing"
)

func TestResourceVersionMap_Marshall_UnMarshall(t *testing.T) {
	rvs := make(ResourceVersionMap)
	loc := location.NewLocation(location.Beijing, location.ResourcePartition1)
	rvs[*loc] = 100

	// marshall
	b, err := json.Marshal(rvs)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	// unmarshall
	var newRVMap ResourceVersionMap
	err = json.Unmarshal(b, &newRVMap)
	assert.Nil(t, err)
	assert.NotNil(t, newRVMap)
	assert.Equal(t, 1, len(newRVMap))
	assert.Equal(t, uint64(100), newRVMap[*loc])
}
