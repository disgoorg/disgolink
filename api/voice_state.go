package api

type VoiceState struct {
	GuildId   string  `json:"guild_id"`
	ChannelId *string `json:"channel_id"`
	UserId    string  `json:"user_id"`
	SessionId string  `json:"session_id"`
}
