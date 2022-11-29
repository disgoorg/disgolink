package lavalink

type LoadResult struct {
	LoadType  LoadType   `json:"loadType"`
	Playlist  *Playlist  `json:"playlist"`
	Tracks    []Track    `json:"tracks"`
	Exception *Exception `json:"exception"`
}

var _ error = (*Exception)(nil)

type Exception struct {
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Cause    *string  `json:"cause,omitempty"`
}

func (e Exception) Error() string {
	return e.Message
}

type Severity string

const (
	SeverityCommon     Severity = "COMMON"
	SeveritySuspicious Severity = "SUSPICIOUS"
	SeverityFault      Severity = "FAULT"
)

type TrackEndReason string

const (
	TrackEndReasonFinished   TrackEndReason = "FINISHED"
	TrackEndReasonLoadFailed TrackEndReason = "LOAD_FAILED"
	TrackEndReasonStopped    TrackEndReason = "STOPPED"
	TrackEndReasonReplaced   TrackEndReason = "REPLACED"
	TrackEndReasonCleanup    TrackEndReason = "CLEANUP"
)

func (e TrackEndReason) MayStartNext() bool {
	switch e {
	case TrackEndReasonFinished, TrackEndReasonLoadFailed:
		return true
	default:
		return false
	}
}

type LoadType string

const (
	LoadTypeTrackLoaded    LoadType = "TRACK_LOADED"
	LoadTypePlaylistLoaded LoadType = "PLAYLIST_LOADED"
	LoadTypeSearchResult   LoadType = "SEARCH_RESULT"
	LoadTypeNoMatches      LoadType = "NO_MATCHES"
	LoadTypeLoadFailed     LoadType = "LOAD_FAILED"
)

type AudioLoadResultHandler interface {
	TrackLoaded(track Track)
	PlaylistLoaded(playlist Playlist)
	SearchResultLoaded(tracks []Track)
	NoMatches()
	LoadFailed(e Exception)
}

var _ AudioLoadResultHandler = (*FunctionalResultHandler)(nil)

func NewResultHandler(trackLoaded func(track Track), playlistLoaded func(playlist Playlist), searchResultLoaded func(tracks []Track), noMatches func(), loadFailed func(e Exception)) AudioLoadResultHandler {
	return FunctionalResultHandler{
		trackLoaded:        trackLoaded,
		playlistLoaded:     playlistLoaded,
		searchResultLoaded: searchResultLoaded,
		noMatches:          noMatches,
		loadFailed:         loadFailed,
	}
}

type FunctionalResultHandler struct {
	trackLoaded        func(track Track)
	playlistLoaded     func(playlist Playlist)
	searchResultLoaded func(tracks []Track)
	noMatches          func()
	loadFailed         func(e Exception)
}

func (h FunctionalResultHandler) TrackLoaded(track Track) {
	if h.trackLoaded != nil {
		h.trackLoaded(track)
	}
}
func (h FunctionalResultHandler) PlaylistLoaded(playlist Playlist) {
	if h.playlistLoaded != nil {
		h.playlistLoaded(playlist)
	}
}
func (h FunctionalResultHandler) SearchResultLoaded(tracks []Track) {
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
