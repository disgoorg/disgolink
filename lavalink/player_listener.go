package lavalink

import "github.com/disgoorg/disgolink/lavalink/protocol"

type PlayerEventListener interface {
	OnPlayerPause(player protocol.Player)
	OnPlayerResume(player protocol.Player)
	OnPlayerUpdate(player protocol.Player, state PlayerState)
	OnTrackStart(player protocol.Player, track AudioTrack)
	OnTrackEnd(player protocol.Player, track AudioTrack, endReason protocol.TrackEndReason)
	OnTrackException(player protocol.Player, track AudioTrack, exception FriendlyException)
	OnTrackStuck(player protocol.Player, track AudioTrack, thresholdMs protocol.Duration)
	OnWebSocketClosed(player protocol.Player, code int, reason string, byRemote bool)
}

type PlayerEventAdapter struct{}

func (a PlayerEventAdapter) OnPlayerPause(player protocol.Player)                     {}
func (a PlayerEventAdapter) OnPlayerResume(player protocol.Player)                    {}
func (a PlayerEventAdapter) OnPlayerUpdate(player protocol.Player, state PlayerState) {}
func (a PlayerEventAdapter) OnTrackStart(player protocol.Player, track AudioTrack)    {}
func (a PlayerEventAdapter) OnTrackEnd(player protocol.Player, track AudioTrack, endReason protocol.TrackEndReason) {
}
func (a PlayerEventAdapter) OnTrackException(player protocol.Player, track AudioTrack, exception FriendlyException) {
}
func (a PlayerEventAdapter) OnTrackStuck(player protocol.Player, track AudioTrack, thresholdMs protocol.Duration) {
}
func (a PlayerEventAdapter) OnWebSocketClosed(player protocol.Player, code int, reason string, byRemote bool) {
}
