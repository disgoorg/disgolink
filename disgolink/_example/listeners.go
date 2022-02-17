package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/lavalink"
)

func checkMusicPlayer(event *events.ApplicationCommandInteractionEvent) *MusicPlayer {
	musicPlayer, ok := musicPlayers[*event.GuildID]
	if !ok {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent("No MusicPlayer found for this guild").Build())
		return nil
	}
	return musicPlayer
}

func onApplicationCommand(event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName {
	case "shuffle":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Queue is empty").Build())
			return
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(musicPlayer.queue), func(i, j int) {
			musicPlayer.queue[i], musicPlayer.queue[j] = musicPlayer.queue[j], musicPlayer.queue[i]
		})
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Queue shuffled").Build())

	case "filter":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		flts := musicPlayer.Filters()
		if flts.Timescale() == nil {
			flts.Timescale().Speed = 2
		} else {
			flts.SetTimescale(nil)
		}
		_ = flts.Commit()

	case "queue":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No songs in queue").Build())
		}
		tracks := ""
		for i, track := range musicPlayer.queue {
			tracks += fmt.Sprintf("%d. [%s](%s)\n", i+1, track.Info().Title, *track.Info().URI)
		}
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().
			SetTitle("Queue:").
			SetDescription(tracks).
			Build(),
		).Build())

	case "pause":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		pause := !musicPlayer.Paused()
		_ = musicPlayer.Pause(pause)
		message := "paused"
		if !pause {
			message = "resumed"
		}
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(message + " music").Build())

	case "play":
		voiceState := event.Member.VoiceState()

		if voiceState == nil || voiceState.ChannelID == nil {
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Please join a VoiceChannel to use this command").Build())
			return
		}
		go func() {
			_ = event.DeferCreateMessage(false)

			query := *data.Options.String("query")
			if searchProvider := data.Options.String("search-provider"); searchProvider != nil {
				switch *searchProvider {
				case "yt":
					query = lavalink.SearchTypeYoutube.Apply(query)
				case "ytm":
					query = lavalink.SearchTypeYoutubeMusic.Apply(query)
				case "sc":
					query = lavalink.SearchTypeSoundCloud.Apply(query)
				}
			} else {
				if !URLPattern.MatchString(query) {
					query = lavalink.SearchTypeYoutube.Apply(query)
				}
			}
			musicPlayer, ok := musicPlayers[*event.GuildID]
			if !ok {
				musicPlayer = NewMusicPlayer(*event.GuildID)
				musicPlayers[*event.GuildID] = musicPlayer
			}
			skipSegments := false
			if option := data.Options.Bool("skip-segments"); option != nil {
				skipSegments = *option
			}

			_ = musicPlayer.Node().RestClient().LoadItemHandler(context.TODO(), query, lavalink.NewResultHandler(
				func(track lavalink.AudioTrack) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, skipSegments, track)
				},
				func(playlist lavalink.AudioPlaylist) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, skipSegments, playlist.Tracks()...)
				},
				func(tracks []lavalink.AudioTrack) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, skipSegments, tracks[0])
				},
				func() {
					_, _ = event.UpdateOriginalMessage(discord.NewMessageUpdateBuilder().SetContent("no tracks found").Build())
				},
				func(e lavalink.FriendlyException) {
					_, _ = event.UpdateOriginalMessage(discord.NewMessageUpdateBuilder().SetContent("error while loading track:\n" + e.Error()).Build())
				},
			))
		}()
	}
}
