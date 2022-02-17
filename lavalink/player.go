package lavalink

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/snowflake"
	"github.com/pkg/errors"
)

type Player interface {
	Track() AudioTrack
	Play(track AudioTrack) error
	PlayAt(track AudioTrack, start Duration, end Duration) error
	Stop() error
	Destroy() error
	Pause(paused bool) error
	Paused() bool
	Position() Duration
	Connected() bool
	Seek(position Duration) error
	Volume() int
	SetVolume(volume int) error
	Filters() Filters
	SetFilters(filters Filters)

	GuildID() snowflake.Snowflake
	ChannelID() *snowflake.Snowflake

	OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)

	Node() Node
	ChangeNode(node Node)

	OnPlayerUpdate(state PlayerState)
	EmitEvent(caller func(l interface{}))
	AddListener(listener interface{})
	RemoveListener(listener interface{})
}

type PlayerState struct {
	Time      time.Time `json:"time"`
	Position  Duration  `json:"position"`
	Connected bool      `json:"connected"`
}

func (s *PlayerState) UnmarshalJSON(data []byte) error {
	type playerState PlayerState
	var v struct {
		Time int64 `json:"time"`
		playerState
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*s = PlayerState(v.playerState)
	s.Time = time.UnixMilli(v.Time)
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
	guildID               snowflake.Snowflake
	channelID             *snowflake.Snowflake
	lastSessionID         *string
	lastVoiceServerUpdate *VoiceServerUpdate
	track                 AudioTrack
	volume                int
	paused                bool
	state                 PlayerState
	filters               Filters
	node                  Node
	listeners             []interface{}
}

func (p *DefaultPlayer) Track() AudioTrack {
	return p.track
}

func (p *DefaultPlayer) Play(track AudioTrack) error {
	encodedTrack, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err = p.Node().Send(PlayCommand{
		GuildID: p.guildID,
		Track:   encodedTrack,
	}); err != nil {
		return errors.Wrap(err, "error while playing track")
	}
	p.track = track
	return nil
}

func (p *DefaultPlayer) PlayAt(track AudioTrack, start Duration, end Duration) error {
	encodedTrack, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err = p.Node().Send(PlayCommand{
		GuildID:   p.guildID,
		Track:     encodedTrack,
		StartTime: &start,
		EndTime:   &end,
	}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	p.track = track
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
		return nil
	}
	p.EmitEvent(func(l interface{}) {
		if listener, ok := l.(PlayerEventListener); ok {
			listener.OnPlayerResume(p)
		}
	})
	return nil
}

func (p *DefaultPlayer) Paused() bool {
	return p.paused
}

func (p *DefaultPlayer) Position() Duration {
	if p.track == nil {
		return 0
	}
	position := p.state.Position
	if !p.paused {
		position += Duration(time.Now().UnixMilli() - p.state.Time.UnixMilli())
	}
	if position > p.track.Info().Length {
		return p.track.Info().Length
	}
	if position < 0 {
		return 0
	}
	return position
}

func (p *DefaultPlayer) Connected() bool {
	return p.state.Connected
}

func (p *DefaultPlayer) Seek(position Duration) error {
	if err := p.Node().Send(SeekCommand{
		GuildID:  p.guildID,
		Position: position,
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

func (p *DefaultPlayer) OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	p.lastVoiceServerUpdate = &voiceServerUpdate
	if err := p.Node().Send(VoiceUpdateCommand{
		GuildID:   voiceServerUpdate.GuildID,
		SessionID: *p.lastSessionID,
		Event:     voiceServerUpdate,
	}); err != nil {
		p.node.Lavalink().Logger().Error("error while sending voice server update: ", err)
	}
}

func (p *DefaultPlayer) OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
	if voiceStateUpdate.ChannelID == nil {
		p.channelID = nil
		if p.Node() != nil {
			if err := p.Destroy(); err != nil {
				p.node.Lavalink().Logger().Error("error while destroying player: ", err)
			}
			p.node = nil
		}
		return
	}
	p.channelID = voiceStateUpdate.ChannelID
	p.lastSessionID = &voiceStateUpdate.SessionID
}

func (p *DefaultPlayer) Node() Node {
	return p.node
}

func (p *DefaultPlayer) ChangeNode(node Node) {
	p.node = node
	if p.lastVoiceServerUpdate != nil {
		p.OnVoiceServerUpdate(*p.lastVoiceServerUpdate)
		if track := p.track; track != nil {
			track.SetPosition(p.Position())
			if err := p.Play(track); err != nil {
				p.node.Lavalink().Logger().Error("error while changing node and resuming track: ", err)
			}
		}
	}
}

func (p *DefaultPlayer) OnPlayerUpdate(state PlayerState) {
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
