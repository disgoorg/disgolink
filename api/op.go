package api

type OpType string

const (
	PlayOp              OpType = "play"
	StopOp              OpType = "stop"
	PauseOP             OpType = "pause"
	SeekOP              OpType = "seek"
	VolumeOP            OpType = "volume"
	DestroyOp           OpType = "destroy"
	StatsOp             OpType = "stats"
	PlayerUpdateOp      OpType = "playerUpdate"
	EventOp             OpType = "event"
	ConfigureResumingOp OpType = "configureResuming"
	FiltersOp           OpType = "filters"
)
