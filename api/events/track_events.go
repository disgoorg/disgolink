package events

import (
	"encoding/json"
	"log"

	"github.com/DisgoOrg/disgolink/api"
)

type RawTrackEvent struct {
	GuildID string          `json:"guildId"`
	Type    string          `json:"type"`
	Event   json.RawMessage `json:"event"`
}

type TrackEvent interface {
	Track() api.Track
}

type genericTrackEvent struct {
	genericPlayerEvent
	RawTrack string `json:"track"`
	track *api.Track
}

func (e genericTrackEvent) Track() *api.Track {
	if e.track == nil {
		e.track = &api.Track{Track: e.RawTrack}
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
	EndReason api.EndReason `json:"reason"`
}

type TrackExceptionEvent struct {
	genericTrackEvent
	Exception api.Exception `json:"exception"`
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
