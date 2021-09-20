package disgolink

type Op string

const (
	OpPlay              Op = "play"
	OpStop              Op = "stop"
	OpPause             Op = "pause"
	OpSeek              Op = "seek"
	OpVolume            Op = "volume"
	OpEqualizer         Op = "equalizer"
	OpDestroy           Op = "destroy"
	OpStats             Op = "stats"
	OpVoiceUpdate       Op = "voiceUpdate"
	OpPlayerUpdate      Op = "playerUpdate"
	OpEvent             Op = "event"
	OpConfigureResuming Op = "configureResuming"
	OpFilters           Op = "filters"
)
