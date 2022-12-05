package disgolink

import "github.com/disgoorg/disgolink/v2/lavalink"

type EventListener interface {
	OnEvent(player Player, event lavalink.Event)
}

func NewListenerFunc[E lavalink.Event](f func(p Player, e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E lavalink.Event] struct {
	f func(p Player, e E)
}

func (l *listenerFunc[E]) OnEvent(p Player, e lavalink.Event) {
	if event, ok := e.(E); ok {
		l.f(p, event)
	}
}
