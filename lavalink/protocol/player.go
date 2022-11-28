package protocol

import (
	"time"

	"github.com/disgoorg/json"
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
	Connected bool   `json:"connected"`
	Ping      int    `json:"ping"`
}

type PlayerState struct {
	Time      time.Time `json:"time"`
	Position  Duration  `json:"position"`
	Connected bool      `json:"connected"`
	Ping      int       `json:"ping"`
}

type PlayerUpdate struct {
	EncodedTrack *json.Nullable[string] `json:"encodedTrack,omitempty"`
	Identifier   *string                `json:"identifier,omitempty"`
	Position     *int                   `json:"position,omitempty"`
	EndTime      *int                   `json:"endTime,omitempty"`
	Volume       *int                   `json:"volume,omitempty"`
	Paused       *bool                  `json:"paused,omitempty"`
	Voice        *VoiceState            `json:"voice,omitempty"`
	Filters      *Filters               `json:"filters,omitempty"`
}
