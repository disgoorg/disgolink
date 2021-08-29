package main

import (
	"fmt"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgolink/api"
)

func NewMusicPlayer(guildID discord.Snowflake) *MusicPlayer {
	player := dgolink.Player(guildID)
	musicPlayer := &MusicPlayer{
		player: player,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

type MusicPlayer struct {
	player      api.Player
	queue       []api.Track
	textChannel core.TextChannel
}

func (p *MusicPlayer) Queue(event *events.SlashCommandEvent, tracks ...api.Track) {
	p.textChannel = event.Interaction.TextChannel()
	for _, track := range tracks {
		p.queue = append(p.queue, track)
	}

	var embed core.EmbedBuilder
	if p.player.Track() == nil {
		var track api.Track
		track, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		p.player.Play(track)
		message := fmt.Sprintf("▶ ️playing [%s](%s)", track.Info().Title(), *track.Info().URI())
		if len(tracks) > 1 {
			message += fmt.Sprintf("\nand queued %d tracks", len(tracks)-1)
		}
		embed.SetDescription(message)
	} else {
		embed.SetDescriptionf("queued %d tracks", len(tracks))
	}
	embed.SetFooter("executed by "+event.Interaction.Member.EffectiveName(), event.Interaction.Member.User.EffectiveAvatarURL(1024))
	if _, err := event.UpdateOriginal(core.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build()); err != nil {
		log.Errorf("error while edit original: %s", err)
	}
}

func (p *MusicPlayer) OnPlayerPause(player api.Player) {

}
func (p *MusicPlayer) OnPlayerResume(player api.Player) {

}
func (p *MusicPlayer) OnPlayerUpdate(player api.Player, state api.State) {
	log.Infof("player update: %+v", state)
}
func (p *MusicPlayer) OnTrackStart(player api.Player, track api.Track) {

}
func (p *MusicPlayer) OnTrackEnd(player api.Player, track api.Track, endReason api.EndReason) {
	if endReason.MayStartNext() && len(p.queue) > 0 {
		var newTrack api.Track
		newTrack, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		player.Play(newTrack)
	}
}
func (p *MusicPlayer) OnTrackException(player api.Player, track api.Track, exception api.Exception) {
	_, _ = p.textChannel.CreateMessage(core.NewMessageCreateBuilder().SetContentf("Track exception: `%s`, `%s`, `%+v`", track.Info().Title(), exception).Build())
}
func (p *MusicPlayer) OnTrackStuck(player api.Player, track api.Track, thresholdMs int) {
	_, _ = p.textChannel.CreateMessage(core.NewMessageCreateBuilder().SetContentf("track stuck: `%s`, %d", track.Info().Title(), thresholdMs).Build())
}
func (p *MusicPlayer) OnWebSocketClosed(player api.Player, code int, reason string, byRemote bool) {
	_, _ = p.textChannel.CreateMessage(core.NewMessageCreateBuilder().SetContentf("websocket closed: `%d`, `%s`, `%t`", code, reason, byRemote).Build())
}
