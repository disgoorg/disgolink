package lavalink

import (
	"fmt"

	"github.com/disgoorg/json"
)

var ErrPluginDataNotFound = func(name string) error {
	return fmt.Errorf("plugin data not found for name %s", name)
}

type Plugins []Plugin

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PluginInfo map[string]json.RawMessage

func (p PluginInfo) Get(name string, v any) error {
	data, ok := p[name]
	if !ok {
		return ErrPluginDataNotFound(name)
	}

	return json.Unmarshal(data, v)
}

func (p *PluginInfo) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*map[string]json.RawMessage)(p))
}

func (p PluginInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]json.RawMessage(p))
}
