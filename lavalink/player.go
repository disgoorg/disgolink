package lavalink

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/disgoorg/disgolink/lavalink/protocol"
	"github.com/disgoorg/snowflake/v2"
)

var ErrPlayerNoNode = errors.New("player has no node")

type AudioPlayer interface {
	GuildID() snowflake.ID
	ChannelID() *snowflake.ID
	Track() *protocol.Track
	Paused() bool
	Position() protocol.Duration
	State() protocol.PlayerState
	Volume() int
	Filters() protocol.Filters

	Update(ctx context.Context, update protocol.PlayerUpdate) error
	Destroy(ctx context.Context) error

	Lavalink() Lavalink
	Node() Node

	EmitEvent(event protocol.Event)
	AddListeners(listeners ...EventListener)
	RemoveListeners(listeners ...EventListener)

	OnEvent(event protocol.Event)
	OnPlayerUpdate(playerUpdate protocol.PlayerUpdate)
	OnVoiceServerUpdate(token string, endpoint string)
	OnVoiceStateUpdate(channelID snowflake.ID, sessionID string)
}

type defaultPlayer struct {
	guildID   snowflake.ID
	channelID *snowflake.ID
	track     *protocol.Track
	volume    int
	paused    bool
	state     protocol.PlayerState
	voice     protocol.VoiceState
	filters   protocol.Filters

	node     Node
	lavalink Lavalink

	listeners []EventListener
}

func (p *defaultPlayer) GuildID() snowflake.ID {
	return p.guildID
}

func (p *defaultPlayer) ChannelID() *snowflake.ID {
	return p.channelID
}

func (p *defaultPlayer) Track() *protocol.Track {
	return p.track
}

func (p *defaultPlayer) Destroy(ctx context.Context) error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	return p.node.Rest().DestroyPlayer(ctx, p.node.SessionID(), p.guildID)
}

func (p *defaultPlayer) Paused() bool {
	return p.paused
}

func (p *defaultPlayer) Position() protocol.Duration {
	if p.track == nil {
		return 0
	}
	position := p.state.Position
	if p.paused {
		return position
	}
	position += protocol.Duration(time.Now().UnixMilli() - p.state.Time.UnixMilli())
	if position > p.track.Info.Length {
		position = p.track.Info.Length
	} else if position < 0 {
		position = 0
	}
	return position
}

func (p *defaultPlayer) Volume() int {
	return p.volume
}

func (p *defaultPlayer) Filters() protocol.Filters {
	return p.filters
}

func (p *defaultPlayer) Node() Node {
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
	return p.node
}

func (p *defaultPlayer) EmitEvent(event protocol.Event) {
	defer func() {
		if r := recover(); r != nil {
			p.lavalink.Logger().Errorf("recovered from panic in event listener: %#v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	for _, listener := range p.listeners {
		listener.OnEvent(event)
	}
}

func (p *defaultPlayer) AddListeners(listeners ...EventListener) {
	p.listeners = append(p.listeners, listeners...)
}

func (p *defaultPlayer) RemoveListener(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range p.listeners {
			if l == listener {
				p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
			}
		}
	}
}

func (p *defaultPlayer) OnEvent(event protocol.Event) {
	switch e := event.(type) {
	case protocol.EventTrackEnd:
		if e.Reason != protocol.TrackEndReasonReplaced && e.Reason != protocol.TrackEndReasonStopped {
			p.track = nil
		}

	case protocol.EventTrackException, protocol.EventTrackStuck:
		p.track = nil

	case protocol.EventWebSocketClosed:
		p.voice = protocol.VoiceState{}
	}
	p.EmitEvent(event)
}

func (p *defaultPlayer) OnPlayerUpdate(state protocol.PlayerState) {
	p.state = state
}

func (p *defaultPlayer) OnVoiceServerUpdate(token string, endpoint string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := p.Node().Rest().UpdatePlayer(ctx, p.node.SessionID(), p.guildID, protocol.PlayerUpdate{
		Voice: &protocol.VoiceState{
			Token:     token,
			Endpoint:  endpoint,
			SessionID: p.voice.SessionID,
		},
	}); err != nil {
		p.node.Lavalink().Logger().Error("error while sending voice server update: ", err)
	}
	p.voice.Token = token
	p.voice.Endpoint = endpoint
}

func (p *defaultPlayer) OnVoiceStateUpdate(channelID *snowflake.ID, sessionID string) {
	if channelID == nil {
		p.channelID = nil
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.Destroy(ctx); err != nil {
			p.node.Lavalink().Logger().Error("error while destroying player: ", err)
		}
		return
	}
	p.channelID = channelID
	p.voice.SessionID = sessionID
}
