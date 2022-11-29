package lavalink

import "github.com/disgoorg/disgolink/lavalink/protocol"

type EventListener interface {
	OnEvent(event protocol.Event)
}

func NewListenerFunc[E protocol.Event](f func(e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E protocol.Event] struct {
	f func(e E)
}

func (l *listenerFunc[E]) OnEvent(e protocol.Event) {
	if event, ok := e.(E); ok {
		l.f(event)
	}
}
