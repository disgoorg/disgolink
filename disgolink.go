package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
	"github.com/DisgoOrg/log"
	"net/http"
)

func NewLavalink(logger log.Logger, httpClient *http.Client, userID discord.Snowflake) api.Lavalink {
	return internal.NewLavalinkImpl(logger, httpClient, userID)
}

func NewDisgolinkByUserID(logger log.Logger, httpClient *http.Client, userID discord.Snowflake) api.Disgolink {
	return &internal.DisgolinkImpl{
		Lavalink: internal.NewLavalinkImpl(logger, httpClient, userID),
	}
}

func NewDisgolink(disgo *core.Bot) api.Disgolink {
	dgolink := NewDisgolinkByUserID(disgo.Logger, disgo.RestServices.RestClient().HTTPClient(), disgo.ApplicationID)
	disgo.EventManager.AddEventListeners(dgolink)
	disgo.VoiceDispatchInterceptor = dgolink
	return dgolink
}
