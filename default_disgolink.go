package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
)

var _ Disgolink = (*defaultDisgolink)(nil)
var _ Lavalink = (*defaultDisgolink)(nil)
var _ core.VoiceDispatchInterceptor = (*defaultDisgolink)(nil)
var _ core.EventListener = (*defaultDisgolink)(nil)

type defaultDisgolink struct {
	Lavalink
}

func (l *defaultDisgolink) OnVoiceServerUpdate(voiceServerUpdate *core.VoiceServerUpdateEvent) {
	l.Lavalink.VoiceServerUpdate(&VoiceServerUpdate{
		Token:    voiceServerUpdate.Token,
		GuildID:  voiceServerUpdate.GuildID,
		Endpoint: voiceServerUpdate.Endpoint,
	})
}

func (l *defaultDisgolink) OnVoiceStateUpdate(voiceStateUpdate *core.VoiceStateUpdateEvent) {
	l.Lavalink.VoiceStateUpdate(&VoiceStateUpdate{
		GuildID:   voiceStateUpdate.GuildID,
		ChannelID: voiceStateUpdate.ChannelID,
		UserID:    voiceStateUpdate.UserID,
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *defaultDisgolink) OnEvent(event interface{}) {

}
