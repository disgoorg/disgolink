package disgolink

import "github.com/disgoorg/disgolink/v4/lavalink"

type TrackLoadingResultHandler interface {
	OnTrack(track lavalink.Track)
	OnPlaylist(playlist lavalink.Playlist)
	OnSearch(tracks []lavalink.Track)
	OnEmpty()
	OnError(err error)
}

var _ TrackLoadingResultHandler = (*functionalTrackLoadingResultHandler)(nil)

func NewTrackLoadingResultHandler(trackLoaded func(track lavalink.Track), playlistLoaded func(playlist lavalink.Playlist), searchResultLoaded func(tracks []lavalink.Track), noMatches func(), loadFailed func(err error)) TrackLoadingResultHandler {
	return functionalTrackLoadingResultHandler{
		onTrack:    trackLoaded,
		onPlaylist: playlistLoaded,
		onSearch:   searchResultLoaded,
		onEmpty:    noMatches,
		onError:    loadFailed,
	}
}

type functionalTrackLoadingResultHandler struct {
	onTrack    func(track lavalink.Track)
	onPlaylist func(playlist lavalink.Playlist)
	onSearch   func(tracks []lavalink.Track)
	onEmpty    func()
	onError    func(err error)
}

func (h functionalTrackLoadingResultHandler) OnTrack(track lavalink.Track) {
	if h.onTrack != nil {
		h.onTrack(track)
	}
}
func (h functionalTrackLoadingResultHandler) OnPlaylist(playlist lavalink.Playlist) {
	if h.onPlaylist != nil {
		h.onPlaylist(playlist)
	}
}
func (h functionalTrackLoadingResultHandler) OnSearch(tracks []lavalink.Track) {
	if h.onSearch != nil {
		h.onSearch(tracks)
	}
}
func (h functionalTrackLoadingResultHandler) OnEmpty() {
	if h.onEmpty != nil {
		h.onEmpty()
	}
}
func (h functionalTrackLoadingResultHandler) OnError(err error) {
	if h.onError != nil {
		h.onError(err)
	}
}
