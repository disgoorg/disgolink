package main

import (
	"context"
	"fmt"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgolink/v4/disgolink"
	"github.com/disgoorg/disgolink/v4/lavalink"
)

func (b *Bot) onTrackStart(event *disgolink.PlayerTrackStartEvent) {
	fmt.Printf("onTrackStart: %v\n", event)
}

func (b *Bot) onTrackEnd(event *disgolink.PlayerTrackEndEvent) {
	fmt.Printf("onTrackEnd: %v\n", event)

	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GetGuildID().String())
	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = event.Track

	case QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		return
	}
	if err := event.Player.Update(context.TODO(), disgolink.WithTrack(nextTrack)); err != nil {
		log.Error("Failed to play next track: ", err)
	}
}

func (b *Bot) onTrackException(event *disgolink.PlayerTrackExceptionEvent) {
	fmt.Printf("onTrackException: %v\n", event)
}

func (b *Bot) onTrackStuck(event *disgolink.PlayerTrackStuckEvent) {
	fmt.Printf("onTrackStuck: %v\n", event)
}

func (b *Bot) onWebSocketClosed(event *disgolink.PlayerWebSocketClosedEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
}
