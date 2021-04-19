package api

type PlayerEventListener interface {
	OnEvent(event PlayerEvent)
}
