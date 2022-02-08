package lavalink

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/snowflake"
	"github.com/pkg/errors"
)

type Player interface {
	Track() AudioTrack
	SetTrack(track AudioTrack)
	Play(track AudioTrack) error
	PlayAt(track AudioTrack, start time.Duration, end time.Duration) error
	Stop() error
	Destroy() error
	Pause(paused bool) error
	Paused() bool
	Position() time.Duration
	Seek(position time.Duration) error
	Volume() int
	SetVolume(volume int) error
	Filters() Filters
	SetFilters(filters Filters)

	GuildID() snowflake.Snowflake
	ChannelID() *snowflake.Snowflake
	SetChannelID(channelID *snowflake.Snowflake)
	LastSessionID() *string
	SetLastSessionID(sessionID string)

	Node() Node
	SetNode(node Node)

	PlayerUpdate(state PlayerState)
	EmitEvent(caller func(l interface{}))
	AddListener(listener interface{})
	RemoveListener(listener interface{})
}

type PlayerState struct {
	Time      time.Time     `json:"time"`
	Position  time.Duration `json:"position"`
	Connected bool          `json:"connected"`
}

func (s *PlayerState) UnmarshalJSON(data []byte) error {
	var v struct {
		Time      int64 `json:"time"`
		Position  int64 `json:"position"`
		Connected bool  `json:"connected"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.Time = time.Unix(v.Time, 0)
	s.Position = time.Duration(v.Position) * time.Millisecond
	s.Connected = v.Connected
	return nil
}

var _ Player = (*DefaultPlayer)(nil)

func NewPlayer(node Node, guildID snowflake.Snowflake) Player {
	return &DefaultPlayer{
		guildID: guildID,
		volume:  100,
		node:    node,
	}
}

type DefaultPlayer struct {
	guildID       snowflake.Snowflake
	channelID     *snowflake.Snowflake
	lastSessionID *string
	track         AudioTrack
	volume        int
	paused        bool
	state         PlayerState
	filters       Filters
	node          Node
	listeners     []interface{}
}

func (p *DefaultPlayer) Track() AudioTrack {
	return p.track
}

func (p *DefaultPlayer) SetTrack(track AudioTrack) {
	p.track = track
}

func (p *DefaultPlayer) Play(track AudioTrack) error {
	t, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err := p.Node().Send(PlayCommand{
		GuildID: p.guildID,
		Track:   t,
	}); err != nil {
		return errors.Wrap(err, "error while playing track")
	}
	return nil
}

func (p *DefaultPlayer) PlayAt(track AudioTrack, start time.Duration, end time.Duration) error {
	t, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err := p.Node().Send(PlayCommand{
		GuildID:   p.guildID,
		Track:     t,
		StartTime: start.Milliseconds(),
		EndTime:   end.Milliseconds(),
	}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	return nil
}

func (p *DefaultPlayer) Stop() error {
	p.track = nil

	if err := p.Node().Send(StopCommand{GuildID: p.guildID}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	return nil
}

func (p *DefaultPlayer) Destroy() error {
	p.track = nil

	if err := p.Node().Send(DestroyCommand{GuildID: p.guildID}); err != nil {
		return errors.Wrap(err, "error while destroying player")
	}
	for _, pl := range p.Node().Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnDestroyPlayer(p)
		}
	}
	return nil
}

func (p *DefaultPlayer) Pause(pause bool) error {
	if p.paused == pause {
		return nil
	}

	if err := p.Node().Send(PauseCommand{
		GuildID: p.guildID,
		Pause:   pause,
	}); err != nil {
		return errors.Wrap(err, "error while pausing player")
	}

	p.paused = pause
	if pause {
		p.EmitEvent(func(l interface{}) {
			if listener, ok := l.(PlayerEventListener); ok {
				listener.OnPlayerPause(p)
			}

		})
	} else {
		p.EmitEvent(func(l interface{}) {
			if listener, ok := l.(PlayerEventListener); ok {
				listener.OnPlayerResume(p)
			}
		})
	}
	return nil
}

func (p *DefaultPlayer) Paused() bool {
	return p.paused
}

func (p *DefaultPlayer) Position() time.Duration {
	if p.track == nil {
		return 0
	}
	if p.paused {
		timeDiff := time.Since(p.state.Time)
		if p.state.Position+timeDiff > p.track.Info().Length {
			return p.track.Info().Length
		}
		return p.state.Position + timeDiff
	}
	return p.state.Position
}

func (p *DefaultPlayer) Seek(position time.Duration) error {
	if err := p.Node().Send(SeekCommand{
		GuildID:  p.guildID,
		Position: position.Milliseconds(),
	}); err != nil {
		return errors.Wrap(err, "error while seeking player")
	}
	return nil
}

func (p *DefaultPlayer) Volume() int {
	return p.volume
}

func (p *DefaultPlayer) SetVolume(volume int) error {
	if volume < 0 {
		volume = 0
	}
	if volume > 1000 {
		volume = 1000
	}
	if err := p.Node().Send(VolumeCommand{
		GuildID: p.guildID,
		Volume:  volume,
	}); err != nil {
		return errors.Wrap(err, "error while setting volume of player")
	}
	p.volume = volume
	return nil
}

func (p *DefaultPlayer) Filters() Filters {
	if p.filters == nil {
		p.filters = NewFilters(p.commitFilters)
	}
	return p.filters
}

func (p *DefaultPlayer) SetFilters(filters Filters) {
	p.filters = filters
}

func (p *DefaultPlayer) GuildID() snowflake.Snowflake {
	return p.guildID
}

func (p *DefaultPlayer) ChannelID() *snowflake.Snowflake {
	return p.channelID
}

func (p *DefaultPlayer) SetChannelID(channelID *snowflake.Snowflake) {
	p.channelID = channelID
}

func (p *DefaultPlayer) LastSessionID() *string {
	return p.lastSessionID
}

func (p *DefaultPlayer) SetLastSessionID(sessionID string) {
	p.lastSessionID = &sessionID
}

func (p *DefaultPlayer) Node() Node {
	return p.node
}

func (p *DefaultPlayer) SetNode(node Node) {
	p.node = node
}

func (p *DefaultPlayer) PlayerUpdate(state PlayerState) {
	p.state = state
}

func (p *DefaultPlayer) EmitEvent(caller func(l interface{})) {
	for _, listener := range p.listeners {
		caller(listener)
	}
}

func (p *DefaultPlayer) AddListener(listener interface{}) {
	p.listeners = append(p.listeners, listener)
}

func (p *DefaultPlayer) RemoveListener(listener interface{}) {
	for i, l := range p.listeners {
		if l == listener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}

func (p *DefaultPlayer) commitFilters(filters Filters) error {
	if err := p.node.Send(FiltersCommand{
		GuildID: p.guildID,
		Filters: filters,
	}); err != nil {
		return errors.Wrap(err, "error while setting filters of player")
	}
	return nil
}
