package lavalink

import "encoding/json"

type LoadType string

const (
	LoadTypeTrackLoaded    LoadType = "TRACK_LOADED"
	LoadTypePlaylistLoaded LoadType = "PLAYLIST_LOADED"
	LoadTypeSearchResult   LoadType = "SEARCH_RESULT"
	LoadTypeNoMatches      LoadType = "NO_MATCHES"
	LoadTypeLoadFailed     LoadType = "LOAD_FAILED"
)

type AudioLoadResultHandler interface {
	TrackLoaded(track AudioTrack)
	PlaylistLoaded(playlist AudioPlaylist)
	SearchResultLoaded(tracks []AudioTrack)
	NoMatches()
	LoadFailed(e FriendlyException)
}

type LoadResult struct {
	LoadType     LoadType           `json:"loadType"`
	PlaylistInfo *AudioPlaylistInfo `json:"playlistInfo"`
	Tracks       []AudioTrack       `json:"tracks"`
	Exception    *FriendlyException `json:"exception"`
}

func (r *LoadResult) UnmarshalJSON(data []byte) error {
	type loadResult LoadResult
	var v struct {
		Tracks []*DefaultAudioTrack `json:"tracks"`
		loadResult
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = LoadResult(v.loadResult)
	if v.Tracks != nil {
		r.Tracks = make([]AudioTrack, len(v.Tracks))
		for i, t := range v.Tracks {
			r.Tracks[i] = t
		}
	}
	return nil
}

var _ AudioLoadResultHandler = (*FunctionalResultHandler)(nil)

func NewResultHandler(trackLoaded func(track AudioTrack), playlistLoaded func(playlist AudioPlaylist), searchResultLoaded func(tracks []AudioTrack), noMatches func(), loadFailed func(e FriendlyException)) AudioLoadResultHandler {
	return FunctionalResultHandler{
		trackLoaded:        trackLoaded,
		playlistLoaded:     playlistLoaded,
		searchResultLoaded: searchResultLoaded,
		noMatches:          noMatches,
		loadFailed:         loadFailed,
	}
}

type FunctionalResultHandler struct {
	trackLoaded        func(track AudioTrack)
	playlistLoaded     func(playlist AudioPlaylist)
	searchResultLoaded func(tracks []AudioTrack)
	noMatches          func()
	loadFailed         func(e FriendlyException)
}

func (h FunctionalResultHandler) TrackLoaded(track AudioTrack) {
	if h.trackLoaded != nil {
		h.trackLoaded(track)
	}
}
func (h FunctionalResultHandler) PlaylistLoaded(playlist AudioPlaylist) {
	if h.playlistLoaded != nil {
		h.playlistLoaded(playlist)
	}
}
func (h FunctionalResultHandler) SearchResultLoaded(tracks []AudioTrack) {
	if h.searchResultLoaded != nil {
		h.searchResultLoaded(tracks)
	}
}
func (h FunctionalResultHandler) NoMatches() {
	if h.noMatches != nil {
		h.noMatches()
	}
}
func (h FunctionalResultHandler) LoadFailed(e FriendlyException) {
	if h.loadFailed != nil {
		h.loadFailed(e)
	}
}
