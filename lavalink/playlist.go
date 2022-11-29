package lavalink

type Playlist struct {
	Info       PlaylistInfo   `json:"info"`
	PluginData map[string]any `json:"pluginData"`
	Tracks     []Track        `json:"tracks"`
}

type PlaylistInfo struct {
	Name          string
	SelectedTrack int
}
