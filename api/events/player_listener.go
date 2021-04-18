package events

type PlayerEventListener interface {
	OnEvent(event PlayerEvent)
}
