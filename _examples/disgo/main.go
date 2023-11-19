package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)

	Token   = os.Getenv("TOKEN")
	GuildId = snowflake.GetEnv("GUILD_ID")

	NodeName      = os.Getenv("NODE_NAME")
	NodeAddress   = os.Getenv("NODE_ADDRESS")
	NodePassword  = os.Getenv("NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(os.Getenv("NODE_SECURE"))
)

func main() {
	slog.Info("starting disgo example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))
	slog.Info("disgolink version: ", slog.String("version", disgolink.Version))

	b := newBot()

	client, err := disgo.New(Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.onApplicationCommand),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		slog.Error("error while building disgo client", slog.Any("err", err))
		os.Exit(1)
	}
	b.Client = client

	registerCommands(client)

	b.Lavalink = disgolink.New(client.ApplicationID(),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
		disgolink.WithListenerFunc(b.onUnknownEvent),
	)
	b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error{
		"play":        b.play,
		"pause":       b.pause,
		"now-playing": b.nowPlaying,
		"stop":        b.stop,
		"players":     b.players,
		"queue":       b.queue,
		"clear-queue": b.clearQueue,
		"queue-type":  b.queueType,
		"shuffle":     b.shuffle,
		"seek":        b.seek,
		"volume":      b.volume,
		"skip":        b.skip,
		"bass-boost":  b.bassBoost,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		slog.Error("failed to open gateway", slog.Any("err", err))
		os.Exit(1)
	}
	defer client.Close(context.TODO())

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})
	if err != nil {
		slog.Error("failed to add node", slog.Any("err", err))
		os.Exit(1)
	}
	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("failed to get node version", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("DisGo example is now running. Press CTRL-C to exit.", slog.String("node_version", version), slog.String("node_session_id", node.SessionID()))
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
