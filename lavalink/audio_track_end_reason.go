package lavalink

type AudioTrackEndReason string

const (
	AudioTrackEndReasonFinished   AudioTrackEndReason = "FINISHED"
	AudioTrackEndReasonLoadFailed AudioTrackEndReason = "LOAD_FAILED"
	AudioTrackEndReasonStopped    AudioTrackEndReason = "STOPPED"
	AudioTrackEndReasonReplaced   AudioTrackEndReason = "REPLACED"
	AudioTrackEndReasonCleanup    AudioTrackEndReason = "CLEANUP"
)

func (e AudioTrackEndReason) MayStartNext() bool {
	switch e {
	case AudioTrackEndReasonFinished, AudioTrackEndReasonLoadFailed:
		return true
	default:
		return false
	}
}
