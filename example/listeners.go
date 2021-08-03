package main

import (
	"fmt"
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/filters"
	"math/rand"
	"time"
)

func checkMusicPlayer(event *events.CommandEvent) *MusicPlayer {
	musicPlayer, ok := musicPlayers[*event.Interaction.GuildID]
	if !ok {
		_ = event.Reply(dapi.NewMessageCreateBuilder().SetEphemeral(true).SetContent("No MusicPlayer found for this guild").Build())
		return nil
	}
	return musicPlayer
}

func commandListener(event *events.CommandEvent) {
	switch event.CommandName() {
	case "shuffle":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.Reply(dapi.NewMessageCreateBuilder().SetContent("Queue is empty").Build())
			return
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(musicPlayer.queue), func(i, j int) {
			musicPlayer.queue[i], musicPlayer.queue[j] = musicPlayer.queue[j], musicPlayer.queue[i]
		})
		_ = event.Reply(dapi.NewMessageCreateBuilder().SetContent("Queue shuffled").Build())

	case "filter":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		flts := musicPlayer.player.Filters()
		if flts.Timescale == nil {
			flts.Timescale = &filters.Timescale{Speed: 2}
		} else {
			flts.Timescale = nil
		}
		flts.Commit()

	case "queue":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.Reply(dapi.NewMessageCreateBuilder().SetContent("No songs in queue").Build())
		}
		tracks := ""
		for i, track := range musicPlayer.queue {
			tracks += fmt.Sprintf("%d. [%s](%s)\n", i+1, track.Info().Title(), *track.Info().URI())
		}
		_ = event.Reply(dapi.NewMessageCreateBuilder().SetEmbeds(dapi.NewEmbedBuilder().
			SetTitle("Queue:").
			SetDescription(tracks).
			Build(),
		).Build())

	case "pause":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		pause := !musicPlayer.player.Paused()
		musicPlayer.player.Pause(pause)
		message := "paused"
		if !pause {
			message = "resumed"
		}
		_ = event.Reply(dapi.NewMessageCreateBuilder().SetContent(message + " music").Build())

	case "play":
		voiceState := event.Interaction.Member.VoiceState()

		if voiceState == nil || voiceState.ChannelID == nil {
			_ = event.Reply(dapi.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Please join a VoiceChannel to use this command").Build())
			return
		}
		go func() {
			_ = event.DeferReply(false)

			query := event.Option("query").String()
			searchProvider := event.Option("search-provider")
			if searchProvider != nil {
				switch searchProvider.String() {
				case "yt":
					query = api.SearchTypeYoutube.Apply(query)
				case "ytm":
					query = api.SearchTypeYoutubeMusic.Apply(query)
				case "sc":
					query = api.SearchTypeSoundCloud.Apply(query)
				}
			} else {
				if !api.URLPattern.MatchString(query) {
					query = string(api.SearchTypeYoutube) + query
				}
			}
			musicPlayer, ok := musicPlayers[*event.Interaction.GuildID]
			if !ok {
				musicPlayer = NewMusicPlayer(*event.Interaction.GuildID)
				musicPlayers[*event.Interaction.GuildID] = musicPlayer
			}
			dgolink.RestClient().LoadItemHandler(query, api.NewResultHandler(
				func(track api.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, track)
				},
				func(playlist *api.Playlist) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, playlist.Tracks...)
				},
				func(tracks []api.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, tracks[0])
				},
				func() {
					_, _ = event.EditOriginal(dapi.NewMessageUpdateBuilder().SetContent("no tracks found").Build())
				},
				func(e *api.Exception) {
					_, _ = event.EditOriginal(dapi.NewMessageUpdateBuilder().SetContent("error while loading:\n" + e.Error()).Build())
				},
			))
		}()
	}
}
