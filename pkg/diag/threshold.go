package diag

const (
	/*
		thresholds about process
	*/
	ThresholdWaitNumber = 10
	ThresholdRunNumber  = 0.6

	/*
		thresholds about memory
	*/
	ThresholdMemoryFree = 0.2

	/*
		thresholds about cpu
	*/
	ThresholdIdleTime    = 0.9
	ThresholdWaitPercent = 0.2

	/*
		thresholds about graph service
	*/

	ThresholdIdleNumQuerisSum600 = 10000

	/*
		thresholds about meta service
	*/

	ThresholdIdleHeartbeatLatencyUsAvg600 = 900


	/*
		thresholds about storage service
	*/

	ThresholdIdleNumLookupErrorsSum600 = 100
)
