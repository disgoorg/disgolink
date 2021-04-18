package api

import "github.com/DisgoOrg/disgolink/api/filters"

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

type GenericOpCommand struct {
	Op      Op     `json:"op"`
	GuildID string `json:"guildId"`
}

func NewCommand(op Op) *GenericOpCommand {
	return &GenericOpCommand{Op: op}
}

type EventCommand struct {
	*GenericOpCommand
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
}

type PlayerCommand struct {
	*GenericOpCommand
	GuildID string `json:"guildId"`
}

func NewPlayerCommand(op Op, p Player) *PlayerCommand {
	return &PlayerCommand{
		GenericOpCommand: NewCommand(op),
		GuildID:          p.GuildID(),
	}
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

type FilterPlayerCommand struct {
	*PlayerCommand
	*filters.Filters
}

type GenericOpEvent struct {
	Op      Op     `json:"op"`
}

type PlayerUpdateEvent struct {
	GenericOpEvent
}
