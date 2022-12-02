package disgolink

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var ErrPlayerNoNode = errors.New("player has no node")

type AudioPlayer interface {
	GuildID() snowflake.ID
	ChannelID() *snowflake.ID
	Track() *lavalink.Track
	Paused() bool
	Position() lavalink.Duration
	State() lavalink.PlayerState
	Volume() int
	Filters() lavalink.Filters

	Update(ctx context.Context, update lavalink.PlayerUpdate) error
	Destroy(ctx context.Context) error

	Lavalink() Client
	Node() Node

	EmitEvent(event lavalink.Event)
	AddListeners(listeners ...EventListener)
	RemoveListeners(listeners ...EventListener)

	OnEvent(event lavalink.Event)
	OnPlayerUpdate(state lavalink.PlayerState)
	OnVoiceServerUpdate(token string, endpoint string)
	OnVoiceStateUpdate(channelID *snowflake.ID, sessionID string)
}

func NewPlayer(lavalink Client, node Node, guildID snowflake.ID) AudioPlayer {
	return &defaultPlayer{
		lavalink: lavalink,
		node:     node,
		guildID:  guildID,
		volume:   100,
	}
}

type defaultPlayer struct {
	guildID   snowflake.ID
	channelID *snowflake.ID
	track     *lavalink.Track
	volume    int
	paused    bool
	state     lavalink.PlayerState
	voice     lavalink.VoiceState
	filters   lavalink.Filters

	node     Node
	lavalink Client

	listeners []EventListener
}

func (p *defaultPlayer) GuildID() snowflake.ID {
	return p.guildID
}

func (p *defaultPlayer) ChannelID() *snowflake.ID {
	return p.channelID
}

func (p *defaultPlayer) Track() *lavalink.Track {
	return p.track
}

func (p *defaultPlayer) Paused() bool {
	return p.paused
}

func (p *defaultPlayer) Position() lavalink.Duration {
	if p.track == nil {
		return 0
	}
	position := p.state.Position
	if p.paused {
		return position
	}
	position += lavalink.Duration(time.Now().UnixMilli() - p.state.Time.UnixMilli())
	if position > p.track.Info.Length {
		position = p.track.Info.Length
	} else if position < 0 {
		position = 0
	}
	return position
}

func (p *defaultPlayer) State() lavalink.PlayerState {
	return p.state
}

func (p *defaultPlayer) Volume() int {
	return p.volume
}

func (p *defaultPlayer) Filters() lavalink.Filters {
	return p.filters
}

func (p *defaultPlayer) Update(ctx context.Context, update lavalink.PlayerUpdate) error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	_, err := p.node.Rest().UpdatePlayer(ctx, p.node.SessionID(), p.guildID, update)
	return err
}

func (p *defaultPlayer) Destroy(ctx context.Context) error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	return p.node.Rest().DestroyPlayer(ctx, p.node.SessionID(), p.guildID)
}

func (p *defaultPlayer) Node() Node {
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
	return p.node
}

func (p *defaultPlayer) Lavalink() Client {
	return p.lavalink
}

func (p *defaultPlayer) EmitEvent(event lavalink.Event) {
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

func (p *defaultPlayer) RemoveListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range p.listeners {
			if l == listener {
				p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
			}
		}
	}
}

func (p *defaultPlayer) OnEvent(event lavalink.Event) {
	switch e := event.(type) {
	case lavalink.EventTrackEnd:
		if e.Reason != lavalink.TrackEndReasonReplaced && e.Reason != lavalink.TrackEndReasonStopped {
			p.track = nil
		}

	case lavalink.EventTrackException, lavalink.EventTrackStuck:
		p.track = nil

	case lavalink.EventWebSocketClosed:
		p.voice = lavalink.VoiceState{}
	}
	p.EmitEvent(event)
}

func (p *defaultPlayer) OnPlayerUpdate(state lavalink.PlayerState) {
	p.state = state
}

func (p *defaultPlayer) OnVoiceServerUpdate(token string, endpoint string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := p.Node().Rest().UpdatePlayer(ctx, p.node.SessionID(), p.guildID, lavalink.PlayerUpdate{
		Voice: &lavalink.VoiceState{
			Token:     token,
			Endpoint:  endpoint,
			SessionID: p.voice.SessionID,
		},
	}); err != nil {
		p.lavalink.Logger().Error("error while sending voice server update: ", err)
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
