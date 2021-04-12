package api

type Track struct {
	info  TrackInfo
	track string
}

func (t *Track) Encode() string {
	return ""
}

func (t *Track) Info() TrackInfo {
	return t.info
}

func (t *Track) Position() int {
	return t.info.Position
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
