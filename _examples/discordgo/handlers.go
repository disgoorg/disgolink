package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"time"
)

func (b *Bot) pause(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found",
			},
		})
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(!player.Paused())); err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error while pausing: `%s`", err),
			},
		})
	}

	status := "playing"
	if player.Paused() {
		status = "paused"
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Player is now %s", status),
		},
	})
}

func (b *Bot) stop(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found",
			},
		})
	}

	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, "", false, false); err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error while disconnecting: `%s`", err),
			},
		})
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Player stopped",
		},
	})
}

func (b *Bot) nowPlaying(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found",
			},
		})
	}

	track := player.Track()
	if track == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No track found",
			},
		})
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Now playing: [`%s`](<%s>)\n\n %s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length)),
		},
	})
}

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

func (b *Bot) play(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	identifier := data.Options[0].StringValue()
	if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYoutube.Apply(identifier)
	}

	if err := b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracks(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loading track: [`%s`](<%s>)", track.Info.Title, *track.Info.URI)),
			})
			toPlay = &track
		},
		func(playlist lavalink.Playlist) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
			})
			toPlay = &playlist.Tracks[0]
		},
		func(tracks []lavalink.Track) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loaded search result: [`%s`](<%s>)", tracks[0].Info.Title, *tracks[0].Info.URI)),
			})
			toPlay = &tracks[0]
		},
		func() {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, "", false, false); err != nil {
		return err
	}

	return b.Lavalink.Player(snowflake.MustParse(event.GuildID)).Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
