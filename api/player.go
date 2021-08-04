package api

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api/filters"
)

type Player interface {
	Track() Track
	SetTrack(track Track)
	Play(track Track)
	PlayAt(track Track, start int, end int)
	Stop()
	Destroy()
	Pause(paused bool)
	Paused() bool
	Position() int
	Seek(position int)
	Volume() int
	SetVolume(volume int)
	Filters() *filters.Filters
	SetFilters(filters *filters.Filters)

	GuildID() dapi.Snowflake
	ChannelID() *dapi.Snowflake
	SetChannelID(channelID *dapi.Snowflake)
	LastSessionID() *string
	SetLastSessionID(sessionID string)

	Node() Node
	ChangeNode(node Node)

	PlayerUpdate(state State)
	EmitEvent(listenerCaller func(listener PlayerEventListener))
	AddListener(playerListener PlayerEventListener)
	RemoveListener(playerListener PlayerEventListener)
}
