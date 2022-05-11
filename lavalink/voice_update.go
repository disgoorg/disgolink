package lavalink

import (
	"github.com/disgoorg/snowflake/v2"
)

type VoiceServerUpdate struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guildId"`
	Endpoint *string      `json:"endpoint"`
}

type VoiceStateUpdate struct {
	GuildID   snowflake.ID  `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	SessionID string        `json:"session_id"`
}
