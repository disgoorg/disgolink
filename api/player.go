package api

import (
	"github.com/DisgoOrg/disgolink/api/events"
	"github.com/DisgoOrg/disgolink/api/filters"
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
	AddListener(playerListener events.PlayerEventListener)
	RemoveListener(playerListener events.PlayerEventListener)
	EmitEvent(playerEvent events.PlayerEvent)
	GuildID() string
	ChannelID() *string
	SetChannelID(channelID *string)
	LastSessionID() *string
	SetLastSessionID(sessionID string)
	Node() Node
	ChangeNode(node Node)
}
