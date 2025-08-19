package disgolink

import (
	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

type Plugin interface {
	Name() string
	Version() string
}

type OpPlugin interface {
	Op() lavalink.Op
	OnOpInvocation(node *Node, data json.RawMessage)
}

type EventPlugin interface {
	Event() lavalink.EventType
	OnEventInvocation(player *Player, data json.RawMessage)
}

type EventPlugins interface {
	EventPlugins() []EventPlugin
}

type PluginEventHandler interface {
	OnNodeOpen(node *Node)
	OnNodeClose(node *Node)
	OnNodeMessageIn(node *Node, data json.RawMessage)
	OnNewPlayer(player *Player)
	OnDestroyPlayer(player *Player)
}
