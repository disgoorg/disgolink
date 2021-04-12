package api

type Playlist struct {
	info   PlaylistInfo
	tracks []Track
}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
