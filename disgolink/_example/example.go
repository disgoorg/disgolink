package main

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	URLPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

	token        = os.Getenv("disgolink_token")
	guildID      = snowflake.GetEnv("guild_id")
	client       bot.Client
	dgolink      disgolink.Link
	musicPlayers = map[snowflake.ID]*MusicPlayer{}
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting _example...")

	var err error
	client, err = disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds|discord.GatewayIntentGuildVoiceStates)),
		bot.WithCacheConfigOpts(cache.WithCacheFlags(cache.FlagVoiceStates)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: onApplicationCommand,
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgolink instance: %s", err)
		return
	}

	defer client.Close(context.TODO())

	dgolink = disgolink.New(client)
	registerNodes()

	defer dgolink.Close()

	_, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	err = client.ConnectGateway(context.TODO())
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func connect(event *events.ApplicationCommandInteractionEvent, voiceState discord.VoiceState) bool {
	if err := event.Client().Connect(context.TODO(), voiceState.GuildID, *voiceState.ChannelID); err != nil {
		_, _ = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetContent("error while connecting to channel:\n"+err.Error()).Build())
		log.Errorf("error while connecting to channel: %s", err)
		return false
	}
	return true
}

func registerNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("lavalink_secure"))
	_, _ = dgolink.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "test",
		Host:        os.Getenv("lavalink_host"),
		Port:        os.Getenv("lavalink_port"),
		Password:    os.Getenv("lavalink_password"),
		Secure:      secure,
		ResumingKey: os.Getenv("lavalink_resuming_key"),
	})
	if os.Getenv("lavalink_resuming_key") != "" {
		_ = dgolink.BestNode().ConfigureResuming(os.Getenv("lavalink_resuming_key"), 20)
	}
}
