package disgolink

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
	"github.com/DisgoOrg/log"
)

var _ Disgolink = (*DisgolinkImpl)(nil)
var _ api.Lavalink = (*DisgolinkImpl)(nil)
var _ dapi.VoiceDispatchInterceptor = (*DisgolinkImpl)(nil)
var _ dapi.EventListener = (*DisgolinkImpl)(nil)

type Disgolink interface {
	api.Lavalink
	dapi.VoiceDispatchInterceptor
	dapi.EventListener
}

func NewDisgolink(logger log.Logger, userID dapi.Snowflake) Disgolink {
	return &DisgolinkImpl{
		LavalinkImpl: internal.NewLavalinkImpl(logger, api.Snowflake(userID)),
	}
}

type DisgolinkImpl struct {
	*internal.LavalinkImpl
}

func (l *DisgolinkImpl) OnVoiceServerUpdate(voiceServerUpdate *dapi.VoiceServerUpdateEvent) {
	l.LavalinkImpl.VoiceServerUpdate(&api.VoiceServerUpdate{
		Token:    voiceServerUpdate.Token,
		GuildID:  api.Snowflake(voiceServerUpdate.GuildID),
		Endpoint: voiceServerUpdate.Endpoint,
	})
}

func (l *DisgolinkImpl) OnVoiceStateUpdate(voiceStateUpdate *dapi.VoiceStateUpdateEvent) {
	l.LavalinkImpl.VoiceStateUpdate(&api.VoiceStateUpdate{
		GuildID:   api.Snowflake(voiceStateUpdate.GuildID),
		ChannelID: (*api.Snowflake)(voiceStateUpdate.ChannelID),
		UserID:    api.Snowflake(voiceStateUpdate.UserID),
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *DisgolinkImpl) OnEvent(event interface{}) {

}
