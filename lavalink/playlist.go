package lavalink

type Playlist struct {
	Info   PlaylistInfo `json:"info"`
	Tracks []Track      `json:"tracks"`
}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
