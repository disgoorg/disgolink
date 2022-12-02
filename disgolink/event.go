package disgolink

import "github.com/disgoorg/disgolink/v2/lavalink"

type EventListener interface {
	OnEvent(event lavalink.Event)
}

func NewListenerFunc[E lavalink.Event](f func(e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E lavalink.Event] struct {
	f func(e E)
}

func (l *listenerFunc[E]) OnEvent(e lavalink.Event) {
	if event, ok := e.(E); ok {
		l.f(event)
	}
}
