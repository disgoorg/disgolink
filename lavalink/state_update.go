package lavalink

import "github.com/DisgoOrg/snowflake"

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
