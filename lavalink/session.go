package lavalink

type Session struct {
	Resuming bool `json:"resuming"`
	Timeout  int  `json:"timeout"`
}

type SessionUpdate struct {
	Resuming *bool `json:"resuming,omitzero"`
	Timeout  *int  `json:"timeout,omitzero"`
}
