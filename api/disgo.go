package api

import (
	dapi "github.com/DisgoOrg/disgo/api"
)

type Disgolink interface {
	Lavalink
	dapi.VoiceDispatchInterceptor
	dapi.EventListener
}

