package main

import (
	"fmt"
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink"
	"github.com/DisgoOrg/disgolink/filters"
	"math/rand"
	"time"
)

func checkMusicPlayer(event *core.SlashCommandEvent) *MusicPlayer {
	musicPlayer, ok := musicPlayers[*event.GuildID]
	if !ok {
		_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContent("No MusicPlayer found for this guild").Build())
		return nil
	}
	return musicPlayer
}

func onSlashCommand(event *core.SlashCommandEvent) {
	switch event.CommandName {
	case "shuffle":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.Create(core.NewMessageCreateBuilder().SetContent("Queue is empty").Build())
			return
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(musicPlayer.queue), func(i, j int) {
			musicPlayer.queue[i], musicPlayer.queue[j] = musicPlayer.queue[j], musicPlayer.queue[i]
		})
		_ = event.Create(core.NewMessageCreateBuilder().SetContent("Queue shuffled").Build())

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
		_ = flts.Commit()

	case "queue":
		musicPlayer := checkMusicPlayer(event)
		if musicPlayers == nil {
			return
		}

		if len(musicPlayer.queue) == 0 {
			_ = event.Create(core.NewMessageCreateBuilder().SetContent("No songs in queue").Build())
		}
		tracks := ""
		for i, track := range musicPlayer.queue {
			tracks += fmt.Sprintf("%d. [%s](%s)\n", i+1, track.Info().Title(), *track.Info().URI())
		}
		_ = event.Create(core.NewMessageCreateBuilder().SetEmbeds(core.NewEmbedBuilder().
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
		_ = musicPlayer.player.Pause(pause)
		message := "paused"
		if !pause {
			message = "resumed"
		}
		_ = event.Create(core.NewMessageCreateBuilder().SetContent(message + " music").Build())

	case "play":
		voiceState := event.Member.VoiceState()

		if voiceState == nil || voiceState.ChannelID == nil {
			_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Please join a VoiceChannel to use this command").Build())
			return
		}
		go func() {
			_ = event.DeferCreate(false)

			query := event.Options["query"].String()
			if searchProvider, ok := event.Options["search-provider"]; ok {
				switch searchProvider.String() {
				case "yt":
					query = disgolink.SearchTypeYoutube.Apply(query)
				case "ytm":
					query = disgolink.SearchTypeYoutubeMusic.Apply(query)
				case "sc":
					query = disgolink.SearchTypeSoundCloud.Apply(query)
				}
			} else {
				if !URLPattern.MatchString(query) {
					query = disgolink.SearchTypeYoutube.Apply(query)
				}
			}
			musicPlayer, ok := musicPlayers[*event.GuildID]
			if !ok {
				musicPlayer = NewMusicPlayer(*event.GuildID)
				musicPlayers[*event.GuildID] = musicPlayer
			}

			dgolink.RestClient().LoadItemHandler(query, disgolink.NewResultHandler(
				func(track disgolink.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, track)
				},
				func(playlist *disgolink.Playlist) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, playlist.Tracks...)
				},
				func(tracks []disgolink.Track) {
					if ok = connect(event, voiceState); !ok {
						return
					}
					musicPlayer.Queue(event, tracks[0])
				},
				func() {
					_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContent("no tracks found").Build())
				},
				func(e *disgolink.Exception) {
					_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContent("error while loading track:\n" + e.Error()).Build())
				},
			))
		}()
	}
}
