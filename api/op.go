package api

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

func NewCommand(op Op) *OpCommand {
	return &OpCommand{Op: op}
}

type OpCommand struct {
	Op      Op     `json:"op"`
	GuildID string `json:"guildId"`
}

func NewPlayerCommand(op Op, p Player) *PlayerCommand {
	return &PlayerCommand{
		OpCommand: NewCommand(op),
		GuildID:   p.GuildID(),
	}
}

type EventCommand struct {
	*OpCommand
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
}

type PlayerCommand struct {
	*OpCommand
	GuildID string `json:"guildId"`
}

type PlayPlayerCommand struct {
	*PlayerCommand
	Track     string `json:"track"`
	StartTime int    `json:"startTime,omitempty"`
	EndTime   int    `json:"endTime,omitempty"`
	NoReplace bool   `json:"noReplace"`
	Paused    bool   `json:"pause"`
}

type DestroyPlayerCommand struct {
	*PlayerCommand
}

type StopPlayerCommand struct {
	*PlayerCommand
}

type PausePlayerCommand struct {
	*PlayerCommand
	Paused bool `json:"pause"`
}

type SeekPlayerCommand struct {
	*PlayerCommand
	Position int `json:"position"`
}
