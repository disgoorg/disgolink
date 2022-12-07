package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/log"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)

	TOKEN    = os.Getenv("TOKEN")
	GUILD_ID = os.Getenv("GUILD_ID")

	NODE_NAME      = os.Getenv("NODE_NAME")
	NODE_ADDRESS   = os.Getenv("NODE_ADDRESS")
	NODE_PASSWORD  = os.Getenv("NODE_PASSWORD")
	NODE_SECURE, _ = strconv.ParseBool(os.Getenv("NODE_SECURE"))
)

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Handlers map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues   *QueueManager
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting discordgo example...")
	log.Info("discordgo version: ", discordgo.VERSION)
	log.Info("disgolink version: ", disgolink.Version)

	b := &Bot{
		Queues: &QueueManager{
			queues: make(map[string]*Queue),
		},
	}

	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	b.Session = session

	session.State.TrackVoice = true
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(b.onApplicationCommand)
	session.AddHandler(b.onVoiceStateUpdate)
	session.AddHandler(b.onVoiceServerUpdate)

	if err = session.Open(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	registerCommands(session)

	b.Lavalink = disgolink.New(snowflake.MustParse(session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"play":        b.play,
		"pause":       b.pause,
		"now-playing": b.nowPlaying,
		"stop":        b.stop,
		"queue":       b.queue,
		"clear-queue": b.clearQueue,
		"queue-type":  b.queueType,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NODE_NAME,
		Address:  NODE_ADDRESS,
		Password: NODE_PASSWORD,
		Secure:   NODE_SECURE,
	})
	if err != nil {
		log.Fatal(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("node version: %s", version)

	log.Info("DiscordGo example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	data := event.ApplicationCommandData()

	handler, ok := b.Handlers[data.Name]
	if !ok {
		log.Info("unknown command: ", data.Name)
		return
	}
	if err := handler(event, data); err != nil {
		log.Error("error handling command: ", err)
	}
}

func (b *Bot) onVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.UserID != session.State.User.ID {
		return
	}

	var guildID *snowflake.ID
	if event.GuildID != "" {
		id := snowflake.MustParse(event.GuildID)
		guildID = &id
	}
	b.Lavalink.OnVoiceStateUpdate(snowflake.MustParse(event.GuildID), guildID, event.SessionID)
	if event.ChannelID == "" {
		b.Queues.Delete(event.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}
