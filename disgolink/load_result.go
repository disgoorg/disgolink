package disgolink

import "github.com/disgoorg/disgolink/v3/lavalink"

type AudioLoadResultHandler interface {
	TrackLoaded(track lavalink.Track)
	PlaylistLoaded(playlist lavalink.Playlist)
	SearchResultLoaded(tracks []lavalink.Track)
	NoMatches()
	LoadFailed(err error)
}

var _ AudioLoadResultHandler = (*FunctionalResultHandler)(nil)

func NewResultHandler(trackLoaded func(track lavalink.Track), playlistLoaded func(playlist lavalink.Playlist), searchResultLoaded func(tracks []lavalink.Track), noMatches func(), loadFailed func(err error)) AudioLoadResultHandler {
	return FunctionalResultHandler{
		trackLoaded:        trackLoaded,
		playlistLoaded:     playlistLoaded,
		searchResultLoaded: searchResultLoaded,
		noMatches:          noMatches,
		loadFailed:         loadFailed,
	}
}

type FunctionalResultHandler struct {
	trackLoaded        func(track lavalink.Track)
	playlistLoaded     func(playlist lavalink.Playlist)
	searchResultLoaded func(tracks []lavalink.Track)
	noMatches          func()
	loadFailed         func(err error)
}

func (h FunctionalResultHandler) TrackLoaded(track lavalink.Track) {
	if h.trackLoaded != nil {
		h.trackLoaded(track)
	}
}
func (h FunctionalResultHandler) PlaylistLoaded(playlist lavalink.Playlist) {
	if h.playlistLoaded != nil {
		h.playlistLoaded(playlist)
	}
}
func (h FunctionalResultHandler) SearchResultLoaded(tracks []lavalink.Track) {
	if h.searchResultLoaded != nil {
		h.searchResultLoaded(tracks)
	}
}
func (h FunctionalResultHandler) NoMatches() {
	if h.noMatches != nil {
		h.noMatches()
	}
}
func (h FunctionalResultHandler) LoadFailed(err error) {
	if h.loadFailed != nil {
		h.loadFailed(err)
	}
}
