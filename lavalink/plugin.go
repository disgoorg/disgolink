package lavalink

import (
	"github.com/disgoorg/json"
)

type Plugins []Plugin

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PluginInfo json.RawMessage

func (p PluginInfo) Unmarshal(v any) error {
	return json.Unmarshal(p, v)
}

func (p *PluginInfo) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.RawMessage)(p))
}

func (p PluginInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.RawMessage(p))
}
