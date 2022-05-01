package lavalink

import "io"

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OpPlugins interface {
	OpPlugins() []OpPlugin
}

type OpPlugin interface {
	Op() OpType
	OnOpInvocation(node Node, data []byte)
}

type EventPlugins interface {
	EventPlugins() []EventPlugin
}

type EventPlugin interface {
	Event() EventType
	OnEventInvocation(node Node, data []byte)
}

type SourcePlugin interface {
	SourceName() string
	Encode(track AudioTrack, w io.Writer) error
	Decode(info AudioTrackInfo, r io.Reader) (AudioTrack, error)
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
