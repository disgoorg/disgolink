package lavalink

import "github.com/disgoorg/disgolink/lavalink/protocol"

type PlayerEventListener interface {
	OnEvent()

	OnPlayerPause()
	OnPlayerResume()
	OnPlayerUpdate(state protocol.PlayerState)
	OnTrackStart(track protocol.Track)
	OnTrackEnd(track protocol.Track, endReason protocol.TrackEndReason)
	OnTrackException(track protocol.Track, exception protocol.Exception)
	OnTrackStuck(track protocol.Track, thresholdMs protocol.Duration)
	OnWebSocketClosed(code int, reason string, byRemote bool)
}

type PlayerEventAdapter struct{}

func (PlayerEventAdapter) OnPlayerPause()                                                      {}
func (PlayerEventAdapter) OnPlayerResume()                                                     {}
func (PlayerEventAdapter) OnPlayerUpdate(state protocol.PlayerState)                           {}
func (PlayerEventAdapter) OnTrackStart(track protocol.Track)                                   {}
func (PlayerEventAdapter) OnTrackEnd(track protocol.Track, endReason protocol.TrackEndReason)  {}
func (PlayerEventAdapter) OnTrackException(track protocol.Track, exception protocol.Exception) {}
func (PlayerEventAdapter) OnTrackStuck(track protocol.Track, thresholdMs protocol.Duration)    {}
func (PlayerEventAdapter) OnWebSocketClosed(code int, reason string, byRemote bool)            {}
