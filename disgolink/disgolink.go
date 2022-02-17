package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgolink/lavalink"
)

func New(disgo *core.Bot, opts ...lavalink.ConfigOpt) Link {
	dgolink := &linkImpl{
		Lavalink: lavalink.New(append([]lavalink.ConfigOpt{
			lavalink.WithLogger(disgo.Logger),
			lavalink.WithHTTPClient(disgo.RestServices.HTTPClient()),
			lavalink.WithUserID(disgo.ApplicationID),
		}, opts...)...),
	}

	disgo.EventManager.AddEventListeners(dgolink)
	return dgolink
}

type Link interface {
	lavalink.Lavalink
	core.EventListener
}

var (
	_ Link               = (*linkImpl)(nil)
	_ lavalink.Lavalink  = (*linkImpl)(nil)
	_ core.EventListener = (*linkImpl)(nil)
)

type linkImpl struct {
	lavalink.Lavalink
}

func (l *linkImpl) OnEvent(event core.Event) {
	switch e := event.(type) {
	case *events.VoiceServerUpdateEvent:
		l.OnVoiceServerUpdate(lavalink.VoiceServerUpdate{
			Token:    e.VoiceServerUpdate.Token,
			GuildID:  e.VoiceServerUpdate.GuildID,
			Endpoint: e.VoiceServerUpdate.Endpoint,
		})

	case *events.GuildVoiceStateUpdateEvent:
		if e.VoiceState.UserID != l.UserID() {
			return
		}
		l.OnVoiceStateUpdate(lavalink.VoiceStateUpdate{
			GuildID:   e.VoiceState.GuildID,
			ChannelID: e.VoiceState.ChannelID,
			SessionID: e.VoiceState.SessionID,
		})
	}
}
