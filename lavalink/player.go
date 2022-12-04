package lavalink

import (
	"github.com/disgoorg/snowflake/v2"
)

type Players []Player

type Player struct {
	GuildID snowflake.ID `json:"guildId"`
	Track   *Track       `json:"track"`
	Volume  int          `json:"volume"`
	Paused  bool         `json:"paused"`
	Voice   VoiceState   `json:"voice"`
	Filters Filters      `json:"filters"`
}

type VoiceState struct {
	Token     string `json:"token"`
	Endpoint  string `json:"endpoint"`
	SessionID string `json:"sessionId"`
	Connected bool   `json:"connected,omitempty"`
	Ping      int    `json:"ping,omitempty"`
}

type PlayerState struct {
	Time      Time     `json:"time"`
	Position  Duration `json:"position"`
	Connected bool     `json:"connected"`
	Ping      int      `json:"ping"`
}
