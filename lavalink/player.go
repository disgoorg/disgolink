package lavalink

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type Player interface {
	PlayingTrack() AudioTrack
	Paused() bool
	Position() Duration
	Connected() bool
	Volume() int
	Filters() Filters
	GuildID() snowflake.ID
	ChannelID() *snowflake.ID
	Node() Node
	Export() PlayerRestoreState

	Play(track AudioTrack) error
	PlayTrack(track AudioTrack, options PlayOptions) error
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

	EmitEvent(caller func(l any))
	AddListener(listener any)
	RemoveListener(listener any)
	OnEvent(event TrackEvent)
}

type PlayOptions struct {
	StartTime Duration
	EndTime   Duration
	NoReplace bool
	Paused    bool
	Volume    int
}

var _ Player = (*DefaultPlayer)(nil)

func NewPlayer(node Node, lavalink Lavalink, guildID snowflake.ID) Player {
	return &DefaultPlayer{
		guildID:  guildID,
		volume:   100,
		node:     node,
		lavalink: lavalink,
	}
}

func newResumingPlayer(node Node, lavalink Lavalink, resumeState PlayerRestoreState) (Player, error) {
	var track AudioTrack
	if resumeState.PlayingTrack != nil {
		var err error
		if track, err = lavalink.DecodeTrack(*resumeState.PlayingTrack); err != nil {
			return nil, err
		}
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

	if resumeState.Filters != nil {
		resumeState.Filters.setCommitFunc(player.commitFilters)
	}

	return player, nil
}

type DefaultPlayer struct {
	guildID               snowflake.ID
	channelID             *snowflake.ID
	lastSessionID         *string
	lastVoiceServerUpdate *VoiceServerUpdate
	track                 AudioTrack
	volume                int
	paused                bool
	state                 PlayerState
	filters               Filters
	node                  Node
	lavalink              Lavalink
	listeners             []any
}

func (p *DefaultPlayer) PlayingTrack() AudioTrack {
	return p.track
}

func (p *DefaultPlayer) PlayTrack(track AudioTrack, options PlayOptions) error {
	encodedTrack, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	payload := PlayerUpdate{
		EncodedTrack: json.NewNullablePtr(encodedTrack),
	}
	if options.StartTime != 0 {
		payload.Position = &options.StartTime
	}
	if options.EndTime != 0 {
		payload.EndTime = &options.EndTime
	}
	if options.Paused {
		payload.Paused = &options.Paused
	}
	if options.Volume != 0 {
		payload.Volume = &options.Volume
	}

	if err = p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, payload, options.NoReplace); err != nil {
		return fmt.Errorf("error while playing track: %w", err)
	}
	p.track = track
	p.paused = options.Paused
	if options.Volume != 0 {
		p.volume = options.Volume
	}
	return nil
}

func (p *DefaultPlayer) Play(track AudioTrack) error {
	encodedTrack, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err = p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		EncodedTrack: json.NewNullablePtr(encodedTrack),
	}, false); err != nil {
		return fmt.Errorf("error while playing track: %w", err)
	}
	p.track = track
	return nil
}

func (p *DefaultPlayer) PlayAt(track AudioTrack, start Duration, end Duration) error {
	encodedTrack, err := p.node.Lavalink().EncodeTrack(track)
	if err != nil {
		return err
	}

	if err = p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		EncodedTrack: json.NewNullablePtr(encodedTrack),
		Position:     &start,
		EndTime:      &end,
	}, false); err != nil {
		return fmt.Errorf("error while playing track: %w", err)
	}
	p.track = track
	return nil
}

func (p *DefaultPlayer) Stop() error {
	p.track = nil
	if p.node == nil {
		return nil
	}

	if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		EncodedTrack: json.NullPtr[string](),
	}, false); err != nil {
		return fmt.Errorf("error while stopping player: %w", err)
	}
	return nil
}

func (p *DefaultPlayer) Destroy() error {
	for _, pl := range p.node.Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnDestroyPlayer(p)
		}
	}
	if p.node != nil {
		if err := p.node.RestClient().DestroyPlayer(context.TODO(), p.guildID); err != nil {
			return fmt.Errorf("error while destroying player: %w", err)
		}
	}
	p.lavalink.RemovePlayer(p.guildID)
	return nil
}

func (p *DefaultPlayer) Pause(paused bool) error {
	if p.paused == paused {
		return nil
	}
	if p.node != nil {
		if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
			Paused: &paused,
		}, false); err != nil {
			return fmt.Errorf("error while pausing player: %w", err)
		}
	}

	p.paused = paused
	if paused {
		p.EmitEvent(func(l any) {
			if listener, ok := l.(PlayerEventListener); ok {
				listener.OnPlayerPause(p)
			}
		})
		return nil
	}
	p.EmitEvent(func(l any) {
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
	if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		Position: &position,
	}, false); err != nil {
		return fmt.Errorf("error while seeking player: %w", err)
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
	if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		Volume: &volume,
	}, false); err != nil {
		return fmt.Errorf("error while setting volume of player: %w", err)
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

func (p *DefaultPlayer) GuildID() snowflake.ID {
	return p.guildID
}

func (p *DefaultPlayer) ChannelID() *snowflake.ID {
	return p.channelID
}

func (p *DefaultPlayer) Export() PlayerRestoreState {
	return PlayerRestoreState{}
}

func (p *DefaultPlayer) OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	if p.lastSessionID == nil {
		return
	}

	var endpoint string
	if voiceServerUpdate.Endpoint != nil {
		endpoint = *voiceServerUpdate.Endpoint
	}
	p.lastVoiceServerUpdate = &voiceServerUpdate
	if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		Voice: &VoiceState{
			Token:     voiceServerUpdate.Token,
			Endpoint:  endpoint,
			SessionID: *p.lastSessionID,
		},
	}, false); err != nil {
		p.node.Lavalink().Logger().Error("error while sending voice server update: ", err)
	}
}

func (p *DefaultPlayer) OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
	if voiceStateUpdate.ChannelID == nil {
		p.channelID = nil
		if p.node != nil {
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

func (p *DefaultPlayer) EmitEvent(caller func(l any)) {
	defer func() {
		if r := recover(); r != nil {
			p.lavalink.Logger().Errorf("recovered from panic in event listener: %#v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	for _, listener := range p.listeners {
		caller(listener)
	}
}

func (p *DefaultPlayer) AddListener(listener any) {
	p.listeners = append(p.listeners, listener)
}

func (p *DefaultPlayer) RemoveListener(listener any) {
	for i, l := range p.listeners {
		if l == listener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}

func (p *DefaultPlayer) OnEvent(event TrackEvent) {
	track, err := p.node.Lavalink().DecodeTrack(event.Track())
	if err != nil {
		p.node.Lavalink().Logger().Errorf("error while decoding track: %s", err)
		return
	}
	if playingTrack := p.track; playingTrack != nil {
		track.SetUserData(playingTrack.UserData())
	}

	switch e := event.(type) {
	case TrackStartEvent:
		p.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnTrackStart(p, track)
			}
		})

	case TrackEndEvent:
		p.track = nil
		p.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnTrackEnd(p, track, e.Reason)
			}
		})

	case TrackExceptionEvent:
		p.track = nil
		p.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnTrackException(p, track, e.Exception)
			}
		})

	case TrackStuckEvent:
		p.track = nil
		p.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnTrackStuck(p, track, e.ThresholdMs)
			}
		})
	}
}

func (p *DefaultPlayer) commitFilters(filters Filters) error {
	if p.node == nil {
		return nil
	}
	if err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		Filters: &filters,
	}, false); err != nil {
		return fmt.Errorf("error while setting filters of player: %w", err)
	}
	return nil
}

type PlayerState struct {
	Time      Time     `json:"time"`
	Position  Duration `json:"position"`
	Connected bool     `json:"connected"`
}

type PlayerRestoreState struct {
	PlayingTrack          *string            `json:"playing_track"`
	Paused                bool               `json:"paused"`
	State                 PlayerState        `json:"state"`
	Volume                int                `json:"volume"`
	Filters               Filters            `json:"filters"`
	GuildID               snowflake.ID       `json:"guild_id"`
	ChannelID             *snowflake.ID      `json:"channel_id"`
	LastSessionID         *string            `json:"last_session_id"`
	LastVoiceServerUpdate *VoiceServerUpdate `json:"last_voice_server_update"`
	Node                  string             `json:"node"`
}

func (s *PlayerRestoreState) UnmarshalJSON(data []byte) error {
	type playerRestoreState PlayerRestoreState
	var v struct {
		Filters json.RawMessage `json:"filters"`
		playerRestoreState
	}
	var err error
	if err = json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("error while unmarshalling player resume state: %w", err)
	}
	*s = PlayerRestoreState(v.playerRestoreState)
	s.Filters, err = UnmarshalFilters(v.Filters)
	return err
}

var UnmarshalFilters = func(data []byte) (Filters, error) {
	var filters *DefaultFilters
	if err := json.Unmarshal(data, &filters); err != nil {
		return nil, fmt.Errorf("error while unmarshalling filters: %w", err)
	}
	return filters, nil
}
