package disgolink

import (
	"github.com/disgoorg/disgolink/v4/lavalink"
)

func newGenericEvent(node *Node) *GenericEvent {
	return &GenericEvent{
		node: node,
	}
}

type GenericEvent struct {
	node *Node
}

func (e *GenericEvent) Client() *Client {
	return e.node.Client
}

func (e *GenericEvent) Node() *Node {
	return e.node
}

type ReadyEvent struct {
	*GenericEvent
	lavalink.ReadyMessage
}

type StatsEvent struct {
	*GenericEvent
	lavalink.StatsMessage
}

type PlayerUpdateEvent struct {
	*GenericEvent
	lavalink.PlayerUpdateMessage
	Player *Player
}

type PlayerStartTrackEvent struct {
	*GenericEvent
	lavalink.TrackStartEvent
	Player *Player
}

type PlayerTrackEndEvent struct {
	*GenericEvent
	lavalink.TrackEndEvent
	Player *Player
}

type PlayerTrackExceptionEvent struct {
	*GenericEvent
	lavalink.TrackExceptionEvent
	Player *Player
}

type PlayerTrackStuckEvent struct {
	*GenericEvent
	lavalink.TrackStuckEvent
	Player *Player
}

type WebSocketClosedEvent struct {
	*GenericEvent
	lavalink.WebSocketClosedEvent
	Player *Player
}

type UnknownEvent struct {
	*GenericEvent
	lavalink.UnknownEvent
	Player *Player
}

type UnknownMessageEvent struct {
	*GenericEvent
	lavalink.UnknownMessage
}
