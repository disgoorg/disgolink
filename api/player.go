package api

import (
	"github.com/DisgoOrg/disgolink/api/filters"
	"github.com/DisgoOrg/disgolink/api/player"
)

type Player interface {
	PlayingTrack() *Track
	PlayTrack(track *Track)
	StopTrack()
	SetPaused(paused bool)
	Resume()
	Paused() bool
	TrackPosition() int
	SeekTo(position int)
	Filters() *filters.Filters
	Commit()
	AddListener(playerListener player.Listener)
	RemoveListener(playerListener player.Listener)
	EmitEvent(playerEvent player.Event)
	GuildID() string
	ChannelID() *string
	SetChannelID(channelID *string)
	LastSessionID() *string
	SetLastSessionID(sessionID string)
	Node() Node
	ChangeNode(node Node)
}
