package api

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"log"
)

type WebsocketEvent string

const (
	WebsocketEventTrackStart     WebsocketEvent = "TrackStartEvent"
	WebsocketEventTrackEnd       WebsocketEvent = "TrackEndEvent"
	WebsocketEventTrackException WebsocketEvent = "TrackExceptionEvent"
	WebsocketEventTrackStuck     WebsocketEvent = "TrackStuckEvent"
	WebSocketEventClosed         WebsocketEvent = "WebSocketClosedEvent"
)

type PlayerUpdateEvent struct {
	GenericOp
	GuildID dapi.Snowflake `json:"guildId"`
	State   State         `json:"state"`
}

type StatsEvent struct {
	GenericOp
	*Stats
}

type GenericWebsocketEvent struct {
	GenericOp
	Type WebsocketEvent `json:"type"`
}

type GenericPlayerEvent struct {
	GenericWebsocketEvent
	GuildID dapi.Snowflake `json:"guildId"`
}

type GenericTrackEvent struct {
	GenericPlayerEvent
	RawTrack string `json:"track"`
}

func (e *GenericTrackEvent) Track() Track {
	track := &DefaultTrack{Track_: &e.RawTrack}
	err := track.DecodeInfo()
	if err != nil {
		// TODO access normal logger
		log.Printf("error while unpacking track info: %s", err)
	}
	return track
}

type TrackStartEvent struct {
	GenericTrackEvent
}

type TrackEndEvent struct {
	GenericTrackEvent
	EndReason EndReason `json:"reason"`
}

type TrackExceptionEvent struct {
	GenericTrackEvent
	Exception Exception `json:"exception"`
}

type TrackStuckEvent struct {
	GenericTrackEvent
	ThresholdMs int `json:"thresholdMs"`
}

type State struct {
	Time      int  `json:"time"`
	Position  int  `json:"position"`
	Connected bool `json:"connected"`
}

type WebSocketClosedEvent struct {
	GenericPlayerEvent
	Code     int    `json:"code"`
	Reason   string `json:"reason"`
	ByRemote bool   `json:"byRemote"`
}
