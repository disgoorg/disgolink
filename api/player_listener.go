package api

type PlayerEventListener interface {
	OnPlayerPause(player Player)
	OnPlayerResume(player Player)
	OnPlayerUpdate(player Player, state State)
	OnTrackStart(player Player, track *Track)
	OnTrackEnd(player Player, track *Track, endReason EndReason)
	OnTrackException(player Player, track *Track, exception Exception)
	OnTrackStuck(player Player, track *Track, thresholdMs int)
	OnWebSocketClosed(player Player, code int, reason string, byRemote bool)
}
