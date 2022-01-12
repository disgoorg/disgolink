package lavalink

type VoiceServerUpdate struct {
	Token    string  `json:"token"`
	GuildID  string  `json:"guildId"`
	Endpoint *string `json:"endpoint"`
}

type VoiceStateUpdate struct {
	GuildID   string  `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SessionID string  `json:"session_id"`
}
