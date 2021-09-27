package disgolink

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
	"net/http"
)

func NewLavalink(logger log.Logger, httpClient *http.Client, userID discord.Snowflake) Lavalink {
	return newDefaultLavalink(logger, httpClient, userID)
}

func NewDisgolinkByUserID(logger log.Logger, httpClient *http.Client, userID discord.Snowflake) Disgolink {
	return &defaultDisgolink{
		Lavalink: newDefaultLavalink(logger, httpClient, userID),
	}
}

func NewDisgolink(disgo *core.Bot) Disgolink {
	dgolink := NewDisgolinkByUserID(disgo.Logger, disgo.RestServices.HTTPClient(), disgo.ApplicationID)
	disgo.EventManager.AddEventListeners(dgolink)
	disgo.EventManager.SetVoiceDispatchInterceptor(dgolink)
	return dgolink
}
