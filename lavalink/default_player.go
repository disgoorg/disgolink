package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/filters"
	"github.com/pkg/errors"
)

var _ Player = (*DefaultPlayer)(nil)

func NewPlayer(node Node, guildID discord.Snowflake) Player {
	return &DefaultPlayer{
		guildID:       guildID,
		channelID:     nil,
		lastSessionID: nil,
		track:         nil,
		volume:        100,
		paused:        false,
		position:      -1,
		connected:     false,
		updateTime:    -1,
		filters:       nil,
		node:          node,
		listeners:     nil,
	}
}

type DefaultPlayer struct {
	guildID       discord.Snowflake
	channelID     *discord.Snowflake
	lastSessionID *string
	track         Track
	volume        int
	paused        bool
	position      int
	connected     bool
	updateTime    int
	filters       *filters.Filters
	node          Node
	listeners     []PlayerEventListener
}

func (p *DefaultPlayer) Track() Track {
	return p.track
}

func (p *DefaultPlayer) SetTrack(track Track) {
	p.track = track
}

func (p *DefaultPlayer) Play(track Track) error {
	if track == nil {
		return p.Stop()
	}

	t := track.Track()
	if t == "" {
		return errors.New("can't play empty track")
	}

	if err := p.Node().Send(PlayPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpPlay, p),
		Track:         t,
	}); err != nil {
		return errors.Wrap(err, "error while playing track")
	}
	return nil
}

func (p *DefaultPlayer) PlayAt(track Track, start int, end int) error {
	if track == nil {
		return p.Stop()
	}

	t := track.Track()
	if t == "" {
		return errors.New("can't play empty track")
	}

	if err := p.Node().Send(PlayPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpPlay, p),
		Track:         t,
		StartTime:     &start,
		EndTime:       &end,
	}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	return nil
}

func (p *DefaultPlayer) Stop() error {
	p.track = nil

	if err := p.Node().Send(StopPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpStop, p),
	}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
	return nil
}

func (p *DefaultPlayer) Destroy() error {
	p.track = nil

	if err := p.Node().Send(DestroyPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpDestroy, p),
	}); err != nil {
		return errors.Wrap(err, "error while destroying player")
	}
	return nil
}

func (p *DefaultPlayer) Pause(pause bool) error {
	if p.paused == pause {
		return nil
	}

	if err := p.Node().Send(&PausePlayerCommand{
		PlayerCommand: NewPlayerCommand(OpPause, p),
		Pause:         pause,
	}); err != nil {
		return errors.Wrap(err, "error while pausing player")
	}

	p.paused = pause
	if pause {
		p.EmitEvent(func(listener PlayerEventListener) {
			listener.OnPlayerPause(p)
		})
	} else {
		p.EmitEvent(func(listener PlayerEventListener) {
			listener.OnPlayerResume(p)
		})
	}
	return nil
}

func (p *DefaultPlayer) Paused() bool {
	return p.paused
}

func (p *DefaultPlayer) Position() int {
	return p.position
}

func (p *DefaultPlayer) Seek(position int) error {
	if err := p.Node().Send(&SeekPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpSeek, p),
		Position:      position,
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
	if err := p.Node().Send(&VolumePlayerCommand{
		PlayerCommand: NewPlayerCommand(OpSeek, p),
		Volume:        volume,
	}); err != nil {
		return errors.Wrap(err, "error while setting volume of player")
	}
	p.volume = volume
	return nil
}

func (p *DefaultPlayer) Filters() *filters.Filters {
	if p.filters == nil {
		p.filters = filters.NewFilters(p.commitFilters)
	}
	return p.filters
}

func (p *DefaultPlayer) SetFilters(filters *filters.Filters) {
	p.filters = filters
}

func (p *DefaultPlayer) GuildID() discord.Snowflake {
	return p.guildID
}

func (p *DefaultPlayer) ChannelID() *discord.Snowflake {
	return p.channelID
}

func (p *DefaultPlayer) SetChannelID(channelID *discord.Snowflake) {
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

func (p *DefaultPlayer) ChangeNode(node Node) {
	p.node = node
}

func (p *DefaultPlayer) PlayerUpdate(state State) {
	p.position = state.Position
	p.connected = state.Connected
}

func (p *DefaultPlayer) EmitEvent(listenerCaller func(listener PlayerEventListener)) {
	for _, listener := range p.listeners {
		listenerCaller(listener)
	}
}

func (p *DefaultPlayer) AddListener(playerListener PlayerEventListener) {
	p.listeners = append(p.listeners, playerListener)
}

func (p *DefaultPlayer) RemoveListener(playerListener PlayerEventListener) {
	for i, listener := range p.listeners {
		if listener == playerListener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
		}
	}
}

func (p *DefaultPlayer) commitFilters(filters *filters.Filters) error {
	if err := p.node.Send(&FilterPlayerCommand{
		PlayerCommand: NewPlayerCommand(OpFilters, p),
		Filters:       filters,
	}); err != nil {
		return errors.Wrap(err, "error while setting filters of player")
	}
	return nil
}
