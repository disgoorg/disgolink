package main

import (
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/DisgoOrg/disgolink/lavalink"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
)

func main() {
	bot := &Bot{
		Link: lavalink.New(
			lavalink.WithUserID(""),
		),
	}
	bot.registerNodes()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

type Bot struct {
	Link lavalink.Lavalink
}

func (b *Bot) messageCreateHandler() {
	command := "!play channelID url"
	channelID := ""
	guildID := ""
	args := strings.Split(command, " ")
	if len(args) < 3 {
		// TODO: send error message
		return
	}
	query := strings.Join(args[2:], " ")
	if !urlPattern.MatchString(query) {
		query = "ytsearch:" + query
	}
	b.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			b.Play(guildID, args[1], channelID, track)
		},
		func(playlist lavalink.AudioPlaylist) {
			b.Play(guildID, args[1], channelID, playlist.Tracks[0])
		},
		func(tracks []lavalink.AudioTrack) {
			b.Play(guildID, args[1], channelID, tracks[0])
		},
		func() {
			// TODO: send error message
		},
		func(ex lavalink.FriendlyException) {
			// TODO: send error message
		},
	))
}

func (b *Bot) Play(guildID string, voiceChannelID string, channelID string, track lavalink.AudioTrack) {
	// TODO: join voice channel

	if err := b.Link.Player(guildID).Play(track); err != nil {
		// TODO: send error message
		return
	}
	// TODO: send playing message
}

func (b *Bot) registerNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	b.Link.AddNode(lavalink.NodeConfig{
		Name:     "test",
		Host:     os.Getenv("LAVALINK_HOST"),
		Port:     os.Getenv("LAVALINK_PORT"),
		Password: os.Getenv("LAVALINK_PASSWORD"),
		Secure:   secure,
	})
}
