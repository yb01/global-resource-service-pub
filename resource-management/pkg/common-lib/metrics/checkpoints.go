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

var ResourceManagementMeasurement_Enabled = true

func SetEnableResourceManagementMeasurement(enabled bool) {
	ResourceManagementMeasurement_Enabled = enabled
}
