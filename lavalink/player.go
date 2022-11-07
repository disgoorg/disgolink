package lavalink

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Player interface {
	GuildID() snowflake.ID
	Track() AudioTrack
	Volume() int
	Paused() bool
	State() PlayerState
	Position() Duration
	VoiceState() VoiceState
	Filters() *Filters

	Node() Node

	Update(ctx context.Context, update PlayerUpdate, noReplace bool) error
	Destroy(ctx context.Context) error

	OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)
	OnPlayerUpdate(state PlayerState)

	EmitEvent(caller func(l any))
	AddListener(listener any)
	RemoveListener(listener any)
	OnEvent(event TrackEvent)
}

var _ Player = (*defaultPlayer)(nil)

func NewPlayer(node Node, lavalink Lavalink, guildID snowflake.ID) Player {
	return &defaultPlayer{
		guildID:  guildID,
		volume:   100,
		node:     node,
		lavalink: lavalink,
	}
}

type defaultPlayer struct {
	guildID    snowflake.ID
	voiceState VoiceState
	track      AudioTrack
	volume     int
	paused     bool
	state      PlayerState
	filters    *Filters
	node       Node
	lavalink   Lavalink
	listeners  []any
}

func (p *defaultPlayer) GuildID() snowflake.ID {
	return p.guildID
}

func (p *defaultPlayer) Track() AudioTrack {
	return p.track
}

func (p *defaultPlayer) Volume() int {
	return p.volume
}

func (p *defaultPlayer) Paused() bool {
	return p.paused
}

func (p *defaultPlayer) State() PlayerState {
	return p.state
}

func (p *defaultPlayer) Position() Duration {
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

func (p *defaultPlayer) VoiceState() VoiceState {
	return p.voiceState
}

func (p *defaultPlayer) Filters() *Filters {
	if p.filters == nil {
		p.filters = new(Filters)
	}
	return p.filters
}

func (p *defaultPlayer) Update(ctx context.Context, playerUpdate PlayerUpdate, noReplace bool) error {
	_, err := p.node.RestClient().UpdatePlayer(ctx, p.guildID, playerUpdate, noReplace)
	if err != nil {
		return err
	}
	// TODO: update local state

	return nil
}

func (p *defaultPlayer) Destroy(ctx context.Context) error {
	for _, pl := range p.node.Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnDestroyPlayer(p)
		}
	}
	if p.node != nil {
		if err := p.node.RestClient().DestroyPlayer(context.TODO(), p.guildID); err != nil {
			return fmt.Errorf("error while destroying defaultPlayer: %w", err)
		}
	}
	p.lavalink.RemovePlayer(p.guildID)
	return nil
}

func (p *defaultPlayer) OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	p.voiceState.Token = voiceServerUpdate.Token
	if voiceServerUpdate.Endpoint != nil {
		p.voiceState.Endpoint = *voiceServerUpdate.Endpoint
	}

	if p.voiceState.SessionID == "" {
		return
	}

	if _, err := p.node.RestClient().UpdatePlayer(context.TODO(), p.guildID, PlayerUpdate{
		Voice: &p.voiceState,
	}, false); err != nil {
		p.node.Lavalink().Logger().Error("error while sending voice server update: ", err)
	}
}

func (p *defaultPlayer) OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
	if voiceStateUpdate.ChannelID == nil {
		if p.node != nil {
			if err := p.Destroy(context.TODO()); err != nil {
				p.node.Lavalink().Logger().Error("error while destroying defaultPlayer: ", err)
			}
			p.node = nil
		}
		return
	}
	p.voiceState.SessionID = voiceStateUpdate.SessionID
}

func (p *defaultPlayer) Node() Node {
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
	return p.node
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

type PlayerState struct {
	Time      Time     `json:"time"`
	Position  Duration `json:"position"`
	Connected bool     `json:"connected"`
}
