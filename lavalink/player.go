package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/filters"
)

type Player interface {
	Track() Track
	SetTrack(track Track)
	Play(track Track) error
	PlayAt(track Track, start int, end int) error
	Stop() error
	Destroy() error
	Pause(paused bool) error
	Paused() bool
	Position() int
	Seek(position int) error
	Volume() int
	SetVolume(volume int) error
	Filters() *filters.Filters
	SetFilters(filters *filters.Filters)

	GuildID() discord.Snowflake
	ChannelID() *discord.Snowflake
	SetChannelID(channelID *discord.Snowflake)
	LastSessionID() *string
	SetLastSessionID(sessionID string)

	Node() Node
	ChangeNode(node Node)

	PlayerUpdate(state State)
	EmitEvent(listenerCaller func(listener PlayerEventListener))
	AddListener(playerListener PlayerEventListener)
	RemoveListener(playerListener PlayerEventListener)
}
