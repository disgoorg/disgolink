package disgolink

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
)

type Disgolink interface {
	api.Lavalink
	dapi.VoiceDispatchInterceptor
	dapi.EventListener
}

func NewDisgolink(userID api.Snowflake) Disgolink {
	return &DisgolinkImpl{
		Lavalink: internal.NewLavalinkImpl(userID),
	}
}

type DisgolinkImpl struct {
	api.Lavalink
	dapi.VoiceDispatchInterceptor
}

func (l *DisgolinkImpl) OnVoiceServerUpdate(voiceServerUpdate *dapi.VoiceServerUpdateEvent) {
	l.Lavalink.VoiceServerUpdate(&api.VoiceServerUpdate{
		Token:    voiceServerUpdate.Token,
		GuildID: api.Snowflake(voiceServerUpdate.GuildID),
		Endpoint: voiceServerUpdate.Endpoint,
	})
}

func (l *DisgolinkImpl) OnVoiceStateUpdate(voiceStateUpdate *dapi.VoiceStateUpdateEvent) {
	l.Lavalink.VoiceStateUpdate(&api.VoiceStateUpdate{
		GuildID:   api.Snowflake(voiceStateUpdate.GuildID),
		ChannelID: (*api.Snowflake)(voiceStateUpdate.ChannelID),
		UserID:    api.Snowflake(voiceStateUpdate.UserID),
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *DisgolinkImpl) OnEvent(event interface{}) {

}
