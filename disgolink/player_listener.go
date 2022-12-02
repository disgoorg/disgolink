package disgolink

import (
	"github.com/disgoorg/disgolink/v2/lavalink"
)

type PlayerEventListener interface {
	OnEvent()

	OnPlayerPause()
	OnPlayerResume()
	OnPlayerUpdate(state lavalink.PlayerState)
	OnTrackStart(track lavalink.Track)
	OnTrackEnd(track lavalink.Track, endReason lavalink.TrackEndReason)
	OnTrackException(track lavalink.Track, exception lavalink.Exception)
	OnTrackStuck(track lavalink.Track, thresholdMs lavalink.Duration)
	OnWebSocketClosed(code int, reason string, byRemote bool)
}

type PlayerEventAdapter struct{}

func (PlayerEventAdapter) OnPlayerPause()                                                      {}
func (PlayerEventAdapter) OnPlayerResume()                                                     {}
func (PlayerEventAdapter) OnPlayerUpdate(state lavalink.PlayerState)                           {}
func (PlayerEventAdapter) OnTrackStart(track lavalink.Track)                                   {}
func (PlayerEventAdapter) OnTrackEnd(track lavalink.Track, endReason lavalink.TrackEndReason)  {}
func (PlayerEventAdapter) OnTrackException(track lavalink.Track, exception lavalink.Exception) {}
func (PlayerEventAdapter) OnTrackStuck(track lavalink.Track, thresholdMs lavalink.Duration)    {}
func (PlayerEventAdapter) OnWebSocketClosed(code int, reason string, byRemote bool)            {}
