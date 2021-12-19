package lavalink

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

func NewDisgolink(disgo *core.Bot) Disgolink {
	dgolink := &disgolinkImpl{
		Lavalink: NewLavalink(
			WithLogger(disgo.Logger),
			WithHTTPClient(disgo.RestServices.HTTPClient()),
			WithUserID(disgo.ApplicationID),
		),
	}

	disgo.EventManager.AddEventListeners(dgolink)
	return dgolink
}

type Disgolink interface {
	Lavalink
	core.EventListener
}

var (
	_ Disgolink          = (*disgolinkImpl)(nil)
	_ Lavalink           = (*disgolinkImpl)(nil)
	_ core.EventListener = (*disgolinkImpl)(nil)
)

type disgolinkImpl struct {
	Lavalink
}

func (l *disgolinkImpl) OnEvent(event core.Event) {
	switch e := event.(type) {
	case *events.VoiceServerUpdateEvent:
		l.VoiceServerUpdate(VoiceServerUpdate{
			Token:    e.VoiceServerUpdate.Token,
			GuildID:  e.VoiceServerUpdate.GuildID,
			Endpoint: e.VoiceServerUpdate.Endpoint,
		})

	case *events.GuildVoiceStateUpdateEvent:
		if e.VoiceState.UserID != l.UserID() {
			return
		}
		l.VoiceStateUpdate(VoiceStateUpdate{
			GuildID:   e.VoiceState.GuildID,
			ChannelID: e.VoiceState.ChannelID,
			SessionID: e.VoiceState.SessionID,
		})
	}
}
