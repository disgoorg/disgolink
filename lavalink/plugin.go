package lavalink

import "io"

type Plugin interface {
	Name() string
	Version() string
}

type WebSocketExtension interface {
	Op() OpType
	OnInvocation(node Node, data []byte)
}

type EventExtension interface {
	Event() EventType
	OnInvocation(node Node, data []byte)
}

type SourceExtension interface {
	SourceName() string
	Encode(track Track, w io.Writer) error
	Decode(trackInfo TrackInfo, r io.Reader) (Track, error)
}

type PluginEventHandler struct {
	OnWebSocketOpen       func(node Node)
	OnWebSocketDestroyed  func(node Node)
	OnWebSocketMessageIn  func(node Node, data []byte)
	OnWebSocketMessageOut func(node Node, data []byte)
	OnNewPlayer           func(node Node, player Player)
	OnDestroyPlayer       func(node Node, player Player)
}
