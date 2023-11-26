package lavalink

import (
	"github.com/disgoorg/json"
)

type Plugins []Plugin

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type RawData json.RawMessage

func (p RawData) String() string {
	return string(p)
}

func (p RawData) Unmarshal(v any) error {
	return json.Unmarshal(p, v)
}

func (p *RawData) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.RawMessage)(p))
}

func (p RawData) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.RawMessage(p))
}
