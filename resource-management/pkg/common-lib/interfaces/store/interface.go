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

package store

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
	"global-resource-service/resource-management/pkg/common-lib/types/location"
)

const (
	Preserve_VirtualNodesAssignments_KeyPrefix = "VirtualNodesAssignments"
	Preserve_NodeStoreStatus_KeyPrefix         = "NodeStoreStatus"
	Preserve_Client_KeyPrefix                  = "Client"
)

type StoreInterface interface {
	PersistNodes([]*types.LogicalNode) bool
	PersistNodeStoreStatus(*NodeStoreStatus) bool
	PersistVirtualNodesAssignments(*VirtualNodeAssignment) bool

	// Interfaces for client object operations to store
	PersistClient(string, *types.Client) error
	GetClient(string) (*types.Client, error)
	// Get all client object, during distributor restart, Admin UI etc
	GetClients() ([]*types.Client, error)
	// UpdateClient will be used with client Add/remove resources
	UpdateClient(string, *types.Client) error

	// For fake storage test only, no need to implement
	InitNodeIdCache()
	GetNodeIdCount() int
	SetTestNodeIdMatch(isMatch bool)
}

type NodeStoreStatus struct {
	// # of regions
	RegionNum int

	// # of max resource partition in each region
	PartitionMaxNum int

	// virutal node number per resource partition
	VirtualNodeNumPerRP int

	// Latest resource version map
	CurrentResourceVerions types.TransitResourceVersionMap
}

func (nsStatus *NodeStoreStatus) GetKey() string {
	return Preserve_NodeStoreStatus_KeyPrefix
}

type VirtualNodeAssignment struct {
	ClientId     string
	VirtualNodes []*VirtualNodeConfig
}

func (assignment *VirtualNodeAssignment) GetKey() string {
	return Preserve_VirtualNodesAssignments_KeyPrefix
}

type VirtualNodeConfig struct {
	Lowerbound float64
	Upperbound float64
	Location   location.Location
}
