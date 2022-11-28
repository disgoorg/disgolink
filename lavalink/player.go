package lavalink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/disgoorg/disgolink/lavalink/protocol"
	"github.com/disgoorg/snowflake/v2"
)

var ErrPlayerNoNode = errors.New("player has no node")

type Player interface {
	GuildID() snowflake.ID
	ChannelID() *snowflake.ID
	Track() protocol.Track
	Paused() bool
	State() protocol.PlayerState
	Volume() int
	Filters() protocol.Filters

	Update(ctx context.Context, update protocol.PlayerUpdate) error
	Destroy() error

	Lavalink() Lavalink
	Node() Node

	EmitEvent(caller func(l any))
	AddListener(listener any)
	RemoveListener(listener any)

	OnEvent(event protocol.Event)
	OnPlayerUpdate(playerUpdate protocol.PlayerUpdate)
	OnVoiceServerUpdate(token string, endpoint string)
	OnVoiceStateUpdate(sessionID string)
}

type defaultPlayer struct {
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

func (p *defaultPlayer) Track() protocol.Track {
	return p.track
}

func (p *defaultPlayer) Destroy() error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	return p.node.Rest().
}

func (p *defaultPlayer) Pause(pause bool) error {
	if p.paused == pause {
		return nil
	}
	if p.node != nil {
		if err := p.node.Send(PauseCommand{
			GuildID: p.guildID,
			Pause:   pause,
		}); err != nil {
			return fmt.Errorf("error while pausing player: %w", err)
		}
	}

	p.paused = pause
	if pause {
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

func (p *defaultPlayer) Paused() bool {
	return p.paused
}

func (p *defaultPlayer) Position() protocol.Duration {
	if p.track == nil {
		return 0
	}
	position := p.state.Position
	if !p.paused {
		position += protocol.Duration(time.Now().UnixMilli() - p.state.Time.UnixMilli())
	}
	if position > p.track.Info().Length {
		return p.track.Info().Length
	}
	if position < 0 {
		return 0
	}
	return position
}

func (p *defaultPlayer) Connected() bool {
	return p.state.Connected
}

func (p *defaultPlayer) Seek(position protocol.Duration) error {
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
		return fmt.Errorf("error while seeking player: %w", err)
	}
	p.state.Position = position
	p.state.Time = protocol.Time{Time: time.Now()}
	return nil
}

func (p *defaultPlayer) Volume() int {
	return p.volume
}

func (p *defaultPlayer) SetVolume(volume int) error {
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
		return fmt.Errorf("error while setting volume of player: %w", err)
	}
	p.volume = volume
	return nil
}

func (p *defaultPlayer) Filters() Filters {
	if p.filters == nil {
		p.filters = NewFilters(p.commitFilters)
	}
	return p.filters
}

func (p *defaultPlayer) SetFilters(filters Filters) {
	p.filters = filters
}

func (p *defaultPlayer) GuildID() snowflake.ID {
	return p.guildID
}

func (p *defaultPlayer) ChannelID() *snowflake.ID {
	return p.channelID
}

func (p *defaultPlayer) Export() PlayerRestoreState {
	return PlayerRestoreState{}
}

func (p *defaultPlayer) OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	p.lastVoiceServerUpdate = &voiceServerUpdate
	if err := p.Node().Send(VoiceUpdateCommand{
		GuildID:   voiceServerUpdate.GuildID,
		SessionID: *p.lastSessionID,
		Event:     voiceServerUpdate,
	}); err != nil {
		p.node.Lavalink().Logger().Error("error while sending voice server update: ", err)
	}
}

func (p *defaultPlayer) OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
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

func (p *defaultPlayer) Node() Node {
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
	return p.node
}

func (p *defaultPlayer) ChangeNode(node Node) {
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

func (p *defaultPlayer) OnPlayerUpdate(state PlayerState) {
	p.state = state
}

func (p *defaultPlayer) EmitEvent(caller func(l any)) {
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

func (p *defaultPlayer) AddListener(listener any) {
	p.listeners = append(p.listeners, listener)
}

func (p *defaultPlayer) RemoveListener(listener any) {
	for i, l := range p.listeners {
		if l == listener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}

func (p *defaultPlayer) OnEvent(event TrackEvent) {
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
		if e.Reason != protocol.AudioTrackEndReasonReplaced && e.Reason != protocol.AudioTrackEndReasonStopped {
			p.track = nil
		}
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

func (p *defaultPlayer) commitFilters(filters Filters) error {
	if p.node == nil {
		return nil
	}
	if err := p.node.Send(FiltersCommand{
		GuildID: p.guildID,
		Filters: filters,
	}); err != nil {
		return fmt.Errorf("error while setting filters of player: %w", err)
	}
	return nil
}

type PlayerState struct {
	Time      protocol.Time     `json:"time"`
	Position  protocol.Duration `json:"position"`
	Connected bool              `json:"connected"`
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
