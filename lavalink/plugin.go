package lavalink

import "io"

type Plugin interface {
	Name() string
	Version() string
}

type OpExtension interface {
	Op() OpType
	OnOpInvocation(node Node, data []byte)
}

type EventExtension interface {
	Event() EventType
	OnEventInvocation(node Node, data []byte)
}

type SourceExtension interface {
	SourceName() string
	Encode(track Track, w io.Writer) error
	Decode(trackInfo TrackInfo, r io.Reader) (Track, error)
}

type PluginEventHandler interface {
	OnWebSocketOpen(node Node)
	OnWebSocketDestroy(node Node)
	OnWebSocketMessageIn(node Node, data []byte)
	OnWebSocketMessageOut(node Node, data []byte)
	OnNewPlayer(player Player)
	OnDestroyPlayer(player Player)
}

type PluginEventAdapter struct {
	OnWebSocketOpenEvent       func(node Node)
	OnWebSocketDestroyEvent    func(node Node)
	OnWebSocketMessageInEvent  func(node Node, data []byte)
	OnWebSocketMessageOutEvent func(node Node, data []byte)
	OnNewPlayerEvent           func(player Player)
	OnDestroyPlayerEvent       func(player Player)
}

func (a *PluginEventAdapter) OnWebSocketOpen(node Node) {
	if a.OnWebSocketOpenEvent == nil {
		return
	}
	a.OnWebSocketOpenEvent(node)
}
func (a *PluginEventAdapter) OnWebSocketDestroy(node Node) {
	if a.OnWebSocketDestroyEvent == nil {
		return
	}
	a.OnWebSocketDestroyEvent(node)
}
func (a *PluginEventAdapter) OnWebSocketMessageIn(node Node, data []byte) {
	if a.OnWebSocketMessageInEvent == nil {
		return
	}
	a.OnWebSocketMessageInEvent(node, data)
}
func (a *PluginEventAdapter) OnWebSocketMessageOut(node Node, data []byte) {
	if a.OnWebSocketMessageOutEvent == nil {
		return
	}
	a.OnWebSocketMessageOutEvent(node, data)
}
func (a *PluginEventAdapter) OnNewPlayer(player Player) {
	if a.OnNewPlayerEvent == nil {
		return
	}
	a.OnNewPlayerEvent(player)
}
func (a *PluginEventAdapter) OnDestroyPlayer(player Player) {
	if a.OnDestroyPlayerEvent == nil {
		return
	}
	a.OnDestroyPlayerEvent(player)
}
