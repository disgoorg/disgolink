package disgolink

import "github.com/disgoorg/disgolink/v2/lavalink"

type Plugin interface {
	Name() string
	Version() string
}

type OpPlugin interface {
	Op() lavalink.Op
	OnOpInvocation(node Node, data []byte)
}

type EventPlugin interface {
	Event() lavalink.EventType
	OnEventInvocation(player Player, data []byte)
}

type PluginEventHandler interface {
	OnNodeOpen(node Node)
	OnNodeClose(node Node)
	OnNodeMessageIn(node Node, data []byte)
	OnNewPlayer(player Player)
	OnDestroyPlayer(player Player)
}
