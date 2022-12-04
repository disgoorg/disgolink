package disgolink

import (
	"runtime/debug"
	"strings"
)

const (
	Name   = "disgolink"
	Module = "github.com/disgoorg/" + Name
	GitHub = "https://" + Module
)

var (
	Version = getVersion()
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(dep.Path, Module) {
				return dep.Version
			}
		}
	}
	return "unknown"
}
