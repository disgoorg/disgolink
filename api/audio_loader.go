package api

type LoadType string

const (
	LoadTypeTrackLoaded    LoadType = "TRACK_LOADED"
	LoadTypePlaylistLoaded LoadType = "PLAYLIST_LOADED"
	LoadTypeSearchResult   LoadType = "SEARCH_RESULT"
	LoadTypeNoMatches      LoadType = "NO_MATCHES"
	LoadTypeLoadFailed     LoadType = "LOAD_FAILED"
)

type Severity string

const (
	SeverityCommon     Severity = "COMMON"
	SeveritySuspicious Severity = "SUSPICIOUS"
	SeverityFault      Severity = "FAULT"
)

type LoadResult struct {
	LoadType     LoadType      `json:"loadType"`
	PlaylistInfo *PlaylistInfo `json:"playlistInfo"`
	Tracks       []*Track      `json:"tracks"`
	Exception    *Exception    `json:"exception"`
}

var _ error = (*Exception)(nil)

func NewExceptionFromErr(err error) *Exception {
	return &Exception{Message: err.Error(), Severity: SeverityFault}
}

func NewException(message string, severity Severity) *Exception {
	return &Exception{Message: message, Severity: severity}
}

type Exception struct {
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

func (e *Exception) Error() string {
	return e.Message
}

type AudioLoaderResultHandler interface {
	TrackLoaded(track *Track)
	PlaylistLoaded(playlist *Playlist)
	SearchResultLoaded(tracks []*Track)
	NoMatches()
	LoadFailed(e *Exception)
}

var _ AudioLoaderResultHandler = (*FunctionalResultHandler)(nil)

func NewResultHandler(trackLoaded func(track *Track), playlistLoaded func(playlist *Playlist), searchResultLoaded func(tracks []*Track), noMatches func(), loadFailed func(e *Exception)) AudioLoaderResultHandler {
	return &FunctionalResultHandler{trackLoaded: trackLoaded, playlistLoaded: playlistLoaded, searchResultLoaded: searchResultLoaded, noMatches: noMatches, loadFailed: loadFailed}
}

type FunctionalResultHandler struct {
	trackLoaded        func(track *Track)
	playlistLoaded     func(playlist *Playlist)
	searchResultLoaded func(tracks []*Track)
	noMatches          func()
	loadFailed         func(e *Exception)
}

func (h *FunctionalResultHandler) TrackLoaded(track *Track) {
	h.trackLoaded(track)
}
func (h *FunctionalResultHandler) PlaylistLoaded(playlist *Playlist) {
	h.playlistLoaded(playlist)
}
func (h *FunctionalResultHandler) SearchResultLoaded(tracks []*Track) {
	h.searchResultLoaded(tracks)
}
func (h *FunctionalResultHandler) NoMatches() {
	h.noMatches()
}
func (h *FunctionalResultHandler) LoadFailed(e *Exception) {
	h.loadFailed(e)
}
