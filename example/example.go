package main

import (
	"github.com/DisgoOrg/disgolink"
	"github.com/DisgoOrg/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgo"
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	"github.com/DisgoOrg/disgolink/api"
)

var (
	token        = os.Getenv("disgolink_token")
	guildID      = dapi.Snowflake(os.Getenv("guild_id"))
	dgo          dapi.Disgo
	dgolink      api.Disgolink
	musicPlayers = map[dapi.Snowflake]*MusicPlayer{}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")

	var err error
	dgo, err = disgo.NewBuilder(token).
		SetGatewayIntents(dapi.GatewayIntentsNonPrivileged).
		SetCacheFlags(dapi.CacheFlagsDefault | dapi.CacheFlagVoiceState).
		SetMemberCachePolicy(dapi.MemberCachePolicyNone).
		AddEventListeners(&events.ListenerAdapter{
			OnCommand: commandListener,
		}).
		Build()
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
		return
	}

	dgolink = disgolink.NewDisgolink(dgo)
	registerNodes()

	defer dgolink.Close()

	_, err = dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), guildID, commands...)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	defer dgo.Close()

	err = dgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func connect(event *events.CommandEvent, voiceState *dapi.VoiceState) bool {
	err := voiceState.VoiceChannel().Connect()
	if err != nil {
		_, _ = event.EditOriginal(dapi.NewMessageUpdateBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
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
