package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
)

type VoiceServerUpdate struct {
	Token    string            `json:"token"`
	GuildID  discord.Snowflake `json:"guildId"`
	Endpoint *string           `json:"endpoint"`
}

type VoiceStateUpdate struct {
	GuildID   discord.Snowflake  `json:"guild_id"`
	ChannelID *discord.Snowflake `json:"channel_id"`
	SessionID string             `json:"session_id"`
}
