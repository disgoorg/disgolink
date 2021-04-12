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


func (l *DisgolinkImpl) OnVoiceServerUpdate(voiceServerUpdateEvent *dapi.VoiceServerUpdateEvent) {

}

func (l *DisgolinkImpl) OnVoiceStateUpdate(voiceStateUpdateEvent *dapi.VoiceStateUpdateEvent) {

}

func (l *DisgolinkImpl) OnEvent(event interface{}) {

}
