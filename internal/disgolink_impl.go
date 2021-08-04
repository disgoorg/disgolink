package internal

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
)

var _ api.Disgolink = (*DisgolinkImpl)(nil)
var _ api.Lavalink = (*DisgolinkImpl)(nil)
var _ dapi.VoiceDispatchInterceptor = (*DisgolinkImpl)(nil)
var _ dapi.EventListener = (*DisgolinkImpl)(nil)

type DisgolinkImpl struct {
	api.Lavalink
}

func (l *DisgolinkImpl) OnVoiceServerUpdate(voiceServerUpdate *dapi.VoiceServerUpdateEvent) {
	l.Lavalink.VoiceServerUpdate(&api.VoiceServerUpdate{
		Token:    voiceServerUpdate.Token,
		GuildID:  voiceServerUpdate.GuildID,
		Endpoint: voiceServerUpdate.Endpoint,
	})
}

func (l *DisgolinkImpl) OnVoiceStateUpdate(voiceStateUpdate *dapi.VoiceStateUpdateEvent) {
	l.Lavalink.VoiceStateUpdate(&api.VoiceStateUpdate{
		GuildID:   voiceStateUpdate.GuildID,
		ChannelID: voiceStateUpdate.ChannelID,
		UserID:    voiceStateUpdate.UserID,
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *DisgolinkImpl) OnEvent(event interface{}) {

}
