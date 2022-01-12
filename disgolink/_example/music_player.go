package main

import (
	"fmt"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/log"

	"github.com/DisgoOrg/disgo/core"
)

func NewMusicPlayer(guildID discord.Snowflake) *MusicPlayer {
	player := dgolink.Player(guildID.String())
	musicPlayer := &MusicPlayer{
		Player: player,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

type MusicPlayer struct {
	lavalink.Player
	queue   []lavalink.Track
	channel core.MessageChannel
}

func (p *MusicPlayer) Queue(event *events.ApplicationCommandInteractionEvent, skipSegments bool, tracks ...lavalink.Track) {
	p.channel = event.Channel()
	for _, track := range tracks {
		p.queue = append(p.queue, track)
	}

	var embed discord.EmbedBuilder
	if p.Track() == nil {
		var track lavalink.Track
		track, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		_ = p.Play(track)
		message := fmt.Sprintf("▶ ️playing [%s](%s)", track.Info().Title(), *track.Info().URI())
		if len(tracks) > 1 {
			message += fmt.Sprintf("\nand queued %d tracks", len(tracks)-1)
		}
		embed.SetDescription(message)
	} else {
		embed.SetDescriptionf("queued %d tracks", len(tracks))
	}
	embed.SetFooter("executed by "+event.Member.EffectiveName(), event.User.EffectiveAvatarURL(1024))
	if _, err := event.UpdateResponse(discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build()); err != nil {
		log.Errorf("error while edit original: %s", err)
	}
}

func (p *MusicPlayer) OnPlayerPause(player lavalink.Player) {

}
func (p *MusicPlayer) OnPlayerResume(player lavalink.Player) {

}
func (p *MusicPlayer) OnPlayerUpdate(player lavalink.Player, state lavalink.PlayerState) {
	log.Infof("player update: %+v", state)
}
func (p *MusicPlayer) OnTrackStart(player lavalink.Player, track lavalink.Track) {

}
func (p *MusicPlayer) OnTrackEnd(player lavalink.Player, track lavalink.Track, endReason lavalink.TrackEndReason) {
	if endReason.MayStartNext() && len(p.queue) > 0 {
		var newTrack lavalink.Track
		newTrack, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		_ = player.Play(newTrack)
	}
}
func (p *MusicPlayer) OnTrackException(player lavalink.Player, track lavalink.Track, exception lavalink.Exception) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Track exception: `%s`, `%s`, `%+v`", track.Info().Title(), exception).Build())
}
func (p *MusicPlayer) OnTrackStuck(player lavalink.Player, track lavalink.Track, thresholdMs int) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("track stuck: `%s`, %d", track.Info().Title(), thresholdMs).Build())
}
func (p *MusicPlayer) OnWebSocketClosed(player lavalink.Player, code int, reason string, byRemote bool) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("websocket closed: `%d`, `%s`, `%t`", code, reason, byRemote).Build())
}
