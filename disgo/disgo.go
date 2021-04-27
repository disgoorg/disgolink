package disgolink

import (
	"fmt"

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

func NewDisgolinkByUserID(logger log.Logger, userID dapi.Snowflake) Disgolink {
	return &DisgolinkImpl{
		Lavalink: internal.NewLavalinkImpl(logger, string(userID)),
	}
}

func NewDisgolink(logger log.Logger, dgo dapi.Disgo) Disgolink {
	dgolink := NewDisgolinkByUserID(logger, dgo.ApplicationID())
	dgo.EventManager().AddEventListeners(dgolink)
	dgo.SetVoiceDispatchInterceptor(dgolink)
	return dgolink
}

type DisgolinkImpl struct {
	api.Lavalink
}

func (l *DisgolinkImpl) OnVoiceServerUpdate(voiceServerUpdate *dapi.VoiceServerUpdateEvent) {
	fmt.Printf("voiceServerUpdate: %+v", voiceServerUpdate)
	l.Lavalink.VoiceServerUpdate(&api.VoiceServerUpdate{
		Token:    voiceServerUpdate.Token,
		GuildID:  string(voiceServerUpdate.GuildID),
		Endpoint: voiceServerUpdate.Endpoint,
	})
}

func (l *DisgolinkImpl) OnVoiceStateUpdate(voiceStateUpdate *dapi.VoiceStateUpdateEvent) {
	l.Lavalink.VoiceStateUpdate(&api.VoiceStateUpdate{
		GuildID:   voiceStateUpdate.GuildID.String(),
		ChannelID: (*string)(voiceStateUpdate.ChannelID),
		UserID:    voiceStateUpdate.UserID.String(),
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *DisgolinkImpl) OnEvent(event interface{}) {

}

