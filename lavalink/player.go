package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/filters"
	"github.com/pkg/errors"
)

type Player interface {
	Track() Track
	SetTrack(track Track)
	Play(track Track) error
	PlayAt(track Track, start int, end int) error
	Stop() error
	Destroy() error
	Pause(paused bool) error
	Paused() bool
	Position() int
	Seek(position int) error
	Volume() int
	SetVolume(volume int) error
	Filters() *filters.Filters
	SetFilters(filters *filters.Filters)

	GuildID() discord.Snowflake
	ChannelID() *discord.Snowflake
	SetChannelID(channelID *discord.Snowflake)
	LastSessionID() *string
	SetLastSessionID(sessionID string)

	Node() Node
	ChangeNode(node Node)

	PlayerUpdate(state PlayerState)
	EmitEvent(listenerCaller func(listener PlayerEventListener))
	AddListener(playerListener PlayerEventListener)
	RemoveListener(playerListener PlayerEventListener)
}

type PlayerState struct {
	Time      int64 `json:"time"`
	Position  int   `json:"position"`
	Connected bool  `json:"connected"`
}

var _ Player = (*DefaultPlayer)(nil)

func NewPlayer(node Node, guildID discord.Snowflake) Player {
	return &DefaultPlayer{
		guildID: guildID,
		volume:  100,
		node:    node,
	}
}

type DefaultPlayer struct {
	guildID       discord.Snowflake
	channelID     *discord.Snowflake
	lastSessionID *string
	track         Track
	volume        int
	paused        bool
	state         PlayerState
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
	t := track.Track()
	if t == "" {
		return errors.New("can't play empty track")
	}

	if err := p.Node().Send(PlayCommand{
		GuildID: p.guildID,
		Track:   t,
	}); err != nil {
		return errors.Wrap(err, "error while playing track")
	}
	return nil
}

func (p *DefaultPlayer) PlayAt(track Track, start int, end int) error {
	t := track.Track()
	if t == "" {
		return errors.New("can't play empty track")
	}

	if err := p.Node().Send(PlayCommand{
		GuildID:   p.guildID,
		Track:     t,
		StartTime: start,
		EndTime:   end,
	}); err != nil {
		return errors.Wrap(err, "error while stopping player")
	}
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
	return p.state.Position
}

func (p *DefaultPlayer) Seek(position int) error {
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

func (p *DefaultPlayer) PlayerUpdate(state PlayerState) {
	p.state = state
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
	if err := p.node.Send(FiltersCommand{
		GuildID: p.guildID,
		Filters: *filters,
	}); err != nil {
		return errors.Wrap(err, "error while setting filters of player")
	}
	return nil
}
