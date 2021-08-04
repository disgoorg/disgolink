package internal

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/filters"
)

func NewPlayer(node api.Node, guildID dapi.Snowflake) api.Player {
	return &PlayerImpl{
		guildID:       guildID,
		channelID:     nil,
		lastSessionID: nil,
		track:         nil,
		volume:        100,
		paused:        false,
		position:      -1,
		connected:     false,
		updateTime:    -1,
		filters:       nil,
		node:          node,
		listeners:     nil,
	}
}

type PlayerImpl struct {
	guildID       dapi.Snowflake
	channelID     *dapi.Snowflake
	lastSessionID *string
	track         api.Track
	volume        int
	paused        bool
	position      int
	connected     bool
	updateTime    int
	filters       *filters.Filters
	node          api.Node
	listeners     []api.PlayerEventListener
}

func (p *PlayerImpl) Track() api.Track {
	return p.track
}

func (p *PlayerImpl) SetTrack(track api.Track) {
	p.track = track
}

func (p *PlayerImpl) Play(track api.Track) {
	t := track.Track()
	if track == nil {
		p.Node().Lavalink().Logger().Errorf("error while playing track: track base64 is nil")
		return
	}
	p.Node().Send(&api.PlayPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpPlay, p),
		Track:         *t,
	})
}

func (p *PlayerImpl) PlayAt(track api.Track, start int, end int) {
	t := track.Track()
	if track == nil {
		p.Node().Lavalink().Logger().Errorf("error while playing track: track base64 is nil")
		return
	}
	p.Node().Send(&api.PlayPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpPlay, p),
		Track:         *t,
		StartTime:     &start,
		EndTime:       &end,
	})
}

func (p *PlayerImpl) Stop() {
	p.track = nil

	p.Node().Send(&api.StopPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpStop, p),
	})
}

func (p *PlayerImpl) Destroy() {
	p.track = nil

	p.Node().Send(&api.DestroyPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpDestroy, p),
	})
}

func (p *PlayerImpl) Pause(pause bool) {
	if p.paused == pause {
		return
	}
	p.Node().Send(&api.PausePlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpPause, p),
		Pause:         pause,
	})
	p.paused = pause
	if pause {
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnPlayerPause(p)
		})
	} else {
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnPlayerResume(p)
		})
	}
}

func (p *PlayerImpl) Paused() bool {
	return p.paused
}

func (p *PlayerImpl) Position() int {
	return p.position
}

func (p *PlayerImpl) Seek(position int) {
	p.Node().Send(&api.SeekPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpSeek, p),
		Position:      position,
	})
}

func (p *PlayerImpl) Volume() int {
	return p.volume
}

func (p *PlayerImpl) SetVolume(volume int) {
	if volume < 0 {
		volume = 0
	}
	if volume > 1000 {
		volume = 1000
	}
	p.Node().Send(&api.VolumePlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpSeek, p),
		Volume:        volume,
	})
	p.volume = volume
}

func (p *PlayerImpl) Filters() *filters.Filters {
	if p.filters == nil {
		p.filters = filters.NewFilters(p.commitFilters)
	}
	return p.filters
}

func (p *PlayerImpl) SetFilters(filters *filters.Filters) {
	p.filters = filters
}

func (p *PlayerImpl) GuildID() dapi.Snowflake {
	return p.guildID
}

func (p *PlayerImpl) ChannelID() *dapi.Snowflake {
	return p.channelID
}

func (p *PlayerImpl) SetChannelID(channelID *dapi.Snowflake) {
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

func (p *PlayerImpl) PlayerUpdate(state api.State) {
	p.position = state.Position
	p.connected = state.Connected
}

func (p *PlayerImpl) EmitEvent(listenerCaller func(listener api.PlayerEventListener)) {
	for _, listener := range p.listeners {
		listenerCaller(listener)
	}
}

func (p *PlayerImpl) AddListener(playerListener api.PlayerEventListener) {
	p.listeners = append(p.listeners, playerListener)
}

func (p *PlayerImpl) RemoveListener(playerListener api.PlayerEventListener) {
	for i, listener := range p.listeners {
		if listener == playerListener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}

func (p *PlayerImpl) commitFilters(filters *filters.Filters) {
	p.node.Send(&api.FilterPlayerCommand{
		PlayerCommand: api.NewPlayerCommand(api.OpFilters, p),
		Filters:       filters,
	})
}
