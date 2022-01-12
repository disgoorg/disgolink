package dgolink

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
)

func New(s *discordgo.Session, opts ...lavalink.ConfigOpt) *Link {
	opts = append(opts, lavalink.WithHTTPClient(s.Client))
	discordgolink := &Link{
		Lavalink: lavalink.NewLavalink(opts...),
	}

	s.AddHandler(discordgolink.ReadyHandler)
	s.AddHandler(discordgolink.VoiceStateUpdateHandler)
	s.AddHandler(discordgolink.VoiceServerUpdateHandler)
	return discordgolink
}

var _ lavalink.Lavalink = (*Link)(nil)

type Link struct {
	lavalink.Lavalink
}

func (l *Link) ReadyHandler(_ *discordgo.Session, ready *discordgo.Ready) {
	l.SetUserID(ready.User.ID)
}

func (l *Link) VoiceStateUpdateHandler(_ *discordgo.Session, voiceStateUpdate discordgo.VoiceStateUpdate) {
	var channelID *string
	if voiceStateUpdate.ChannelID != "" {
		channelID = &voiceStateUpdate.ChannelID
	}
	l.VoiceStateUpdate(lavalink.VoiceStateUpdate{
		GuildID:   voiceStateUpdate.GuildID,
		ChannelID: channelID,
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *Link) VoiceServerUpdateHandler(_ *discordgo.Session, voiceServerUpdate *discordgo.VoiceServerUpdate) {
	var endpoint *string
	if voiceServerUpdate.Endpoint != "" {
		endpoint = &voiceServerUpdate.Endpoint
	}
	l.VoiceServerUpdate(lavalink.VoiceServerUpdate{
		GuildID:  voiceServerUpdate.GuildID,
		Token:    voiceServerUpdate.Token,
		Endpoint: endpoint,
	})
}
