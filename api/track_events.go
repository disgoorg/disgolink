package api

import (
	"encoding/json"
	"log"
)

type RawTrackEvent struct {
	GuildID string          `json:"guildId"`
	Type    string          `json:"type"`
	Event   json.RawMessage `json:"event"`
}

type TrackEvent interface {
	Track() Track
}

type genericTrackEvent struct {
	genericPlayerEvent
	RawTrack string `json:"track"`
	track *Track
}

func (e genericTrackEvent) Track() *Track {
	if e.track == nil {
		e.track = &Track{Track: e.RawTrack}
		err := e.track.DecodeInfo()
		if err != nil {
			log.Printf("error while unpacking track info: %s", err)
			return nil
		}
	}
	return e.track
}

type TrackEndEvent struct {
	genericTrackEvent
	EndReason EndReason `json:"reason"`
}

type TrackExceptionEvent struct {
	genericTrackEvent
	Exception Exception `json:"exception"`
}

type TrackStartEvent struct {
	genericTrackEvent
}

type TrackStuckEvent struct {
	genericTrackEvent
	thresholdMs int
}

func (e TrackStuckEvent) ThresholdMs() int {
	return e.thresholdMs
}
