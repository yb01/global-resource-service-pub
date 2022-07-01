package types

import "global-resource-service/resource-management/pkg/common-lib/types"

// WatchRequest is the request body of the Watch API call
// ResourceVersionMap is part of the return of the LIST API call
type WatchRequest struct {
	ResourceVersions types.TransitResourceVersionMap `json:"resource_versions"`
}

// ClientRegistrationRequest is the request body when a client register to the resource management service
// TBD: Optionally, client can set its customized name and initial resource request
type ClientRegistrationRequest struct {
	ClientInfo               types.ClientInfoType  `json:"client_info,omitempty"`
	InitialRequestedResource types.ResourceRequest `json:"init_resource_request,omitempty"`
}

// ClientRegistrationResponse is the response body for approved client registration request
// ClientId is required for an approved client registration to the resource management service
// GrantedResource is an info to client on the resource level the List OP it can request
type ClientRegistrationResponse struct {
	ClientId        string                `json:"client_id"`
	GrantedResource types.ResourceRequest `json:"granted_resource,omitempty"`
}

// ListNodeResponse is the response body for listing nodes from a client
// NodeList is the list of LogicalNodes returned from Distributor allocated for this client
// ResourceVersions are the list of RVs from each RP
type ListNodeResponse struct {
	NodeList         []*types.LogicalNode     `json:"node_list",omitempty`
	ResourceVersions types.TransitResourceVersionMap `json:"resource_version_map,omitempty"`
}
