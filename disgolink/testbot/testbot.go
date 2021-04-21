package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
	"github.com/DisgoOrg/disgo/api/events"
	dapi "github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/disgolink"
	"github.com/sirupsen/logrus"
)

const guildID = "817327181659111454"

var logger = logrus.New()
var dgolink disgolink.Disgolink

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("starting testbot...")

	dgo, err := disgo.NewBuilder(endpoints.Token(os.Getenv("token"))).
		SetLogger(logger).
		SetIntents(api.IntentsGuilds | api.IntentsGuildMembers | api.IntentsGuildVoiceStates).
		SetCacheFlags(api.CacheFlagsDefault | api.CacheFlagVoiceState).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		AddEventListeners(&events.ListenerAdapter{
			OnSlashCommand: slashCommandListener,
		}).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	dgolink = disgolink.NewDisgolink(logger, dgo.ApplicationID())

	dgo.EventManager().AddEventListeners(dgolink)
	dgo.SetVoiceDispatchInterceptor(dgolink)

	port, _ := strconv.Atoi(os.Getenv("lavalink_port"))
	secure, _ := strconv.ParseBool(os.Getenv("lavalink_secure"))
	dgolink.AddNode(&dapi.NodeOptions{
		Name:     "test",
		Host:     os.Getenv("lavalink_host"),
		Port:     port,
		Password: os.Getenv("lavalink_password"),
		Secure:   secure,
	})

	_, err = dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), guildID, commands...)
	if err != nil {
		logger.Errorf("error while registering guild commands: %s", err)
	}

	err = dgo.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to discord: %s", err)
	}

	defer dgo.Close()

	logger.Infof("testbot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.CommandName {
	case "play":
		voiceState := event.Member.VoiceState()

		if voiceState == nil || voiceState.ChannelID == nil {
			_ = event.Reply(api.NewInteractionResponseBuilder().SetContent("Please join a VoiceChannel to use this command").Build())
			return
		}
		go func() {
			_ = event.Acknowledge()

			query := event.Option("query").String()
			searchProvider := event.Option("search-provider")
			if searchProvider != nil {
				switch searchProvider.String() {
				case "yt":
					query = string(dapi.SearchTypeYoutube) + query
				case "ytm":
					query = string(dapi.SearchTypeYoutubeMusic) + query
				case "sc":
					query = string(dapi.SearchTypeSoundCloud) + query
				}
			} else {
				if !dapi.URLPattern.MatchString(query) {
					query = string(dapi.SearchTypeYoutube) + query
				}
			}

			dgolink.RestClient().LoadItemAsync(query, dapi.NewResultHandler(
				func(track *dapi.Track) {
					queueOrPlay(event, voiceState, track)
				},
				func(playlist *dapi.Playlist) {
					track := playlist.SelectedTrack()
					if track == nil {
						track = playlist.Tracks[0]
					}
					queueOrPlay(event, voiceState, track)
				},
				func(tracks []*dapi.Track) {
					queueOrPlay(event, voiceState, tracks[0])
				},
				func() {
					_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("no tracks found").Build())
				},
				func(e *dapi.Exception) {
					_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("error while loading:\n" + e.Error()).Build())
				},
			))
		}()
	}
}

func connect(event *events.SlashCommandEvent, voiceState *api.VoiceState) {
	err := voiceState.VoiceChannel().Connect()
	if err != nil {
		_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
		return
	}
}

func queueOrPlay(, track *dapi.Track) {

	player := dgolink.Player(event.GuildID.String())
	player.Play(track)
	_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("playing [" + track.Info.Title + "](" + *track.Info.URI + ")").Build())
}
