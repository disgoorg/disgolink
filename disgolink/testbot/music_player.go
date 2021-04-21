package main

import (
	"fmt"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	dapi "github.com/DisgoOrg/disgolink/api"
)

func NewMusicPlayer(guildID string) *MusicPlayer {
	player := dgolink.Player(guildID)
	musicPlayer := &MusicPlayer{
		player: player,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

type MusicPlayer struct {
	player dapi.Player
	queue  []*dapi.Track
}

func (p *MusicPlayer) Queue(event *events.SlashCommandEvent, tracks ...*dapi.Track) {
	for _, track := range tracks {
		p.queue = append(p.queue, track)
	}

	var embed api.EmbedBuilder
	if p.player.Track() == nil {
		var track *dapi.Track
		track, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		p.player.Play(track)
		message := fmt.Sprintf("▶ ️playing [%s](%s)", track.Info.Title, *track.Info.URI)
		if len(tracks) > 1 {
			message += fmt.Sprintf("\nand queued %d tracks", len(tracks)-1)
		}
		embed.SetDescription(message)
	} else {
		embed.SetDescriptionf("queued %d tracks", len(tracks))
	}
	embed.SetFooterBy("executed by "+event.Member.EffectiveName(), event.Member.User.AvatarURL())
	if _, err := event.EditOriginal(api.NewFollowupMessageBuilder().SetEmbeds(embed.Build()).Build()); err != nil {
		logger.Errorf("error while edit original: %s", err)
	}
}

func (p *MusicPlayer) OnPlayerPause(player dapi.Player) {

}
func (p *MusicPlayer) OnPlayerResume(player dapi.Player) {

}
func (p *MusicPlayer) OnPlayerUpdate(player dapi.Player, state dapi.State) {

}
func (p *MusicPlayer) OnTrackStart(player dapi.Player, track *dapi.Track) {

}
func (p *MusicPlayer) OnTrackEnd(player dapi.Player, track *dapi.Track, endReason dapi.EndReason) {
	if endReason.MayStartNext() && len(p.queue) > 0 {
		var newTrack *dapi.Track
		newTrack, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		p.player.Play(newTrack)
	}
}
func (p *MusicPlayer) OnTrackException(player dapi.Player, track *dapi.Track, exception dapi.Exception) {

}
func (p *MusicPlayer) OnTrackStuck(player dapi.Player, track *dapi.Track, thresholdMs int) {

}
func (p *MusicPlayer) OnWebSocketClosed(player dapi.Player, code int, reason string, byRemote bool) {

}
