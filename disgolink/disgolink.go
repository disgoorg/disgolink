package disgolink

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/lavalink"
)

func New(client bot.Client, opts ...lavalink.ConfigOpt) Link {
	disgolink := &linkImpl{
		Lavalink: lavalink.New(append([]lavalink.ConfigOpt{
			lavalink.WithLogger(client.Logger()),
			lavalink.WithHTTPClient(client.Rest().HTTPClient()),
			lavalink.WithUserID(client.ApplicationID()),
		}, opts...)...),
	}

	client.EventManager().AddEventListeners(disgolink)
	return disgolink
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
	case *events.VoiceServerUpdate:
		l.OnVoiceServerUpdate(lavalink.VoiceServerUpdate{
			Token:    e.Token,
			GuildID:  e.GuildID,
			Endpoint: e.Endpoint,
		})

	case *events.GuildVoiceStateUpdate:
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
