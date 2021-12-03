package lavalink

type OpType string

const (
	OpTypePlay              OpType = "play"
	OpTypeStop              OpType = "stop"
	OpTypePause             OpType = "pause"
	OpTypeSeek              OpType = "seek"
	OpTypeVolume            OpType = "volume"
	OpTypeDestroy           OpType = "destroy"
	OpTypeStats             OpType = "stats"
	OpTypeVoiceUpdate       OpType = "voiceUpdate"
	OpTypePlayerUpdate      OpType = "playerUpdate"
	OpTypeEvent             OpType = "event"
	OpTypeConfigureResuming OpType = "configureResuming"
	OpTypeFilters           OpType = "filters"
)
