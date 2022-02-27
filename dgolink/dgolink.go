package dgolink

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/bwmarrin/discordgo"
)

func New(s *discordgo.Session, opts ...lavalink.ConfigOpt) *Link {
	discordgolink := &Link{
		Lavalink: lavalink.New(append([]lavalink.ConfigOpt{
			lavalink.WithHTTPClient(s.Client),
			lavalink.WithUserIDFromBotToken(s.Token),
		}, opts...)...),
	}

	s.AddHandler(discordgolink.OnVoiceStateUpdateHandler)
	s.AddHandler(discordgolink.OnVoiceServerUpdateHandler)
	return discordgolink
}

var _ lavalink.Lavalink = (*Link)(nil)

type Link struct {
	lavalink.Lavalink
}

func (l *Link) OnVoiceStateUpdateHandler(_ *discordgo.Session, voiceStateUpdate *discordgo.VoiceStateUpdate) {
	if snowflake.Snowflake(voiceStateUpdate.UserID) != l.UserID() {
		return
	}
	var channelID *snowflake.Snowflake
	if voiceStateUpdate.ChannelID != "" {
		channelID = (*snowflake.Snowflake)(&voiceStateUpdate.ChannelID)
	}
	l.OnVoiceStateUpdate(lavalink.VoiceStateUpdate{
		GuildID:   snowflake.Snowflake(voiceStateUpdate.GuildID),
		ChannelID: channelID,
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *Link) OnVoiceServerUpdateHandler(_ *discordgo.Session, voiceServerUpdate *discordgo.VoiceServerUpdate) {
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
