package api

func NewPlaylist(result *LoadResult) *Playlist {
	return &Playlist{
		Info:   result.PlaylistInfo,
		Tracks: result.Tracks,
	}
}

type Playlist struct {
	Info   *PlaylistInfo
	Tracks []*Track
}

func (p Playlist) SelectedTrack() *Track {
	if p.Info.SelectedTrack == -1 {
		return nil
	}
	return p.Tracks[p.Info.SelectedTrack]
}

type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
