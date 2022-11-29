package lavalink

type PlayerEventListener interface {
	OnEvent()

	OnPlayerPause()
	OnPlayerResume()
	OnPlayerUpdate(state PlayerState)
	OnTrackStart(track Track)
	OnTrackEnd(track Track, endReason TrackEndReason)
	OnTrackException(track Track, exception Exception)
	OnTrackStuck(track Track, thresholdMs Duration)
	OnWebSocketClosed(code int, reason string, byRemote bool)
}

type PlayerEventAdapter struct{}

func (PlayerEventAdapter) OnPlayerPause()                                           {}
func (PlayerEventAdapter) OnPlayerResume()                                          {}
func (PlayerEventAdapter) OnPlayerUpdate(state PlayerState)                         {}
func (PlayerEventAdapter) OnTrackStart(track Track)                                 {}
func (PlayerEventAdapter) OnTrackEnd(track Track, endReason TrackEndReason)         {}
func (PlayerEventAdapter) OnTrackException(track Track, exception Exception)        {}
func (PlayerEventAdapter) OnTrackStuck(track Track, thresholdMs Duration)           {}
func (PlayerEventAdapter) OnWebSocketClosed(code int, reason string, byRemote bool) {}
