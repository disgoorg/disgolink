package api

type VoiceServer struct {
	Token    string  `json:"token"`
	GuildId  string  `json:"guild_id"`
	Endpoint *string `json:"endpoint"`
}
