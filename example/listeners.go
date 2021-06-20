package main

import (
	"fmt"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	dapi "github.com/DisgoOrg/disgolink/api"
)

func commandListener(event events.CommandEvent) {
	switch event.CommandName {
	case "queue":
		musicPlayer, ok := musicPlayers[event.Interaction.GuildID.String()]
		if !ok {
			_ = event.Reply(api.NewMessageCreateBuilder().SetContent("No MusicPlayer found for this guild").Build())
			return
		}
		tracks := ""
		for i, track := range musicPlayer.queue {
			tracks += fmt.Sprintf("%d. [%s](%s)\n", i+1, track.Info().Title(), *track.Info().URI())
		}
		_ = event.Reply(api.NewMessageCreateBuilder().SetEmbeds(api.NewEmbedBuilder().
			SetTitle("Queue:").
			SetDescription(tracks).
			Build(),
		).Build())
	case "pause":
		musicPlayer, ok := musicPlayers[event.Interaction.GuildID.String()]
		if !ok {
			_ = event.Reply(api.NewMessageCreateBuilder().SetContent("No MusicPlayer found for this guild").Build())
			return
		}
		pause := !musicPlayer.player.Paused()
		musicPlayer.player.Pause(pause)
		message := "paused"
		if !pause {
			message = "resumed"
		}
		_ = event.Reply(api.NewMessageCreateBuilder().SetContent(message + "music").Build())
	case "play":
		voiceState := event.Interaction.Member.VoiceState()

		if voiceState == nil || voiceState.ChannelID == nil {
			_ = event.Reply(api.NewMessageCreateBuilder().SetContent("Please join a VoiceChannel to use this command").Build())
			return
		}
		go func() {
			_ = event.DeferReply(false)

			query := event.Option("query").String()
			searchProvider := event.Option("search-provider")
			if searchProvider != nil {
				switch searchProvider.String() {
				case "yt":
					query = dapi.SearchTypeYoutube.Apply(query)
				case "ytm":
					query = dapi.SearchTypeYoutubeMusic.Apply(query)
				case "sc":
					query = dapi.SearchTypeSoundCloud.Apply(query)
				}
			} else {
				if !dapi.URLPattern.MatchString(query) {
					query = string(dapi.SearchTypeYoutube) + query
				}
			}
			musicPlayer, ok := musicPlayers[event.Interaction.GuildID.String()]
			if !ok {
				musicPlayer = NewMusicPlayer(event.Interaction.GuildID.String())
				musicPlayers[event.Interaction.GuildID.String()] = musicPlayer
			}
			dgolink.RestClient().LoadItemAsync(query, dapi.NewResultHandler(
				func(track dapi.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, track)
				},
				func(playlist *dapi.Playlist) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, playlist.Tracks...)
				},
				func(tracks []dapi.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, tracks[0])
				},
				func() {
					_, _ = event.EditOriginal(api.NewMessageUpdateBuilder().SetContent("no tracks found").Build())
				},
				func(e *dapi.Exception) {
					_, _ = event.EditOriginal(api.NewMessageUpdateBuilder().SetContent("error while loading:\n" + e.Error()).Build())
				},
			))
		}()
	}
}
