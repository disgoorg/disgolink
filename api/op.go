package api

type OpType string

const (
	PlayOp              OpType = "play"
	StopOp              OpType = "stop"
	PauseOP             OpType = "pause"
	SeekOP              OpType = "seek"
	VolumeOP            OpType = "volume"
	EqualizerOP         OpType = "equalizer"
	DestroyOp           OpType = "destroy"
	StatsOp             OpType = "stats"
	VoiceUpdateOp       OpType = "voiceUpdate"
	PlayerUpdateOp      OpType = "playerUpdate"
	EventOp             OpType = "event"
	ConfigureResumingOp OpType = "configureResuming"
	FiltersOp           OpType = "filters"
)

func NewCommand(op OpType) *OpCommand {
	return &OpCommand{Op: op}
}

type OpCommand struct {
	Op      OpType `json:"op"`
	GuildID string `json:"guildId"`
}

func NewPlayerCommand(op OpType, guildID string) *OpPlayerCommand {
	return &OpPlayerCommand{
		OpCommand: NewCommand(op),
		GuildID:   guildID,
	}
}

type OpPlayerCommand struct {
	*OpCommand
	GuildID string `json:"guildId"`
}

type EventCommand struct {
	*OpCommand
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
}

type OpPlayPlayer struct {
	*OpPlayerCommand
	Track     string `json:"track"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
	NoReplace bool   `json:"noReplace"`
	Paused    bool   `json:"pause"`
}

type OpDestroyPlayer struct {
	*OpPlayerCommand
}

type OpStopPlayer struct {
	*OpPlayerCommand
}

type OpPausePlayer struct {
	*OpPlayerCommand
	Paused bool `json:"pause"`
}

type OpSeekPlayer struct {
	*OpPlayerCommand
	Position int `json:"position"`
}
