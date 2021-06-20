package main

import (
	"github.com/DisgoOrg/disgolink"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	dapi "github.com/DisgoOrg/disgolink/api"
	"github.com/sirupsen/logrus"
)

const guildID = "817327181659111454"

var logger = logrus.New()
var dgolink dapi.Disgolink
var musicPlayers = map[string]*MusicPlayer{}

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting example...")

	dgo, err := disgo.NewBuilder(os.Getenv("token")).
		SetLogger(logger).
		SetGatewayIntents(api.GatewayIntentsNonPrivileged).
		SetCacheFlags(api.CacheFlagsDefault | api.CacheFlagVoiceState).
		SetMemberCachePolicy(api.MemberCachePolicyNone).
		AddEventListeners(&events.ListenerAdapter{
			OnCommand: commandListener,
		}).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	dgolink = disgolink.NewDisgolink(logger, dgo)
	registerNodes()

	defer dgolink.Close()

	_, err = dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), guildID, commands...)
	if err != nil {
		logger.Errorf("error while registering guild commands: %s", err)
	}

	defer dgo.Close()

	err = dgo.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to discord: %s", err)
	}

	logger.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func connect(event events.CommandEvent, voiceState *api.VoiceState) bool {
	err := voiceState.VoiceChannel().Connect()
	if err != nil {
		_, _ = event.EditOriginal(api.NewMessageUpdateBuilder().SetContent("error while connecting to channel:\n" + err.Error()).Build())
		logger.Errorf("error while connecting to channel: %s", err)
		return false
	}
	return true
}

func registerNodes() {
	port, _ := strconv.Atoi(os.Getenv("lavalink_port"))
	secure, _ := strconv.ParseBool(os.Getenv("lavalink_secure"))
	dgolink.AddNode(&dapi.NodeOptions{
		Name:     "test",
		Host:     os.Getenv("lavalink_host"),
		Port:     port,
		Password: os.Getenv("lavalink_password"),
		Secure:   secure,
	})
}
