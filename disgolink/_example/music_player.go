package main

import (
	"fmt"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

func NewMusicPlayer(client bot.Client, guildID snowflake.ID) *MusicPlayer {
	player := dgolink.Player(guildID)
	musicPlayer := &MusicPlayer{
		Player: player,
		client: client,
	}
	player.AddListener(musicPlayer)
	return musicPlayer
}

var _ lavalink.PlayerEventListener = (*MusicPlayer)(nil)

type MusicPlayer struct {
	lavalink.Player
	queue     []lavalink.AudioTrack
	client    bot.Client
	channelID snowflake.ID
}

func (p *MusicPlayer) Queue(event *events.ApplicationCommandInteractionEvent, tracks ...lavalink.AudioTrack) {
	p.channelID = event.ChannelID()
	for _, track := range tracks {
		p.queue = append(p.queue, track)
	}

	var embed discord.EmbedBuilder
	if p.PlayingTrack() == nil {
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
	embed.SetFooter("executed by "+event.Member().EffectiveName(), event.User().EffectiveAvatarURL())
	if _, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build()); err != nil {
		log.Errorf("error while edit original: %s", err)
	}
}

func (p *MusicPlayer) OnPlayerPause(player lavalink.Player) {

}
func (p *MusicPlayer) OnPlayerResume(player lavalink.Player) {

}
func (p *MusicPlayer) OnPlayerUpdate(player lavalink.Player, state lavalink.PlayerState) {
	log.Infof("player update: %d, %d", state.Position, player.Position())
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
	_, _ = p.client.Rest().CreateMessage(p.channelID, discord.NewMessageCreateBuilder().SetContentf("AudioTrack exception: `%s`, `%+v`", track.Info().Title, exception).Build())
}
func (p *MusicPlayer) OnTrackStuck(player lavalink.Player, track lavalink.AudioTrack, thresholdMs lavalink.Duration) {
	_, _ = p.client.Rest().CreateMessage(p.channelID, discord.NewMessageCreateBuilder().SetContentf("track stuck: `%s`, %d", track.Info().Title, thresholdMs).Build())
}
func (p *MusicPlayer) OnWebSocketClosed(player lavalink.Player, code int, reason string, byRemote bool) {
	_, _ = p.client.Rest().CreateMessage(p.channelID, discord.NewMessageCreateBuilder().SetContentf("websocket closed: `%d`, `%s`, `%t`", code, reason, byRemote).Build())
}
