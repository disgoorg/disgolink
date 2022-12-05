package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/json"
	"time"
)

func (b *Bot) players(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	var description string
	b.Lavalink.ForPlayers(func(player disgolink.Player) {
		description += fmt.Sprintf("GuildID: `%s`\n", player.GuildID())
	})

	return event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Players:\n%s", description),
	})
}

func (b *Bot) pause(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*event.GuildID())
	if player == nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: "No player found",
		})
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(!player.Paused())); err != nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Error while pausing: `%s`", err),
		})
	}

	status := "playing"
	if player.Paused() {
		status = "paused"
	}
	return event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Player is now %s", status),
	})
}

func (b *Bot) stop(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*event.GuildID())
	if player == nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: "No player found",
		})
	}

	if err := b.Client.Disconnect(context.TODO(), *event.GuildID()); err != nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Error while disconnecting: `%s`", err),
		})
	}

	return event.CreateMessage(discord.MessageCreate{
		Content: "Player stopped",
	})
}

func (b *Bot) nowPlaying(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(*event.GuildID())
	if player == nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: "No player found",
		})
	}

	track := player.Track()
	if track == nil {
		return event.CreateMessage(discord.MessageCreate{
			Content: "No track found",
		})
	}

	return event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Now playing: [`%s`](<%s>)\n\n %s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length)),
	})
}

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

func (b *Bot) play(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
	identifier := data.String("identifier")
	if source, ok := data.OptString("source"); ok {
		identifier = lavalink.SearchType(source).Apply(identifier)
	} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYoutube.Apply(identifier)
	}

	voiceState, ok := b.Client.Caches().VoiceStates().Get(*event.GuildID(), event.User().ID)
	if !ok {
		return event.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command",
		})
	}

	if err := event.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracks(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			_, _ = b.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded track: [`%s`](<%s>)", track.Info.Title, *track.Info.URI)),
			})
			toPlay = &track
		},
		func(playlist lavalink.Playlist) {
			_, _ = b.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
			})
			toPlay = &playlist.Tracks[0]
		},
		func(tracks []lavalink.Track) {
			_, _ = b.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded search result: [`%s`](<%s>)", tracks[0].Info.Title, *tracks[0].Info.URI)),
			})
			toPlay = &tracks[0]
		},
		func() {
			_, _ = b.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			_, _ = b.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Client.Connect(context.TODO(), *event.GuildID(), *voiceState.ChannelID); err != nil {
		return err
	}

	return b.Lavalink.Player(*event.GuildID()).Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
