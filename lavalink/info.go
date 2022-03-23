package lavalink

import (
	"runtime/debug"
	"strings"
)

//goland:noinspection GoUnusedConst
const (
	GitHub = "https://github.com/disgoorg/disgolink"
	Name   = "disgolink"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Version = getVersion()
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(GitHub, dep.Path) {
				return dep.Version
			}
		}
	}
	return "unknown"
}
