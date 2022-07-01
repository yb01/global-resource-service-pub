package metrics

type ResourceManagementCheckpoint string

const (
	Aggregator_Received  ResourceManagementCheckpoint = "AGG_RECEIVED"
	Distributor_Received ResourceManagementCheckpoint = "DIS_RECEIVED"
	Distributor_Sending  ResourceManagementCheckpoint = "DIS_SENDING"
	Distributor_Sent     ResourceManagementCheckpoint = "DIS_SENT"
	Serializer_Sent      ResourceManagementCheckpoint = "SER_SENT"
)

var ResourceManagementMeasurement_Enabled = true
