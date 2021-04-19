package api

type PlayerEvent interface {
	Player() Player
}

type genericPlayerEvent struct {
	player Player
}

func (e genericPlayerEvent) Player() Player {
	return e.player
}

type PlayerPauseEvent struct {
	genericPlayerEvent
}

type PlayerResumeEvent struct {
	genericPlayerEvent
}
