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

type AudioLoaderResultHandler interface {
	TrackLoaded(track AudioTrack)
	PlaylistLoaded(playlist Playlist)
	SearchResultLoaded(tracks []AudioTrack)
	NoMatches()
	LoadFailed(e Exception)
}

type LoadResult struct {
	LoadType     LoadType      `json:"loadType"`
	PlaylistInfo *PlaylistInfo `json:"playlistInfo"`
	Tracks       []AudioTrack  `json:"tracks"`
	Exception    *Exception    `json:"exception"`
}

func (r *LoadResult) UnmarshalJSON(data []byte) error {
	var v struct {
		LoadType     LoadType        `json:"loadType"`
		PlaylistInfo *PlaylistInfo   `json:"playlistInfo"`
		Tracks       []*DefaultTrack `json:"tracks"`
		Exception    *Exception      `json:"exception"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.LoadType = v.LoadType
	r.PlaylistInfo = v.PlaylistInfo
	r.Tracks = defaultTracksToTracks(v.Tracks)
	r.Exception = v.Exception
	return nil
}

func defaultTracksToTracks(defaultTracks []*DefaultTrack) []AudioTrack {
	if defaultTracks == nil {
		return nil
	}
	tracks := make([]AudioTrack, len(defaultTracks))
	for i := 0; i < len(defaultTracks); i++ {
		tracks[i] = defaultTracks[i]
	}
	return tracks
}

var _ AudioLoaderResultHandler = (*FunctionalResultHandler)(nil)

func NewResultHandler(trackLoaded func(track AudioTrack), playlistLoaded func(playlist Playlist), searchResultLoaded func(tracks []AudioTrack), noMatches func(), loadFailed func(e Exception)) AudioLoaderResultHandler {
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
	playlistLoaded     func(playlist Playlist)
	searchResultLoaded func(tracks []AudioTrack)
	noMatches          func()
	loadFailed         func(e Exception)
}

func (h FunctionalResultHandler) TrackLoaded(track AudioTrack) {
	if h.trackLoaded != nil {
		h.trackLoaded(track)
	}
}
func (h FunctionalResultHandler) PlaylistLoaded(playlist Playlist) {
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
func (h FunctionalResultHandler) LoadFailed(e Exception) {
	if h.loadFailed != nil {
		h.loadFailed(e)
	}
}
