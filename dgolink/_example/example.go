package main

import (
	"github.com/DisgoOrg/disgolink/dgolink"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

	token = os.Getenv("DISCORD_TOKEN")
)

func main() {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	bot := &Bot{
		Link: dgolink.New(session),
	}
	session.AddHandler(bot.messageCreateHandler)

	if err = session.Open(); err != nil {
		panic(err)
	}
	defer session.Close()
	bot.registerNodes()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

type Bot struct {
	Link *dgolink.Link
}

func (b *Bot) messageCreateHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}
	args := strings.Split(e.Content, " ")
	switch args[0] {
	case "!play":
		if len(args) < 3 {
			_, _ = s.ChannelMessageSend(e.ChannelID, "Please provide a channel id and something to play")
			return
		}
		query := strings.Join(args[2:], " ")
		if !urlPattern.MatchString(query) {
			query = "ytsearch:" + query
		}
		b.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
			func(track lavalink.Track) {
				play(s, b.Link, e.GuildID, args[1], e.ChannelID, track)
			},
			func(playlist lavalink.Playlist) {
				play(s, b.Link, e.GuildID, args[1], e.ChannelID, playlist.Tracks[0])
			},
			func(tracks []lavalink.Track) {
				play(s, b.Link, e.GuildID, args[1], e.ChannelID, tracks[0])
			},
			func() {
				_, _ = s.ChannelMessageSend(e.ChannelID, "no matches found for: "+query)
			},
			func(ex lavalink.Exception) {
				_, _ = s.ChannelMessageSend(e.ChannelID, "error while loading track: "+ex.Message)
			},
		))

	}
}

func play(s *discordgo.Session, link *dgolink.Link, guildID string, voiceChannelID string, channelID string, track lavalink.Track) {
	if err := s.ChannelVoiceJoinManual(guildID, voiceChannelID, false, false); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while joining voice channel: "+err.Error())
		return
	}
	if err := link.Player(guildID).Play(track); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while playing track: "+err.Error())
		return
	}
	_, _ = s.ChannelMessageSend(channelID, "Playing: "+track.Info().Title())
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
