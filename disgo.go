package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
)

type Disgolink interface {
	Lavalink
	core.VoiceDispatchInterceptor
	core.EventListener
}
