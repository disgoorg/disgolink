package events

import (
	"github.com/DisgoOrg/disgolink/api"
)

type PlayerEvent interface {
	Player() api.Player
}

type genericPlayerEvent struct {
	player api.Player
}

func (e genericPlayerEvent) Player() api.Player {
	return e.player
}

type PlayerPauseEvent struct {
	genericPlayerEvent
}

type PlayerResumeEvent struct {
	genericPlayerEvent
}
