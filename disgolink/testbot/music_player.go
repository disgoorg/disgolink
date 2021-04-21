package main

import "github.com/DisgoOrg/disgolink/api"

func NewMusicPlayer(player api.Player) *MusicPlayer {
	musicPlayer := &MusicPlayer{
		player: player,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

type MusicPlayer struct {
	player api.Player
	queue  []*api.Track
}

func (p MusicPlayer) Queue(tracks ...*api.Track) {
	p.queue = append(p.queue, tracks...)

	if p.player.Track() == nil {
		var track *api.Track
		track, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		p.player.Play(track)
	}
}

func (p *MusicPlayer) OnPlayerPause(player api.Player) {
	logger.Infof("OnPlayerPause")
}
func (p *MusicPlayer) OnPlayerResume(player api.Player) {
	logger.Infof("OnPlayerResume")
}
func (p *MusicPlayer) OnPlayerUpdate(player api.Player, state api.State) {
	logger.Infof("OnPlayerUpdate")
}
func (p *MusicPlayer) OnTrackStart(player api.Player, track *api.Track) {
	logger.Infof("OnTrackStart")
}
func (p *MusicPlayer) OnTrackEnd(player api.Player, track *api.Track, endReason api.EndReason) {
	logger.Infof("OnTrackEnd")
}
func (p *MusicPlayer) OnTrackException(player api.Player, track *api.Track, exception api.Exception) {
	logger.Infof("OnTrackException")
}
func (p *MusicPlayer) OnTrackStuck(player api.Player, track *api.Track, thresholdMs int) {
	logger.Infof("OnTrackStuck")
}
func (p *MusicPlayer) OnWebSocketClosed(player api.Player, code int, reason string, byRemote bool) {
	logger.Infof("OnWebSocketClosed")
}
