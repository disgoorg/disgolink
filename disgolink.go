package disgolink

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
	"github.com/DisgoOrg/log"
	"net/http"
)

func NewLavalink(logger log.Logger, httpClient *http.Client, userID dapi.Snowflake) api.Lavalink {
	return internal.NewLavalinkImpl(logger, httpClient, userID)
}

func NewDisgolinkByUserID(logger log.Logger, httpClient *http.Client, userID dapi.Snowflake) api.Disgolink {
	return &internal.DisgolinkImpl{
		Lavalink: internal.NewLavalinkImpl(logger, httpClient, userID),
	}
}

func NewDisgolink(dgo dapi.Disgo) api.Disgolink {
	dgolink := NewDisgolinkByUserID(dgo.Logger(), dgo.RestClient().HTTPClient(), dgo.ApplicationID())
	dgo.EventManager().AddEventListeners(dgolink)
	dgo.SetVoiceDispatchInterceptor(dgolink)
	return dgolink
}
