package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgolink"
	"github.com/DisgoOrg/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgolink/api"
)

var (
	token        = os.Getenv("disgolink_token")
	guildID      = discord.Snowflake(os.Getenv("guild_id"))
	disgo        core.Disgo
	dgolink      api.Disgolink
	musicPlayers = map[discord.Snowflake]*MusicPlayer{}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting _example...")

	var err error
	disgo, err = core.NewBuilder(token).
		SetGatewayConfig(gateway.Config{
			GatewayIntents: discord.GatewayIntentsNonPrivileged,
		}).
		SetCacheConfig(core.CacheConfig{
			CacheFlags:        core.CacheFlagsDefault | core.CacheFlagVoiceState,
			MemberCachePolicy: core.MemberCachePolicyNone,
		}).
		AddEventListeners(&events.ListenerAdapter{
			OnSlashCommand: onSlashCommand,
		}).
		Build()
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
		return
	}

	defer disgo.Close()

	dgolink = disgolink.NewDisgolink(disgo)
	registerNodes()

	defer dgolink.Close()

	_, err = disgo.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	err = disgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	log.Infof("_example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func connect(event *events.SlashCommandEvent, voiceState *core.VoiceState) bool {
	err := voiceState.VoiceChannel().Connect()
	if err != nil {
		_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
		log.Errorf("error while connecting to channel: %s", err)
		return false
	}
	return true
}

func registerNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("lavalink_secure"))
	dgolink.AddNode(&api.NodeOptions{
		Name:     "test",
		Host:     os.Getenv("lavalink_host"),
		Port:     os.Getenv("lavalink_port"),
		Password: os.Getenv("lavalink_password"),
		Secure:   secure,
	})
}
