package lavalink

type NodeStatus int

// Indicates how far along the client is to connecting
const (
	Connecting NodeStatus = iota
	Reconnecting
	Disconnected
)
