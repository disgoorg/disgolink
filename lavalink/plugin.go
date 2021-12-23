package lavalink

type WebsocketMessageInHandler interface {
	OnWebsocketMessageIn(node Node, data []byte) bool
}
