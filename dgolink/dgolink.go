package dgolink

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/bwmarrin/discordgo"
)

func New(s *discordgo.Session, opts ...lavalink.ConfigOpt) *Link {
	discordgolink := &Link{
		Lavalink: lavalink.New(append([]lavalink.ConfigOpt{lavalink.WithHTTPClient(s.Client), lavalink.WithUserIDFromBotToken(s.Token)}, opts...)...),
	}

	s.AddHandler(discordgolink.VoiceStateUpdateHandler)
	s.AddHandler(discordgolink.VoiceServerUpdateHandler)
	return discordgolink
}

var _ lavalink.Lavalink = (*Link)(nil)

type Link struct {
	lavalink.Lavalink
}

func (l *Link) VoiceStateUpdateHandler(_ *discordgo.Session, voiceStateUpdate *discordgo.VoiceStateUpdate) {
	var channelID *string
	if voiceStateUpdate.ChannelID != "" {
		channelID = &voiceStateUpdate.ChannelID
	}
	l.OnVoiceStateUpdate(lavalink.VoiceStateUpdate{
		GuildID:   snowflake.Snowflake(voiceStateUpdate.GuildID),
		ChannelID: (*snowflake.Snowflake)(channelID),
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *Link) VoiceServerUpdateHandler(_ *discordgo.Session, voiceServerUpdate *discordgo.VoiceServerUpdate) {
	var endpoint *string
	if voiceServerUpdate.Endpoint != "" {
		endpoint = &voiceServerUpdate.Endpoint
	}
	l.OnVoiceServerUpdate(lavalink.VoiceServerUpdate{
		GuildID:  snowflake.Snowflake(voiceServerUpdate.GuildID),
		Token:    voiceServerUpdate.Token,
		Endpoint: endpoint,
	})
}
