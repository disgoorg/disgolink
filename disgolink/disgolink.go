package disgolink

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/lavalink"
)

func New(client bot.Client, opts ...lavalink.ConfigOpt) Link {
	dgolink := &linkImpl{
		Lavalink: lavalink.New(append([]lavalink.ConfigOpt{
			lavalink.WithLogger(client.Logger()),
			lavalink.WithHTTPClient(client.Rest().HTTPClient()),
			lavalink.WithUserID(client.ApplicationID()),
		}, opts...)...),
	}

	client.EventManager().AddEventListeners(dgolink)
	return dgolink
}

type Link interface {
	lavalink.Lavalink
	bot.EventListener
}

var (
	_ Link              = (*linkImpl)(nil)
	_ bot.EventListener = (*linkImpl)(nil)
)

type linkImpl struct {
	lavalink.Lavalink
}

func (l *linkImpl) OnEvent(event bot.Event) {
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
