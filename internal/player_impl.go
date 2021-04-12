package internal

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/filters"
	"github.com/DisgoOrg/disgolink/api/player"
)

type PlayerImpl struct {
	track      *api.Track
	paused     bool
	updateTime int
	position   int
	filters    *filters.Filters
	link       api.Link
	listeners  []player.Listener
}

func (p *PlayerImpl) PlayingTrack() *api.Track {
	return p.track
}
func (p *PlayerImpl) PlayTrack(track *api.Track) {
	p.position = track.Position()

	p.link.Node().Send(&api.OpPlayPlayer{
		OpPlayerCommand: api.NewPlayerCommand(api.PlayOp, p.link.GuildID()),
		Track:           track.Encode(),
		StartTime:       p.position,
		Paused:          p.paused,
	})

}
func (p *PlayerImpl) StopTrack() {
	p.track = nil

	p.link.Node().Send(&api.OpStopPlayer{
		OpPlayerCommand: api.NewPlayerCommand(api.StopOp, p.link.GuildID()),
	})

}
func (p *PlayerImpl) SetPaused(paused bool) {
	if p.paused == paused {
		return
	}
	p.link.Node().Send(&api.OpPausePlayer{
		OpPlayerCommand: api.NewPlayerCommand(api.PauseOP, p.link.GuildID()),
		Paused:          paused,
	})
	p.paused = paused
}
func (p *PlayerImpl) Paused() bool {
	return p.paused
}
func (p *PlayerImpl) TrackPosition() int {
	// TODO
	return 0
}
func (p *PlayerImpl) SeekTo(position int) {

}
func (p *PlayerImpl) Filters() *filters.Filters {
	return p.filters
}
func (p *PlayerImpl) AddListener(playerListener player.Listener) {
	p.listeners = append(p.listeners, playerListener)
}
func (p *PlayerImpl) RemoveListener(playerListener player.Listener) {
	for i, listener := range p.listeners {
		if listener == playerListener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}
func (p *PlayerImpl) EmitEvent(playerEvent player.Event) {
	for _, listener := range p.listeners {
		listener.OnEvent(playerEvent)
	}
}
func (p *PlayerImpl) Link() api.Link {
	return p.link
}
