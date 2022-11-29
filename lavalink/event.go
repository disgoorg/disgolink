package lavalink

type EventListener interface {
	OnEvent(event Event)
}

func NewListenerFunc[E Event](f func(e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E Event] struct {
	f func(e E)
}

func (l *listenerFunc[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		l.f(event)
	}
}
