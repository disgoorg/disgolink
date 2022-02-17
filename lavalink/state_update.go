package lavalink

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/snowflake"
)

type VoiceServerUpdate struct {
	Token    string              `json:"token"`
	GuildID  snowflake.Snowflake `json:"guildId"`
	Endpoint *string             `json:"endpoint"`
}

type VoiceStateUpdate struct {
	GuildID   snowflake.Snowflake  `json:"guild_id"`
	ChannelID *snowflake.Snowflake `json:"channel_id"`
	SessionID string               `json:"session_id"`
}

type PlayerState struct {
	Time      time.Time `json:"time"`
	Position  Duration  `json:"position"`
	Connected bool      `json:"connected"`
}

func (s *PlayerState) UnmarshalJSON(data []byte) error {
	type playerState PlayerState
	var v struct {
		Time int64 `json:"time"`
		playerState
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*s = PlayerState(v.playerState)
	s.Time = time.UnixMilli(v.Time)
	return nil
}
