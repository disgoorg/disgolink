package dgolink

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake/v2"
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
	userID, err := snowflake.Parse(voiceStateUpdate.UserID)
	if err != nil {
		l.Logger().Error("Failed to parse user ID: ", err)
		return
	}
	if userID != l.UserID() {
		return
	}
	var channelID *snowflake.ID
	if voiceStateUpdate.ChannelID != "" {
		id, err := snowflake.Parse(voiceStateUpdate.ChannelID)
		if err != nil {
			l.Logger().Error("Failed to parse channel ID: ", err)
			return
		}
		channelID = &id
	}
	guildID, err := snowflake.Parse(voiceStateUpdate.GuildID)
	if err != nil {
		l.Logger().Error("Failed to parse guild ID: ", err)
		return
	}
	l.OnVoiceStateUpdate(lavalink.VoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: channelID,
		SessionID: voiceStateUpdate.SessionID,
	})
}

func (l *Link) OnVoiceServerUpdateHandler(_ *discordgo.Session, voiceServerUpdate *discordgo.VoiceServerUpdate) {
	var endpoint *string
	if voiceServerUpdate.Endpoint != "" {
		endpoint = &voiceServerUpdate.Endpoint
	}

	guildID, err := snowflake.Parse(voiceServerUpdate.GuildID)
	if err != nil {
		l.Logger().Error("Failed to parse guild ID: ", err)
		return
	}
	l.OnVoiceServerUpdate(lavalink.VoiceServerUpdate{
		GuildID:  guildID,
		Token:    voiceServerUpdate.Token,
		Endpoint: endpoint,
	})
}
