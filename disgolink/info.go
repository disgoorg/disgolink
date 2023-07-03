package disgolink

import (
	"runtime/debug"
)

const (
	Name   = "disgolink"
	Module = "github.com/disgoorg/disgolink/v3"
	GitHub = "https://github.com/disgoorg/disgo"
)

var (
	Version = getVersion()
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
