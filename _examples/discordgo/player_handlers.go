package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/log"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.EventPlayerPause) {
	fmt.Printf("onPlayerPause: %v\n", event)
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.EventPlayerResume) {
	fmt.Printf("onPlayerResume: %v\n", event)
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.EventTrackStart) {
	fmt.Printf("onTrackStart: %v\n", event)
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.EventTrackEnd) {
	fmt.Printf("onTrackEnd: %v\n", event)

	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GuildID().String())
	var (
		nextTrack lavalink.Track
		ok bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = *player.Track()

	case QueueTypeRepeatQueue:
		lastTrack, _ := b.Lavalink.BestNode().DecodeTrack(context.TODO(), event.EncodedTrack)
		queue.Add(*lastTrack)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Error("Failed to play next track: ", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.EventTrackException) {
	fmt.Printf("onTrackException: %v\n", event)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.EventTrackStuck) {
	fmt.Printf("onTrackStuck: %v\n", event)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.EventWebSocketClosed) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
}
