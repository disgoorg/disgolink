package api

type VoiceServer struct {
	Token    string    `json:"token"`
	GuildID  Snowflake `json:"guild_id"`
	Endpoint *string   `json:"endpoint"`
}

type VoiceState struct {
	GuildID   Snowflake  `json:"guild_id"`
	ChannelID *Snowflake `json:"channel_id"`
	UserID    Snowflake  `json:"user_id"`
	SessionID string     `json:"session_id"`
}