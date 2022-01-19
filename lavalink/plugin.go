package lavalink

import "io"

type OpExtensions interface {
	OpExtensions() []OpExtension
}

type OpExtension interface {
	Op() OpType
	OnOpInvocation(node Node, data []byte)
}

type EventExtensions interface {
	EventExtensions() []EventExtension
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
	OnNodeOpen(node Node)
	OnNodeDestroy(node Node)
	OnNodeMessageIn(node Node, data []byte)
	OnNodeMessageOut(node Node, data []byte)
	OnNewPlayer(player Player)
	OnDestroyPlayer(player Player)
}

type PluginEventAdapter struct{}

func (a PluginEventAdapter) OnNodeOpen(node Node)                    {}
func (a PluginEventAdapter) OnNodeDestroy(node Node)                 {}
func (a PluginEventAdapter) OnNodeMessageIn(node Node, data []byte)  {}
func (a PluginEventAdapter) OnNodeMessageOut(node Node, data []byte) {}
func (a PluginEventAdapter) OnNewPlayer(player Player)               {}
func (a PluginEventAdapter) OnDestroyPlayer(player Player)           {}