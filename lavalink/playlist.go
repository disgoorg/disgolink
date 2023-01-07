package lavalink

type Playlist struct {
	Info       PlaylistInfo   `json:"info"`
	PluginInfo map[string]any `json:"pluginInfo"`
	Tracks     []Track        `json:"tracks"`
}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
