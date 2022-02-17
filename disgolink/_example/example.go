package main

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgolink/disgolink"
	"github.com/DisgoOrg/snowflake"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/log"
)

var (
	URLPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

	token        = os.Getenv("disgolink_token")
	guildID      = snowflake.GetSnowflakeEnv("guild_id")
	disgo        *core.Bot
	dgolink      disgolink.Link
	musicPlayers = map[snowflake.Snowflake]*MusicPlayer{}
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting _example...")

	var err error
	disgo, err = bot.New(token,
		bot.WithGatewayOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds|discord.GatewayIntentGuildVoiceStates)),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsDefault)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: onApplicationCommand,
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgolink instance: %s", err)
		return
	}

	defer disgo.Close(context.TODO())

	dgolink = disgolink.New(disgo)
	registerNodes()

	defer dgolink.Close()

	_, err = disgo.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	err = disgo.ConnectGateway(context.TODO())
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func connect(event *events.ApplicationCommandInteractionEvent, voiceState *core.VoiceState) bool {
	channel := voiceState.Channel()
	err := channel.Connect(context.TODO())
	if err != nil {
		_, _ = event.UpdateOriginalMessage(discord.NewMessageUpdateBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
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
	_ = dgolink.BestNode().ConfigureResuming("test", 20)
}
