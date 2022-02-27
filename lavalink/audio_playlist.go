package lavalink

func NewAudioPlaylist(name string, selectedTrackIndex int, tracks []AudioTrack) AudioPlaylist {
	return BasicAudioPlaylist{
		PlaylistName:       name,
		SelectedTrackIndex: selectedTrackIndex,
		PlaylistTracks:     tracks,
	}
}

type AudioPlaylist interface {
	Name() string
	Tracks() []AudioTrack
	SelectedTrack() AudioTrack
}

type BasicAudioPlaylist struct {
	PlaylistName       string
	SelectedTrackIndex int
	PlaylistTracks     []AudioTrack
}

func (p BasicAudioPlaylist) Name() string {
	return p.PlaylistName
}

func (p BasicAudioPlaylist) Tracks() []AudioTrack {
	return p.PlaylistTracks
}

func (p BasicAudioPlaylist) SelectedTrack() AudioTrack {
	if p.SelectedTrackIndex == -1 || p.SelectedTrackIndex >= len(p.PlaylistTracks) {
		return nil
	}
	return p.PlaylistTracks[p.SelectedTrackIndex]
}
