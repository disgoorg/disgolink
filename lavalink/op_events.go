package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
	"log"
)

type OpEventType string

const (
	OpEventTypeTrackStart     OpEventType = "TrackStartEvent"
	OpEventTypeTrackEnd       OpEventType = "TrackEndEvent"
	OpEventTypeTrackException OpEventType = "TrackExceptionEvent"
	OpEventTypeTrackStuck     OpEventType = "TrackStuckEvent"
	OpEventTypeClosed         OpEventType = "WebSocketClosedEvent"
)

type OpEvent interface {
	Op() OpType
	opEvent()
}

type UnmarshalEvent struct {
	OpEvent
}

func (e *UnmarshalEvent) UnmarshalEvent(data []byte) error {
	var opType struct {
		Op OpType `json:"op"`
	}
	if err := json.Unmarshal(data, &opType); err != nil {
		return err
	}

	var (
		event OpEvent
		err   error
	)

	switch opType.Op {
	case OpTypeEvent:
		var v EventEvent
		err = json.Unmarshal(data, &v)
		event = v
	}

	re
}

type EventEvent struct {
}

// --------------------------------------------------

type PlayerUpdateEvent struct {
	GenericOp
	GuildID discord.Snowflake `json:"guildId"`
	State   State             `json:"state"`
}

type StatsEvent struct {
	GenericOp
	*Stats
}

type GenericWebsocketEvent struct {
	GenericOp
	Type OpEventType `json:"type"`
}

type GenericPlayerEvent struct {
	GenericWebsocketEvent
	GuildID discord.Snowflake `json:"guildId"`
}

type GenericTrackEvent struct {
	GenericPlayerEvent
	RawTrack string `json:"track"`
}

func (e *GenericTrackEvent) Track() Track {
	track := &DefaultTrack{Base64Track: &e.RawTrack}
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
