package arikawalink

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

func New(s *session.Session, opts ...lavalink.ConfigOpt) *Link {
	link := &Link{
		Lavalink: lavalink.New(opts...),
	}

	if link.UserID() == "" {
		user, err := s.Me()
		if err != nil {
			link.Logger().Errorf("Failed to get user ID: %s", err)
		} else {
			link.SetUserID(snowflake.ParseString(user.ID))
		}
	}

	s.AddHandler(link.VoiceStateUpdateHandler)
	s.AddHandler(link.VoiceServerUpdateHandler)
	return link
}

var (
	_ lavalink.Lavalink = (*Link)(nil)
)

type Link struct {
	lavalink.Lavalink
}

func (l *Link) VoiceStateUpdateHandler(e *gateway.VoiceStateUpdateEvent) {
	var channelID *snowflake.Snowflake
	if e.ChannelID != 0 {
		cid := snowflake.ParseString(e.ChannelID)
		channelID = &cid
	}
	l.VoiceStateUpdate(lavalink.VoiceStateUpdate{
		GuildID:   snowflake.ParseString(e.GuildID),
		ChannelID: channelID,
		SessionID: e.SessionID,
	})
}

func (l *Link) VoiceServerUpdateHandler(e *gateway.VoiceServerUpdateEvent) {
	var endpoint *string
	if e.Endpoint != "" {
		endpoint = &e.Endpoint
	}
	l.VoiceServerUpdate(lavalink.VoiceServerUpdate{
		GuildID:  snowflake.ParseString(e.GuildID),
		Token:    e.Token,
		Endpoint: endpoint,
	})
}
