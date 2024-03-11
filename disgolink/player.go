package disgolink

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var ErrPlayerNoNode = errors.New("player has no node")

type Player interface {
	GuildID() snowflake.ID
	ChannelID() *snowflake.ID
	Track() *lavalink.Track
	Paused() bool
	Position() lavalink.Duration
	State() lavalink.PlayerState
	Volume() int
	Filters() lavalink.Filters

	Update(ctx context.Context, opts ...lavalink.PlayerUpdateOpt) error
	Destroy(ctx context.Context) error

	Lavalink() Client
	Node() Node

	Restore(player lavalink.Player)
	OnEvent(event lavalink.Event)
	OnPlayerUpdate(state lavalink.PlayerState)
	OnVoiceServerUpdate(ctx context.Context, token string, endpoint string)
	OnVoiceStateUpdate(ctx context.Context, channelID *snowflake.ID, sessionID string)
}

func NewPlayer(logger *slog.Logger, lavalink Client, node Node, guildID snowflake.ID) Player {
	return &playerImpl{
		logger:   logger,
		lavalink: lavalink,
		node:     node,
		guildID:  guildID,
		volume:   100,
	}
}

type playerImpl struct {
	logger   *slog.Logger
	node     Node
	lavalink Client

	guildID   snowflake.ID
	channelID *snowflake.ID
	track     *lavalink.Track
	volume    int
	paused    bool
	state     lavalink.PlayerState
	voice     lavalink.VoiceState
	filters   lavalink.Filters
}

func (p *playerImpl) GuildID() snowflake.ID {
	return p.guildID
}

func (p *playerImpl) ChannelID() *snowflake.ID {
	return p.channelID
}

func (p *playerImpl) Track() *lavalink.Track {
	return p.track
}

func (p *playerImpl) Paused() bool {
	return p.paused
}

func (p *playerImpl) Position() lavalink.Duration {
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

func (p *playerImpl) State() lavalink.PlayerState {
	return p.state
}

func (p *playerImpl) Volume() int {
	return p.volume
}

func (p *playerImpl) Filters() lavalink.Filters {
	return p.filters
}

func (p *playerImpl) Update(ctx context.Context, opts ...lavalink.PlayerUpdateOpt) error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	update := lavalink.DefaultPlayerUpdate()
	update.Apply(opts)

	updatedPlayer, err := p.node.Rest().UpdatePlayer(ctx, p.node.SessionID(), p.guildID, *update)
	if err != nil {
		return err
	}

	p.track = updatedPlayer.Track
	if updatedPlayer.Track != nil {
		p.state.Position = updatedPlayer.Track.Info.Position
	} else {
		p.state.Position = 0
	}
	p.state.Time = lavalink.Now()
	p.volume = updatedPlayer.Volume

	p.voice = updatedPlayer.Voice
	p.filters = updatedPlayer.Filters

	// dispatch artificial player resume/pause event
	if update.Paused != nil {
		var event lavalink.Event
		if p.paused && !*update.Paused {
			event = lavalink.PlayerResumeEvent{
				GuildID_: p.guildID,
			}
		} else if !p.paused && *update.Paused {
			event = lavalink.PlayerPauseEvent{
				GuildID_: p.guildID,
			}
		}
		p.paused = updatedPlayer.Paused
		go p.OnEvent(event)
	}

	return nil
}

func (p *playerImpl) Destroy(ctx context.Context) error {
	if p.node == nil {
		return ErrPlayerNoNode
	}

	err := p.node.Rest().DestroyPlayer(ctx, p.node.SessionID(), p.guildID)
	if err != nil {
		return err
	}

	p.lavalink.ForPlugins(func(plugin Plugin) {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnDestroyPlayer(p)
		}
	})

	p.lavalink.RemovePlayer(p.guildID)

	return nil
}

func (p *playerImpl) Node() Node {
	if p.node == nil {
		p.node = p.lavalink.BestNode()
	}
	return p.node
}

func (p *playerImpl) Lavalink() Client {
	return p.lavalink
}

func (p *playerImpl) Restore(player lavalink.Player) {
	p.track = player.Track
	p.state = player.State
	p.paused = player.Paused
	p.voice = player.Voice
	p.filters = player.Filters
	p.volume = player.Volume
}

func (p *playerImpl) OnEvent(event lavalink.Event) {
	switch e := event.(type) {
	case lavalink.UnknownEvent:
		p.lavalink.ForPlugins(func(plugin Plugin) {
			if pl, ok := plugin.(EventPlugin); ok && pl.Event() == e.Type() {
				pl.OnEventInvocation(p, e.Data)
			}
			if pl, ok := plugin.(EventPlugins); ok {
				for _, pls := range pl.EventPlugins() {
					if pls.Event() == e.Type() {
						pls.OnEventInvocation(p, e.Data)
					}
				}
			}
		})
	case lavalink.PlayerPauseEvent:
		p.paused = true

	case lavalink.PlayerResumeEvent:
		p.paused = false

	case lavalink.TrackEndEvent:
		if p.track != nil {
			e.Track = *p.track
		}
		if e.Reason != lavalink.TrackEndReasonReplaced && e.Reason != lavalink.TrackEndReasonStopped {
			p.track = nil
		}

	case lavalink.TrackExceptionEvent:
		if p.track != nil {
			e.Track = *p.track
		}
		p.track = nil

	case lavalink.TrackStuckEvent:
		if p.track != nil {
			e.Track = *p.track
		}
		p.track = nil

	case lavalink.WebSocketClosedEvent:
		p.voice = lavalink.VoiceState{}
		p.state.Connected = false
	}
}

func (p *playerImpl) OnPlayerUpdate(state lavalink.PlayerState) {
	p.state = state
}

func (p *playerImpl) OnVoiceServerUpdate(ctx context.Context, token string, endpoint string) {
	if _, err := p.Node().Rest().UpdatePlayer(ctx, p.node.SessionID(), p.guildID, lavalink.PlayerUpdate{
		Voice: &lavalink.VoiceState{
			Token:     token,
			Endpoint:  endpoint,
			SessionID: p.voice.SessionID,
		},
	}); err != nil {
		p.logger.ErrorContext(ctx, "error while sending voice server update", slog.Any("err", err))
	}
	p.voice.Token = token
	p.voice.Endpoint = endpoint
}

func (p *playerImpl) OnVoiceStateUpdate(ctx context.Context, channelID *snowflake.ID, sessionID string) {
	if channelID == nil {
		p.channelID = nil
		if err := p.Destroy(ctx); err != nil {
			p.logger.ErrorContext(ctx, "error while destroying player", slog.Any("err", err))
		}
		p.lavalink.RemovePlayer(p.guildID)
		return
	}
	p.channelID = channelID
	p.voice.SessionID = sessionID
}
