package lavalink

type TrackEndReason string

//goland:noinspection GoUnusedConst
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
