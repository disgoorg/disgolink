package disgolink

import "github.com/disgoorg/disgolink/v3/lavalink"

type EventListener interface {
	OnEvent(player Player, event lavalink.Message)
}

func NewListenerFunc[E lavalink.Message](f func(p Player, e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E lavalink.Message] struct {
	f func(p Player, e E)
}

func (l *listenerFunc[E]) OnEvent(p Player, e lavalink.Message) {
	if event, ok := e.(E); ok {
		l.f(p, event)
	}
}
