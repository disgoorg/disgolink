package protocol

type Playlist struct {
	Info       PlaylistInfo   `json:"info"`
	PluginData map[string]any `json:"pluginData"`
}

type PlaylistInfo struct {
	Name          string
	SelectedTrack int
}
