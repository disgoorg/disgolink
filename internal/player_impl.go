package internal

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/events"
	"github.com/DisgoOrg/disgolink/api/filters"
)

func NewPlayer(node api.Node, guildID string) api.Player {
	return &PlayerImpl{
		guildID:       guildID,
		channelID:     nil,
		lastSessionID: nil,
		track:         nil,
		paused:        false,
		position:      -1,
		updateTime:    -1,
		filters:       nil,
		connected:     false,
		node:          node,
		listeners:     nil,
	}
}

type PlayerImpl struct {
	guildID       string
	channelID     *string
	lastSessionID *string
	track         *api.Track
	paused        bool
	position      int
	updateTime    int
	filters       *filters.Filters
	connected     bool
	node          api.Node
	listeners     []events.PlayerEventListener
}

func (p *PlayerImpl) GuildID() string {
	return p.guildID
}
func (p *PlayerImpl) ChannelID() *string {
	return p.channelID
}
func (p *PlayerImpl) SetChannelID(channelID *string) {
	p.channelID = channelID
}
func (p *PlayerImpl) LastSessionID() *string {
	return p.lastSessionID
}
func (p *PlayerImpl) SetLastSessionID(sessionID string) {
	p.lastSessionID = &sessionID
}
func (p *PlayerImpl) Node() api.Node {
	return p.node
}
func (p *PlayerImpl) ChangeNode(node api.Node) {
	p.node = node
}

func (p *PlayerImpl) PlayingTrack() *api.Track {
	return p.track
}
func (p *PlayerImpl) PlayTrack(track *api.Track) {
	p.position = track.Info.Position

	p.Node().Send(&api.PlayPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpPlay, p),
		Track:         track.Track,
		StartTime:     p.position,
		Paused:        p.paused,
	})

}
func (p *PlayerImpl) StopTrack() {
	p.track = nil

	p.Node().Send(&api.StopPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpStop, p),
	})

}
func (p *PlayerImpl) SetPaused(paused bool) {
	if p.paused == paused {
		return
	}
	p.Node().Send(&api.PausePlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpPause, p),
		Paused:        paused,
	})
	p.paused = paused
}

func (p *PlayerImpl) Resume() {
	p.SetPaused(false)
}

func (p *PlayerImpl) Paused() bool {
	return p.paused
}
func (p *PlayerImpl) TrackPosition() int {
	return p.position
}
func (p *PlayerImpl) SeekTo(position int) {
	p.Node().Send(&api.SeekPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpSeek, p),
		Position:      position,
	})
}
func (p *PlayerImpl) Filters() *filters.Filters {
	if p.filters == nil {
		p.filters = filters.NewFilters(p.commitFilters)
	}
	return p.filters
}
func (p *PlayerImpl) Commit() {
	if p.filters == nil {
		return
	}
	p.filters.Commit()
}

func (p *PlayerImpl) AddListener(playerListener events.PlayerEventListener) {
	p.listeners = append(p.listeners, playerListener)
}
func (p *PlayerImpl) RemoveListener(playerListener events.PlayerEventListener) {
	for i, listener := range p.listeners {
		if listener == playerListener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}
func (p *PlayerImpl) EmitEvent(playerEvent events.PlayerEvent) {
	for _, listener := range p.listeners {
		listener.OnEvent(playerEvent)
	}
}

func (p *PlayerImpl) commitFilters(filters *filters.Filters) {
	p.node.Send(&api.FilterPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpFilters, p),
		Filters: filters,
	})
}
