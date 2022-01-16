package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgolink/lavalink"
)

func New(disgo *core.Bot, opts ...lavalink.ConfigOpt) Link {
	opts = append(opts, lavalink.WithLogger(disgo.Logger),
		lavalink.WithHTTPClient(disgo.RestServices.HTTPClient()),
		lavalink.WithUserID(disgo.ApplicationID.String()))
	dgolink := &linkImpl{
		Lavalink: lavalink.New(opts...),
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
		l.VoiceServerUpdate(lavalink.VoiceServerUpdate{
			Token:    e.VoiceServerUpdate.Token,
			GuildID:  e.VoiceServerUpdate.GuildID.String(),
			Endpoint: e.VoiceServerUpdate.Endpoint,
		})

	case *events.GuildVoiceStateUpdateEvent:
		if e.VoiceState.UserID.String() != l.UserID() {
			return
		}
		l.VoiceStateUpdate(lavalink.VoiceStateUpdate{
			GuildID:   e.VoiceState.GuildID.String(),
			ChannelID: (*string)(e.VoiceState.ChannelID),
			SessionID: e.VoiceState.SessionID,
		})
	}
}
