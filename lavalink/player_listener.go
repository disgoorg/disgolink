package lavalink

type PlayerEventListener interface {
	OnPlayerPause(player Player)
	OnPlayerResume(player Player)
	OnPlayerUpdate(player Player, state PlayerState)
	OnTrackStart(player Player, track Track)
	OnTrackEnd(player Player, track Track, endReason TrackEndReason)
	OnTrackException(player Player, track Track, exception Exception)
	OnTrackStuck(player Player, track Track, thresholdMs int)
	OnWebSocketClosed(player Player, code int, reason string, byRemote bool)
}
