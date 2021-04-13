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
					query = dapi.YoutubeSearchPrefix + query
				case "ytm":
					query = dapi.YoutubeMusicSearchPrefix + query
				case "sc":
					query = dapi.SoundCloudSearchPrefix + query
				}
			}

			result, err := dgolink.RestClient().LoadItem(query)
			if err != nil || result.Exception != nil {
				var errStr string
				if err != nil {
					errStr = err.Error()
				} else {
					errStr = result.Exception.Error.Error()
				}
				_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("error while loading:\n" + errStr).Build())
				return
			}
			if result.Tracks == nil || len(result.Tracks) == 0 {
				_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("no tracks found").Build())
				return
			}
			var track *dapi.Track
			if result.PlaylistInfo.SelectedTrack != -1 {
				track = result.Tracks[result.PlaylistInfo.SelectedTrack]
			} else {
				track = result.Tracks[0]
			}
			err = voiceState.VoiceChannel().Connect()
			if err != nil {
				_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
				return
			}
			dgolink.Player(event.GuildID.String()).PlayTrack(track)
			_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().SetContent("playing [" + track.Info.Title + "](" + track.Info.URI + ")").Build())
		}()
	}
}
