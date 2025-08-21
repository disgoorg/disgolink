package lavalink

import (
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

type Players []Player

type Player struct {
	GuildID snowflake.ID `json:"guildId"`
	Track   *Track       `json:"track"`
	Volume  int          `json:"volume"`
	Paused  bool         `json:"paused"`
	State   PlayerState  `json:"state"`
	Voice   VoiceState   `json:"voice"`
	Filters Filters      `json:"filters"`
}

type VoiceState struct {
	Token     string       `json:"token"`
	Endpoint  string       `json:"endpoint"`
	SessionID string       `json:"sessionId"`
	ChannelID snowflake.ID `json:"channelId"`
}

type PlayerState struct {
	Time      Timestamp `json:"time"`
	Position  Duration  `json:"position"`
	Connected bool      `json:"connected"`
	Ping      int       `json:"ping"`
}

type PlayerUpdateTrack struct {
	Encoded    omit.Omit[*string] `json:"encoded,omitzero"`
	Identifier *string            `json:"identifier,omitzero"`
	UserData   any                `json:"userData,omitzero"`
}

type PlayerUpdate struct {
	Track     *PlayerUpdateTrack `json:"track,omitzero"`
	Position  *Duration          `json:"position,omitzero"`
	EndTime   *Duration          `json:"endTime,omitzero"`
	Volume    *int               `json:"volume,omitzero"`
	Paused    *bool              `json:"paused,omitzero"`
	Voice     *VoiceState        `json:"voice,omitzero"`
	Filters   *Filters           `json:"filters,omitzero"`
	NoReplace bool               `json:"-"`
}
