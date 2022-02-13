package main

import (
	"fmt"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"

	"github.com/DisgoOrg/disgo/core"
)

func NewMusicPlayer(guildID snowflake.Snowflake) *MusicPlayer {
	player := dgolink.Player(guildID)
	musicPlayer := &MusicPlayer{
		Player: player,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

type MusicPlayer struct {
	lavalink.Player
	queue   []lavalink.AudioTrack
	channel core.MessageChannel
}

func (p *MusicPlayer) Queue(event *events.ApplicationCommandInteractionEvent, skipSegments bool, tracks ...lavalink.AudioTrack) {
	p.channel = event.Channel()
	for _, track := range tracks {
		p.queue = append(p.queue, track)
	}

	var embed discord.EmbedBuilder
	if p.Track() == nil {
		var track lavalink.AudioTrack
		track, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		_ = p.Play(track)
		message := fmt.Sprintf("▶ ️playing [%s](%s)", track.Info().Title, *track.Info().URI)
		if len(tracks) > 1 {
			message += fmt.Sprintf("\nand queued %d tracks", len(tracks)-1)
		}
		embed.SetDescription(message)
	} else {
		embed.SetDescriptionf("queued %d tracks", len(tracks))
	}
	embed.SetFooter("executed by "+event.Member.EffectiveName(), event.User.EffectiveAvatarURL(1024))
	if _, err := event.UpdateOriginalMessage(discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build()); err != nil {
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
func (p *MusicPlayer) OnTrackStart(player lavalink.Player, track lavalink.AudioTrack) {

}
func (p *MusicPlayer) OnTrackEnd(player lavalink.Player, track lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	if endReason.MayStartNext() && len(p.queue) > 0 {
		var newTrack lavalink.AudioTrack
		newTrack, p.queue = p.queue[len(p.queue)-1], p.queue[:len(p.queue)-1]
		_ = player.Play(newTrack)
	}
}
func (p *MusicPlayer) OnTrackException(player lavalink.Player, track lavalink.AudioTrack, exception lavalink.FriendlyException) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("AudioTrack exception: `%s`, `%+v`", track.Info().Title, exception).Build())
}
func (p *MusicPlayer) OnTrackStuck(player lavalink.Player, track lavalink.AudioTrack, thresholdMs int) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("track stuck: `%s`, %d", track.Info().Title, thresholdMs).Build())
}
func (p *MusicPlayer) OnWebSocketClosed(player lavalink.Player, code int, reason string, byRemote bool) {
	_, _ = p.channel.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("websocket closed: `%d`, `%s`, `%t`", code, reason, byRemote).Build())
}
