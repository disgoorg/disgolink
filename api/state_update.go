package api

import dapi "github.com/DisgoOrg/disgo/api"

type VoiceServerUpdate struct {
	Token    string         `json:"token"`
	GuildID  dapi.Snowflake `json:"guildId"`
	Endpoint *string        `json:"endpoint"`
}

type VoiceStateUpdate struct {
	GuildID   dapi.Snowflake  `json:"guild_id"`
	ChannelID *dapi.Snowflake `json:"channel_id"`
	UserID    dapi.Snowflake  `json:"user_id"`
	SessionID string          `json:"session_id"`
}
