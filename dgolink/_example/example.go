package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/dgolink"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake"
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
		Link:           dgolink.New(session),
		PlayerManagers: map[string]*PlayerManager{},
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
	Link           *dgolink.Link
	PlayerManagers map[string]*PlayerManager
}

type PlayerManager struct {
	lavalink.PlayerEventAdapter
	Player        lavalink.Player
	Queue         []lavalink.AudioTrack
	QueueMu       sync.Mutex
	RepeatingMode RepeatingMode
}

func (m *PlayerManager) AddQueue(tracks ...lavalink.AudioTrack) {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	m.Queue = append(m.Queue, tracks...)
}

func (m *PlayerManager) PopQueue() lavalink.AudioTrack {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	if len(m.Queue) == 0 {
		return nil
	}
	var track lavalink.AudioTrack
	track, m.Queue = m.Queue[0], m.Queue[1:]
	return track
}

func (m *PlayerManager) OnTrackEnd(player lavalink.Player, track lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	if !endReason.MayStartNext() {
		return
	}
	switch m.RepeatingMode {
	case RepeatingModeOff:
		if nextTrack := m.PopQueue(); nextTrack != nil {
			if err := player.Play(nextTrack); err != nil {
				fmt.Println("error playing next track:", err)
			}
		}
	case RepeatingModeSong:
		if err := player.Play(track.Clone()); err != nil {
			fmt.Println("error playing next track:", err)
		}

	case RepeatingModeQueue:
		m.AddQueue(track)
		if nextTrack := m.PopQueue(); nextTrack != nil {
			if err := player.Play(nextTrack); err != nil {
				fmt.Println("error playing next track:", err)
			}
		}
	}
}

type RepeatingMode int

const (
	RepeatingModeOff = iota
	RepeatingModeSong
	RepeatingModeQueue
)

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
		_ = b.Link.BestRestClient().LoadItemHandler(context.TODO(), query, lavalink.NewResultHandler(
			func(track lavalink.AudioTrack) {
				b.play(s, e.GuildID, args[1], e.ChannelID, track)
			},
			func(playlist lavalink.AudioPlaylist) {
				b.play(s, e.GuildID, args[1], e.ChannelID, playlist.Tracks()...)
			},
			func(tracks []lavalink.AudioTrack) {
				b.play(s, e.GuildID, args[1], e.ChannelID, tracks[0])
			},
			func() {
				_, _ = s.ChannelMessageSend(e.ChannelID, "no matches found for: "+query)
			},
			func(ex lavalink.FriendlyException) {
				_, _ = s.ChannelMessageSend(e.ChannelID, "error while loading track: "+ex.Message)
			},
		))

	}
}

func (b *Bot) play(s *discordgo.Session, guildID string, voiceChannelID string, channelID string, tracks ...lavalink.AudioTrack) {
	if err := s.ChannelVoiceJoinManual(guildID, voiceChannelID, false, false); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while joining voice channel: "+err.Error())
		return
	}

	manager, ok := b.PlayerManagers[guildID]
	if !ok {
		manager = &PlayerManager{
			Player:        b.Link.Player(snowflake.Snowflake(guildID)),
			RepeatingMode: RepeatingModeOff,
		}
		b.PlayerManagers[guildID] = manager
		manager.Player.AddListener(manager)
	}
	manager.AddQueue(tracks...)

	track := manager.PopQueue()
	if err := manager.Player.Play(track); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while playing track: "+err.Error())
		return
	}
	_, _ = s.ChannelMessageSend(channelID, "Playing: "+track.Info().Title)
}

func (b *Bot) registerNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	b.Link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:     "test",
		Host:     os.Getenv("LAVALINK_HOST"),
		Port:     os.Getenv("LAVALINK_PORT"),
		Password: os.Getenv("LAVALINK_PASSWORD"),
		Secure:   secure,
	})
}
