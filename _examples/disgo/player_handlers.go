package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgolink/v4/disgolink"
	"github.com/disgoorg/disgolink/v4/lavalink"
)

func (b *Bot) onTrackStart(event *disgolink.PlayerTrackStartEvent) {
	fmt.Printf("onTrackStart: %v\n", event)
}

func (b *Bot) onTrackEnd(event *disgolink.PlayerTrackEndEvent) {
	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GetGuildID())
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
		slog.Error("Failed to play next track", slog.Any("err", err))
	}
}

func (b *Bot) onTrackException(event *disgolink.PlayerTrackExceptionEvent) {
	slog.Info("track exception", slog.Any("event", event))
}

func (b *Bot) onTrackStuck(event *disgolink.PlayerTrackStuckEvent) {
	slog.Info("track stuck", slog.Any("event", event))
}

func (b *Bot) onWebSocketClosed(event *disgolink.PlayerWebSocketClosedEvent) {
	slog.Info("websocket closed", slog.Any("event", event))
}

func (b *Bot) onUnknownPlayerEvent(event *disgolink.UnknownPlayerEvent) {
	slog.Info("unknown player event", slog.Any("op", event.Op()), slog.Any("event", event.Type()), slog.String("data", string(event.Data)))
}

func (b *Bot) onUnknownEvent(event *disgolink.UnknownEvent) {
	slog.Info("unknown event", slog.Any("op", event.Op()), slog.String("data", string(event.Data)))
}
