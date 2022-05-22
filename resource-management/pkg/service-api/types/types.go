package types

import "global-resource-service/resource-management/pkg/common-lib/types"

// The request content of the Watch API call
// ResourceVersionMap is part of the return of the LIST API call
type WatchRequest struct {
	ResourceVersions types.ResourceVersionMap `json:"resource_versions"`
}
