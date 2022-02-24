package lavalink

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/snowflake"
	"github.com/pkg/errors"
)

type Player interface {
	PlayingTrack() AudioTrack
	Paused() bool
	Position() Duration
	Connected() bool
	Volume() int
	Filters() Filters
	GuildID() snowflake.Snowflake
	ChannelID() *snowflake.Snowflake
	Node() Node
	Export() PlayerRestoreState

	Play(track AudioTrack) error
	PlayAt(track AudioTrack, start Duration, end Duration) error
	Stop() error
	Destroy() error
	Pause(paused bool) error
	Seek(position Duration) error
	SetVolume(volume int) error
	SetFilters(filters Filters)
	ChangeNode(node Node)

	OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)
	OnPlayerUpdate(state PlayerState)

	EmitEvent(caller func(l interface{}))
	AddListener(listener interface{})
	RemoveListener(listener interface{})
}

var _ Player = (*DefaultPlayer)(nil)

func NewPlayer(node Node, lavalink Lavalink, guildID snowflake.Snowflake) Player {
	return &DefaultPlayer{
		guildID:  guildID,
		volume:   100,
		node:     node,
		lavalink: lavalink,
	}
}

func newResumingPlayer(node Node, lavalink Lavalink, resumeState PlayerRestoreState) (Player, error) {
	track, err := lavalink.DecodeTrack(resumeState.PlayingTrack)
	if err != nil {
		return nil, err
	}
	player := &DefaultPlayer{
		guildID:               resumeState.GuildID,
		channelID:             resumeState.ChannelID,
		lastSessionID:         resumeState.LastSessionID,
		lastVoiceServerUpdate: resumeState.LastVoiceServerUpdate,
		track:                 track,
		volume:                resumeState.Volume,
		paused:                resumeState.Paused,
		state:                 resumeState.State,
		filters:               resumeState.Filters,
		node:                  node,
		lavalink:              lavalink,
	}
	player.Filters().setCommitFunc(player.commitFilters)
	return player, nil
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
	lavalink              Lavalink
	listeners             []interface{}
}

func (p *DefaultPlayer) PlayingTrack() AudioTrack {
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
		return errors.Wrap(err, "error while playing track")
	}
	p.track = track
	return nil
}

func (p *DefaultPlayer) Stop() error {
	p.track = nil

	if p.node == nil {
		return nil
	}

	if err := p.node.Send(StopCommand{GuildID: p.guildID}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	return nil
}

func (p *DefaultPlayer) Destroy() error {
	for _, pl := range p.Node().Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnDestroyPlayer(p)
		}
	}
	if p.node != nil {
		if err := p.node.Send(DestroyCommand{GuildID: p.guildID}); err != nil {
			return errors.Wrap(err, "error while destroying player")
		}
	}
	p.lavalink.RemovePlayer(p.guildID)
	return nil
}

func (p *DefaultPlayer) Pause(pause bool) error {
	if p.paused == pause {
		return nil
	}
	if p.node != nil {
		if err := p.node.Send(PauseCommand{
			GuildID: p.guildID,
			Pause:   pause,
		}); err != nil {
			return errors.Wrap(err, "error while pausing player")
		}
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
	if p.PlayingTrack() == nil {
		return errors.New("no track is playing")
	}
	if p.PlayingTrack().Info().IsStream {
		return errors.New("cannot seek streams")
	}
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
	if p.node == nil {
		return nil
	}
	if volume < 0 {
		volume = 0
	}
	if volume > 1000 {
		volume = 1000
	}
	if err := p.node.Send(VolumeCommand{
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

func (p *DefaultPlayer) Export() PlayerRestoreState {
	return PlayerRestoreState{}
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
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
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
	if p.node == nil {
		return nil
	}
	if err := p.node.Send(FiltersCommand{
		GuildID: p.guildID,
		Filters: filters,
	}); err != nil {
		return errors.Wrap(err, "error while setting filters of player")
	}
	return nil
}

type PlayerState struct {
	Time      Time     `json:"time"`
	Position  Duration `json:"position"`
	Connected bool     `json:"connected"`
}

type PlayerRestoreState struct {
	PlayingTrack          string               `json:"playing_track"`
	Paused                bool                 `json:"paused"`
	State                 PlayerState          `json:"state"`
	Volume                int                  `json:"volume"`
	Filters               Filters              `json:"filters"`
	GuildID               snowflake.Snowflake  `json:"guild_id"`
	ChannelID             *snowflake.Snowflake `json:"channel_id"`
	LastSessionID         *string              `json:"last_session_id"`
	LastVoiceServerUpdate *VoiceServerUpdate   `json:"last_voice_server_update"`
	Node                  string               `json:"node"`
}

func (s *PlayerRestoreState) UnmarshalJSON(data []byte) error {
	type playerRestoreState PlayerRestoreState
	var v struct {
		Filters json.RawMessage `json:"filters"`
		playerRestoreState
	}
	var err error
	if err = json.Unmarshal(data, &v); err != nil {
		return errors.Wrap(err, "error while unmarshalling player resume state")
	}
	*s = PlayerRestoreState(v.playerRestoreState)
	s.Filters, err = UnmarshalFilters(v.Filters)
	return err
}

var UnmarshalFilters = func(data []byte) (Filters, error) {
	var filters *DefaultFilters
	if err := json.Unmarshal(data, &filters); err != nil {
		return nil, errors.Wrap(err, "error while unmarshalling filters")
	}
	return filters, nil
}
