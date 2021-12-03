package lavalink

type EndReason string

const (
	Finished   EndReason = "FINISHED"
	LoadFailed EndReason = "LOAD_FAILED"
	Stopped    EndReason = "STOPPED"
	Replaced   EndReason = "REPLACED"
	Cleanup    EndReason = "CLEANUP"
)

func (e EndReason) MayStartNext() bool {
	switch e {
	case Finished, LoadFailed:
		return true
	default:
		return false
	}
}
