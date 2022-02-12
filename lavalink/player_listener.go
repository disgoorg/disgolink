package lavalink

type PlayerEventListener interface {
	OnPlayerPause(player Player)
	OnPlayerResume(player Player)
	OnPlayerUpdate(player Player, state PlayerState)
	OnTrackStart(player Player, track AudioTrack)
	OnTrackEnd(player Player, track AudioTrack, endReason AudioTrackEndReason)
	OnTrackException(player Player, track AudioTrack, exception FriendlyException)
	OnTrackStuck(player Player, track AudioTrack, thresholdMs Duration)
	OnWebSocketClosed(player Player, code int, reason string, byRemote bool)
}

type PlayerEventAdapter struct{}

func (a PlayerEventAdapter) OnPlayerPause(player Player)                     {}
func (a PlayerEventAdapter) OnPlayerResume(player Player)                    {}
func (a PlayerEventAdapter) OnPlayerUpdate(player Player, state PlayerState) {}
func (a PlayerEventAdapter) OnTrackStart(player Player, track AudioTrack)    {}
func (a PlayerEventAdapter) OnTrackEnd(player Player, track AudioTrack, endReason AudioTrackEndReason) {
}
func (a PlayerEventAdapter) OnTrackException(player Player, track AudioTrack, exception FriendlyException) {
}
func (a PlayerEventAdapter) OnTrackStuck(player Player, track AudioTrack, thresholdMs Duration) {}
func (a PlayerEventAdapter) OnWebSocketClosed(player Player, code int, reason string, byRemote bool) {
}
