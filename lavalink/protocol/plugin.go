package protocol

type Plugins []Plugin

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
