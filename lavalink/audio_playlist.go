package lavalink

func NewAudioPlaylist(result LoadResult) AudioPlaylist {
	return AudioPlaylist{
		Info:   *result.PlaylistInfo,
		Tracks: result.Tracks,
	}
}

type AudioPlaylist struct {
	Info   AudioPlaylistInfo
	Tracks []AudioTrack
}

func (p AudioPlaylist) SelectedTrack() AudioTrack {
	if p.Info.SelectedTrack == -1 {
		return nil
	}
	if p.Info.SelectedTrack >= len(p.Tracks) {
		return nil
	}
	return p.Tracks[p.Info.SelectedTrack]
}

type AudioPlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}
