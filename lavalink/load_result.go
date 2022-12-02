package lavalink

type LoadResult struct {
	LoadType     LoadType       `json:"loadType"`
	PlaylistInfo PlaylistInfo   `json:"playlistInfo"`
	PluginInfo   map[string]any `json:"pluginInfo"`
	Tracks       []Track        `json:"tracks"`
	Exception    Exception      `json:"exception"`
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
