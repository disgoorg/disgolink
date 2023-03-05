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
	State   PlayerState  `json:"state"`
	Voice   VoiceState   `json:"voice"`
	Filters Filters      `json:"filters"`
}

type VoiceState struct {
	Token     string `json:"token"`
	Endpoint  string `json:"endpoint"`
	SessionID string `json:"sessionId"`
}

type PlayerState struct {
	Time      Timestamp `json:"time"`
	Position  Duration  `json:"position"`
	Connected bool      `json:"connected"`
	Ping      int       `json:"ping"`
}
