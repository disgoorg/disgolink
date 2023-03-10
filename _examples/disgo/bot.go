package main

import (
	"context"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v2/disgolink"
)

func newBot() *Bot {
	return &Bot{
		Queues: &QueueManager{
			queues: make(map[snowflake.ID]*Queue),
		},
	}
}

type Bot struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Handlers map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error
	Queues   *QueueManager
}

func (b *Bot) onApplicationCommand(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	handler, ok := b.Handlers[data.CommandName()]
	if !ok {
		log.Info("unknown command: ", data.CommandName())
		return
	}
	if err := handler(event, data); err != nil {
		log.Error("error handling command: ", err)
	}
}

func (b *Bot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
}

func (b *Bot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
