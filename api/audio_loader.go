package api

type LoadType string

const (
	TrackLoaded    LoadType = "TRACK_LOADED"
	PlaylistLoaded LoadType = "PLAYLIST_LOADED"
	SearchResult   LoadType = "SEARCH_RESULT"
	NoMatches      LoadType = "NO_MATCHES"
	LoadFailed     LoadType = "LOAD_FAILED"
)

type Severity string

const (
	Common     Severity = "COMMON"
	Suspicious Severity = "SUSPICIOUS"
	Fault      Severity = "FAULT"
)

var _ error = (*Error)(nil)

type Error string

func (e Error) Error() string {
	return string(e)
}

type LoadResult struct {
	LoadType     LoadType     `json:"loadType"`
	PlaylistInfo PlaylistInfo `json:"playlistInfo"`
	Tracks       []Track      `json:"tracks"`
	Exception    Exception    `json:"exception"`
}

type Exception struct {
	Error    Error    `json:"message"`
	Severity Severity `json:"severity"`
}

type AudioLoaderResultHandler interface {
	TrackLoaded(track Track)
	PlaylistLoaded(playlist Playlist)
	NoMatches()
	LoadFailed(e Exception)
}

var _ AudioLoaderResultHandler = (*FunctionalResultHandler)(nil)

func NewFunctionalResultHandler(trackLoaded func(track Track), playlistLoaded func(playlist Playlist), noMatches func(), loadFailed func(e Exception)) AudioLoaderResultHandler {
	return &FunctionalResultHandler{trackLoaded: trackLoaded, playlistLoaded: playlistLoaded, noMatches: noMatches, loadFailed: loadFailed}
}

type FunctionalResultHandler struct {
	trackLoaded    func(track Track)
	playlistLoaded func(playlist Playlist)
	noMatches      func()
	loadFailed     func(e Exception)
}

func (h *FunctionalResultHandler) TrackLoaded(track Track) {
	h.trackLoaded(track)
}
func (h *FunctionalResultHandler) PlaylistLoaded(playlist Playlist) {
	h.playlistLoaded(playlist)
}
func (h *FunctionalResultHandler) NoMatches() {
	h.noMatches()
}
func (h *FunctionalResultHandler) LoadFailed(e Exception) {
	h.loadFailed(e)
}
