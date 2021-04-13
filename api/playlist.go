package api

func NewPlaylist(result *LoadResult, searchResult bool) *Playlist {
	return &Playlist{
		Info:         result.PlaylistInfo,
		Tracks:       result.Tracks,
		SearchResult: searchResult,
	}
}

type Playlist struct {
	Info         *PlaylistInfo
	Tracks       []*Track
	SearchResult bool
}

func (p Playlist) SelectedTrack() *Track {
	return p.Tracks[p.Info.SelectedTrack]
}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
