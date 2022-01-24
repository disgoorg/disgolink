package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/DisgoOrg/disgolink/arikawalink"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

	token = os.Getenv("DISCORD_TOKEN")
)

func main() {
	log.Println("Starting...")
	s := session.New("Bot " + token)
	s.AddIntents(gateway.IntentGuildVoiceStates)
	s.AddIntents(gateway.IntentGuildMessages)
	bot := &Bot{
		Link:    arikawalink.New(s),
		Session: s,
	}
	s.AddHandler(bot.readyHandler)
	s.AddHandler(bot.messageCreateHandler)

	if err := s.Open(context.TODO()); err != nil {
		panic(err)
	}
	defer s.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

type Bot struct {
	Session *session.Session
	Link    *arikawalink.Link
}

func (b *Bot) readyHandler(_ *gateway.ReadyEvent) {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	b.Link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:     "test",
		Host:     os.Getenv("LAVALINK_HOST"),
		Port:     os.Getenv("LAVALINK_PORT"),
		Password: os.Getenv("LAVALINK_PASSWORD"),
		Secure:   secure,
	})
}

func (b *Bot) messageCreateHandler(e *gateway.MessageCreateEvent) {
	if e.Author.Bot {
		return
	}
	args := strings.Split(e.Content, " ")
	switch args[0] {
	case "!play":
		if len(args) < 3 {
			_, _ = b.Session.Client.SendMessage(e.ChannelID, "Please provide a channel id and something to play")
			return
		}
		query := strings.Join(args[2:], " ")
		if !urlPattern.MatchString(query) {
			query = "ytsearch:" + query
		}
		channelID, _ := discord.ParseSnowflake(args[1])
		_ = b.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
			func(track lavalink.AudioTrack) {
				b.play(e.GuildID, discord.ChannelID(channelID), e.ChannelID, track)
			},
			func(playlist lavalink.AudioPlaylist) {
				b.play(e.GuildID, discord.ChannelID(channelID), e.ChannelID, playlist.Tracks[0])
			},
			func(tracks []lavalink.AudioTrack) {
				b.play(e.GuildID, discord.ChannelID(channelID), e.ChannelID, tracks[0])
			},
			func() {
				_, _ = b.Session.Client.SendMessage(e.ChannelID, "no matches found for: "+query)
			},
			func(ex lavalink.FriendlyException) {
				_, _ = b.Session.Client.SendMessage(e.ChannelID, "error while loading track: "+ex.Message)
			},
		))

	}
}

func (b *Bot) play(guildID discord.GuildID, voiceChannelID discord.ChannelID, channelID discord.ChannelID, track lavalink.AudioTrack) {
	if err := b.Session.Gateway().Send(context.TODO(), &gateway.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: voiceChannelID,
		SelfMute:  false,
		SelfDeaf:  false,
	}); err != nil {
		_, _ = b.Session.Client.SendMessage(channelID, "error while joining voice channel: "+err.Error())
		return
	}
	if err := b.Link.Player(snowflake.ParseUInt64(uint64(guildID))).Play(track); err != nil {
		_, _ = b.Session.Client.SendMessage(channelID, "error while playing track: "+err.Error())
		return
	}
	_, _ = b.Session.Client.SendMessage(channelID, "Playing: "+track.Info().Title())
}
