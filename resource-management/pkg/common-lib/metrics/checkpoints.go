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

package metrics

type ResourceManagementCheckpoint int

const (
	Aggregator_Received ResourceManagementCheckpoint = 0

	Distributor_Received ResourceManagementCheckpoint = 1
	Distributor_Sending  ResourceManagementCheckpoint = 2
	Distributor_Sent     ResourceManagementCheckpoint = 3

	Serializer_Encoded ResourceManagementCheckpoint = 4
	Serializer_Sent    ResourceManagementCheckpoint = 5

	Len_ResourceManagementCheckpoint = 6
)

type ResourceManagementCheckpointName string

const (
	Aggregator_Received_Name ResourceManagementCheckpointName = "AGG_RECEIVED"

	Distributor_Received_Name ResourceManagementCheckpointName = "DIS_RECEIVED"
	Distributor_Sending_Name  ResourceManagementCheckpointName = "DIS_SENDING"
	Distributor_Sent_Name     ResourceManagementCheckpointName = "DIS_SENT"

	Serializer_Encoded_Name ResourceManagementCheckpointName = "SER_ENCODED"
	Serializer_Sent_Name    ResourceManagementCheckpointName = "SER_SENT"
)
