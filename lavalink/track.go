package lavalink

type Track struct {
	Encoded string    `json:"encoded"`
	Info    TrackInfo `json:"info"`
}

type TrackInfo struct {
	Identifier string   `json:"identifier"`
	Author     string   `json:"author"`
	Length     Duration `json:"length"`
	IsStream   bool     `json:"isStream"`
	Title      string   `json:"title"`
	URI        *string  `json:"uri"`
	SourceName string   `json:"sourceName"`
	Position   Duration `json:"position"`
	ArtworkURL *string  `json:"artwork"`
	ISRC       *string  `json:"isrc"`
}
