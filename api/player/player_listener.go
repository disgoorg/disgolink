package player

type Listener interface {
	OnEvent(event Event)
}
