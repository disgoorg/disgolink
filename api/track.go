package api

type Track struct {
	Track string    `json:"track"`
	Info  TrackInfo `json:"info"`
}

type TrackInfo struct {
	Identifier string `json:"identifier"`
	IsSeekable bool   `json:"isSeekable"`
	Author     string `json:"author"`
	Length     int    `json:"length"`
	IsStream   bool   `json:"isStream"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
}
