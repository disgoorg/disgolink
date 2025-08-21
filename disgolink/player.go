package disgolink

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

var ErrPlayerNoNode = errors.New("player has no node")

func newPlayer(logger *slog.Logger, lavalink *Client, node *Node, guildID snowflake.ID) *Player {
	return &Player{
		logger:  logger.With(slog.String("name", "disgolink_player"), slog.Int64("guild_id", int64(guildID))),
		Client:  lavalink,
		Node:    node,
		GuildID: guildID,
		Volume:  100,
	}
}

type Player struct {
	logger *slog.Logger
	Node   *Node
	Client *Client

	GuildID   snowflake.ID
	Track     *lavalink.Track
	Volume    int
	Paused    bool
	State     lavalink.PlayerState
	Voice     lavalink.VoiceState
	Filters   lavalink.Filters
}

func (p *Player) Position() lavalink.Duration {
	if p.Track == nil {
		return 0
	}
	position := p.State.Position
	if p.Paused {
		return position
	}
	position += lavalink.Duration(time.Now().UnixMilli() - p.State.Time.UnixMilli())
	if position > p.Track.Info.Length {
		position = p.Track.Info.Length
	} else if position < 0 {
		position = 0
	}
	return position
}

func (p *Player) Update(ctx context.Context, opts ...PlayerUpdateOpt) error {
	if p.Node == nil {
		return ErrPlayerNoNode
	}

	update := defaultPlayerUpdate()
	playerUpdateApply(&update, opts)

	updatedPlayer, err := p.Node.Rest.UpdatePlayer(ctx, p.Node.SessionID, p.GuildID, update)
	if err != nil {
		return err
	}

	p.Volume = updatedPlayer.Volume

	p.Voice = updatedPlayer.Voice
	p.Filters = updatedPlayer.Filters

	// dispatch artificial player resume/pause event
	if update.Paused != nil {
		var event lavalink.Event
		if p.Paused && !*update.Paused {
			event = lavalink.PlayerResumeEvent{
				GuildID: p.GuildID,
			}
		} else if !p.Paused && *update.Paused {
			event = lavalink.PlayerPauseEvent{
				GuildID: p.GuildID,
			}
		}
		p.Paused = updatedPlayer.Paused
		go p.OnEvent(event)
	}

	return nil
}

func (p *Player) Destroy(ctx context.Context) error {
	if err := p.Node.Rest.DestroyPlayer(ctx, p.Node.SessionID, p.GuildID); err != nil {
		return err
	}

	// check if this player already got destroyed
	if player := p.Client.ExistingPlayer(p.GuildID); player == nil {
		return nil
	}

	for plugin := range p.Client.Plugins() {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnDestroyPlayer(p)
		}
	}

	p.Client.RemovePlayer(p.GuildID)

	return nil
}

func (p *Player) OnEvent(event lavalink.Event) {
	switch e := event.(type) {
	case lavalink.UnknownEvent:
		for plugin := range p.Client.Plugins() {
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
		}
	case lavalink.PlayerPauseEvent:
		p.Paused = true

	case lavalink.PlayerResumeEvent:
		p.Paused = false

	case lavalink.TrackStartEvent:
		p.Track = &e.Track

	case lavalink.TrackEndEvent:
		p.Track = nil

	case lavalink.WebSocketClosedEvent:
		p.Voice = lavalink.VoiceState{}
		p.State.Connected = false
	}
}

func (p *Player) OnPlayerUpdate(state lavalink.PlayerState) {
	p.State = state
}

func (p *Player) OnVoiceServerUpdate(ctx context.Context, token string, endpoint string) {
	p.Voice.Token = token
	p.Voice.Endpoint = endpoint

	// we should receive a session id from a VOICE_STATE_UPDATE event before the VOICE_SERVER_UPDATE event.
	if p.Voice.SessionID == "" {
		return
	}

	if err := p.sendVoiceUpdate(ctx); err != nil {
		p.logger.ErrorContext(ctx, "error while sending voice server update", slog.Any("err", err))
	}
}

func (p *Player) OnVoiceStateUpdate(ctx context.Context, channelID *snowflake.ID, sessionID string) {
	if channelID == nil {
		p.Voice = lavalink.VoiceState{}
		if err := p.Destroy(ctx); err != nil {
			p.logger.ErrorContext(ctx, "error while destroying player", slog.Any("err", err))
		}
		p.Client.RemovePlayer(p.GuildID)
		return
	}
	p.Voice.ChannelID = *channelID
	if sessionID != p.Voice.SessionID {
		p.voice.SessionID = sessionID
		if err := p.sendVoiceUpdate(ctx); err != nil {
			p.logger.ErrorContext(ctx, "error while sending voice update", slog.Any("err", err))
		}
	}
}

func (p *Player) sendVoiceUpdate(ctx context.Context) error {
	if p.Voice.SessionID == "" || p.Voice.Token == "" || p.Voice.Endpoint == "" || p.Voice.ChannelID == 0 {
		return nil
	}

	if _, err := p.Node.Rest.UpdatePlayer(ctx, p.Node.SessionID, p.GuildID, lavalink.PlayerUpdate{
		Voice: &p.Voice,
	}); err != nil {
		return fmt.Errorf("error while sending voice update: %w", err)
	}
	return nil
}

func (p *Player) Move(ctx context.Context, node *Node) error {
	if p.Node == node {
		return nil
	}

	if err := p.Destroy(ctx); err != nil {
		return fmt.Errorf("error while destroying player: %w", err)
	}

	p.Node = node

	opts := []PlayerUpdateOpt{
		WithPosition(p.Position()),
		WithVolume(p.Volume),
		WithPaused(p.Paused),
		WithVoice(p.Voice),
		WithFilters(p.Filters),
	}
	if p.Track != nil {
		opts = append(opts, WithTrack(*p.Track))
	}

	if err := p.Update(ctx, opts...); err != nil {
		return fmt.Errorf("error while updating player: %w", err)
	}

	return nil
}
