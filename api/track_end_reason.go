package api

/*
type EndReason struct {
	Name         string
	MayStartNext bool
}

var (
	Finished   = EndReason{Name: "FINISHED", MayStartNext: true}
	LoadFailed = EndReason{Name: "LOAD_FAILED", MayStartNext: true}
	Stopped    = EndReason{Name: "STOPPED", MayStartNext: false}
	Replaced   = EndReason{Name: "REPLACED", MayStartNext: false}
	Cleanup    = EndReason{Name: "CLEANUP", MayStartNext: false}
)*/

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
