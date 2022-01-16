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

type PlayerEventAdapter struct{}

func (a PlayerEventAdapter) OnPlayerPause(player Player)                                      {}
func (a PlayerEventAdapter) OnPlayerResume(player Player)                                     {}
func (a PlayerEventAdapter) OnPlayerUpdate(player Player, state PlayerState)                  {}
func (a PlayerEventAdapter) OnTrackStart(player Player, track Track)                          {}
func (a PlayerEventAdapter) OnTrackEnd(player Player, track Track, endReason TrackEndReason)  {}
func (a PlayerEventAdapter) OnTrackException(player Player, track Track, exception Exception) {}
func (a PlayerEventAdapter) OnTrackStuck(player Player, track Track, thresholdMs int)         {}
func (a PlayerEventAdapter) OnWebSocketClosed(player Player, code int, reason string, byRemote bool) {
}
