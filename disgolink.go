package disgolink

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
	"github.com/DisgoOrg/log"
)

func NewLavalink(logger log.Logger, userID string) api.Lavalink {
	return internal.NewLavalinkImpl(logger, userID)
}

func NewDisgolinkByUserID(logger log.Logger, userID dapi.Snowflake) api.Disgolink {
	return &internal.DisgolinkImpl{
		Lavalink: internal.NewLavalinkImpl(logger, string(userID)),
	}
}

func NewDisgolink(logger log.Logger, dgo dapi.Disgo) api.Disgolink {
	dgolink := NewDisgolinkByUserID(logger, dgo.ApplicationID())
	dgo.EventManager().AddEventListeners(dgolink)
	dgo.SetVoiceDispatchInterceptor(dgolink)
	return dgolink
}
