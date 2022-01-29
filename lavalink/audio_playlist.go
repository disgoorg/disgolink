package lavalink

func NewAudioPlaylist(info AudioPlaylistInfo, tracks []AudioTrack) AudioPlaylist {
	return AudioPlaylist{
		Info:   info,
		Tracks: tracks,
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
