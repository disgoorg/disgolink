package protocol

import "github.com/disgoorg/json"

type Session struct {
	ResumingKey *string `json:"resumingKey"`
	Timeout     int     `json:"timeout"`
}

type SessionUpdate struct {
	ResumingKey *json.Nullable[string] `json:"resumingKey"`
	Timeout     *int                   `json:"timeout"`
}
