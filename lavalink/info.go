package lavalink

import (
	"runtime/debug"
	"strings"
)

const (
	Name   = "disgolink"
	GitHub = "https://github.com/disgoorg/" + Name
)

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
