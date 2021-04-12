package player

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/filters"
)

type Player interface {
	PlayingTrack() *api.Track
	PlayTrack(track *api.Track)
	StopTrack()
	SetPaused(paused bool)
	Paused() bool
	TrackPosition() int
	SeekTo(position int)
	Filters() *filters.Filters
	AddListener(playerListener Listener)
	RemoveListener(playerListener Listener)
	EmitEvent(playerEvent Event)
	Link() api.Link
}