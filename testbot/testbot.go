package testbot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
	"github.com/DisgoOrg/disgo/api/events"
	"github.com/DisgoOrg/disgolink"
	"github.com/sirupsen/logrus"
)

const guildID = "817327181659111454"

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting testbot...")

	lavalink := disgolink.NewDisgolink()

	dgo, err := disgo.NewBuilder(endpoints.Token(os.Getenv("token"))).
		SetLogger(logger).
		SetIntents(api.IntentsGuilds | api.IntentsGuildMessages | api.IntentsGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyVoice).
		AddEventListeners(&events.ListenerAdapter{
			OnSlashCommand: slashCommandListener,
		}).
		SetVoiceDispatchInterceptor(lavalink).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	_, err = dgo.RestClient().SetGuildCommands(dgo.SelfUserID(), guildID, commands...)
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

}
