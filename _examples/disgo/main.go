package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/log"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)

	token   = os.Getenv("TOKEN")
	guildID = snowflake.GetEnv("GUILD_ID")

	nodeName      = os.Getenv("NODE_NAME")
	nodeAddress   = os.Getenv("NODE_ADDRESS")
	nodePassword  = os.Getenv("NODE_PASSWORD")
	nodeSecure, _ = strconv.ParseBool(os.Getenv("NODE_SECURE"))
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting disgo example...")
	log.Info("disgo version: ", disgo.Version)
	log.Info("disgolink version: ", disgolink.Version)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagVoiceStates),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// registerCommands(client)

	lavalinkClient, err := disgolink.New(client.ApplicationID())
	if err != nil {
		log.Fatal(err)
	}

	client.AddEventListeners(
		bot.NewListenerFunc(onApplicationCommand(client, lavalinkClient)),
		bot.NewListenerFunc(onVoiceStateUpdate(lavalinkClient)),
		bot.NewListenerFunc(onVoiceServerUpdate(lavalinkClient)),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.TODO())

	node, err := lavalinkClient.AddNode(ctx, disgolink.NodeConfig{
		Name:     nodeName,
		Address:  nodeAddress,
		Password: nodePassword,
		Secure:   nodeSecure,
	})
	if err != nil {
		log.Fatal(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("node version: %s", version)

	log.Info("DisGo example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onApplicationCommand(client bot.Client, lavalinkClient disgolink.Client) func(event *events.ApplicationCommandInteractionCreate) {
	return func(event *events.ApplicationCommandInteractionCreate) {
		data := event.SlashCommandInteractionData()
		switch data.CommandName() {
		case "play":
			identifier := data.String("identifier")
			if source, ok := data.OptString("source"); ok {
				identifier = lavalink.SearchType(source).Apply(identifier)
			} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
				identifier = lavalink.SearchTypeYoutube.Apply(identifier)
			}

			voiceState, ok := client.Caches().VoiceStates().Get(*event.GuildID(), event.User().ID)
			if !ok {
				_ = event.CreateMessage(discord.MessageCreate{
					Content: "You need to be in a voice channel to use this command",
				})
				return
			}

			if err := event.DeferCreateMessage(false); err != nil {
				log.Error(err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var toPlay *lavalink.Track
			lavalinkClient.BestNode().LoadTracks(ctx, identifier, disgolink.NewResultHandler(
				func(track lavalink.Track) {
					_, _ = client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
						Content: json.Ptr(fmt.Sprintf("Loaded track: (`%s`)[%s]", track.Info.Title, *track.Info.URI)),
					})
					toPlay = &track
				},
				func(playlist lavalink.Playlist) {
					_, _ = client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
						Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
					})
					toPlay = &playlist.Tracks[0]
				},
				func(tracks []lavalink.Track) {
					_, _ = client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
						Content: json.Ptr(fmt.Sprintf("Loaded search result: (`%s`)[%s]", tracks[0].Info.Title, *tracks[0].Info.URI)),
					})
					toPlay = &tracks[0]
				},
				func() {
					_, _ = client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
						Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
					})
				},
				func(err error) {
					_, _ = client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
						Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
					})
				},
			))
			if toPlay == nil {
				return
			}

			if err := client.Connect(context.TODO(), *event.GuildID(), *voiceState.ChannelID); err != nil {
				log.Error(err)
				return
			}

			player := lavalinkClient.Player(*event.GuildID())
			if err := player.Update(context.TODO(), lavalink.PlayerUpdate{
				EncodedTrack: json.NewNullablePtr(toPlay.Encoded),
			}); err != nil {
				log.Error(err)
			}
		}
	}
}

func onVoiceStateUpdate(lavalinkClient disgolink.Client) func(event *events.GuildVoiceStateUpdate) {
	return func(event *events.GuildVoiceStateUpdate) {
		lavalinkClient.OnVoiceStateUpdate(event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	}
}

func onVoiceServerUpdate(lavalinkClient disgolink.Client) func(event *events.VoiceServerUpdate) {
	return func(event *events.VoiceServerUpdate) {
		lavalinkClient.OnVoiceServerUpdate(event.GuildID, event.Token, *event.Endpoint)
	}
}
