package lavalink

type Playlist struct {
	Info       PlaylistInfo `json:"info"`
	PluginInfo PluginInfo   `json:"pluginInfo"`
	Tracks     []Track      `json:"tracks"`
}

func (Playlist) loadResultData() {}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
