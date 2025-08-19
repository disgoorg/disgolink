package disgolink

import (
	"runtime/debug"
)

const (
	Name   = "disgolink"
	Module = "github.com/disgoorg/disgolink/v4"
	GitHub = "https://github.com/disgoorg/disgolink"
)

var (
	Version = getVersion()

	SemVersion = "semver:" + Version
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == Module {
				return dep.Version
			}
		}
	}
	return "unknown"
}
