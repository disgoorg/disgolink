package disgolink

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/internal"
	"github.com/DisgoOrg/log"
)

func NewLavalink(logger log.Logger, userID string) api.Lavalink {
	return internal.NewLavalinkImpl(logger, userID)
}
