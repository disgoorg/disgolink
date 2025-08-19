package lavalink

import (
	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

type Op string

const (
	OpReady        Op = "ready"
	OpStats        Op = "stats"
	OpPlayerUpdate Op = "playerUpdate"
	OpEvent        Op = "event"
)

type EventType string

const (
	EventTypeTrackStart      EventType = "TrackStartEvent"
	EventTypeTrackEnd        EventType = "TrackEndEvent"
	EventTypeTrackException  EventType = "TrackExceptionEvent"
	EventTypeTrackStuck      EventType = "TrackStuckEvent"
	EventTypeWebSocketClosed EventType = "WebSocketClosedEvent"

	EventTypePlayerPause  EventType = "PlayerPauseEvent"  // not actually sent by lavalink
	EventTypePlayerResume EventType = "PlayerResumeEvent" // not actually sent by lavalink
)

func UnmarshalMessage(data []byte) (Message, error) {
	var v struct {
		Op    Op        `json:"op"`
		Event EventType `json:"type"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	var (
		message Message
		err     error
	)

	switch v.Op {
	case OpReady:
		var m ReadyMessage
		err = json.Unmarshal(data, &m)
		message = m
	case OpStats:
		var m StatsMessage
		err = json.Unmarshal(data, &m)
		message = m
	case OpPlayerUpdate:
		var m PlayerUpdateMessage
		err = json.Unmarshal(data, &m)
		message = m
	case OpEvent:
		switch v.Event {
		case EventTypeTrackStart:
			var m TrackStartEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackEnd:
			var m TrackEndEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackException:
			var m TrackExceptionEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackStuck:
			var m TrackStuckEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeWebSocketClosed:
			var m WebSocketClosedEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypePlayerPause:
			var m PlayerPauseEvent
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypePlayerResume:
			var m PlayerResumeEvent
			err = json.Unmarshal(data, &m)
			message = m
		default:
			var m UnknownEvent
			err = json.Unmarshal(data, &m)
			message = m
		}
	default:
		var m UnknownMessage
		err = json.Unmarshal(data, &m)
		message = m
	}
	if err != nil {
		return nil, err
	}
	return message, nil
}

type Message interface {
	Op() Op
}

type ReadyMessage struct {
	Resumed   bool   `json:"resumed"`
	SessionID string `json:"sessionId"`
}

func (ReadyMessage) Op() Op { return OpReady }

type PlayerUpdateMessage struct {
	State   PlayerState  `json:"state"`
	GuildID snowflake.ID `json:"guildId"`
}

func (PlayerUpdateMessage) Op() Op { return OpPlayerUpdate }

type StatsMessage struct {
	Stats
}

func (StatsMessage) Op() Op { return OpStats }

type UnknownMessage struct {
	Op_  Op              `json:"op"`
	Data json.RawMessage `json:"-"`
}

func (m *UnknownMessage) UnmarshalJSON(data []byte) error {
	type unknownMessage UnknownMessage
	if err := json.Unmarshal(data, (*unknownMessage)(m)); err != nil {
		return err
	}
	m.Data = data
	return nil
}

func (m UnknownMessage) MarshalJSON() ([]byte, error) {
	return m.Data, nil
}

func (m UnknownMessage) Op() Op { return m.Op_ }

type Event interface {
	Op() Op
	Type() EventType
	GetGuildID() snowflake.ID
}

type TrackStartEvent struct {
	Track   Track        `json:"track"`
	GuildID snowflake.ID `json:"guildID"`
}

func (TrackStartEvent) Op() Op                     { return OpEvent }
func (TrackStartEvent) Type() EventType            { return EventTypeTrackStart }
func (e TrackStartEvent) GetGuildID() snowflake.ID { return e.GuildID }

type TrackEndEvent struct {
	Track   Track          `json:"track"`
	Reason  TrackEndReason `json:"reason"`
	GuildID snowflake.ID   `json:"guildID"`
}

func (TrackEndEvent) Op() Op                     { return OpEvent }
func (TrackEndEvent) Type() EventType            { return EventTypeTrackEnd }
func (e TrackEndEvent) GetGuildID() snowflake.ID { return e.GuildID }

type TrackEndReason string

const (
	TrackEndReasonFinished   TrackEndReason = "finished"
	TrackEndReasonLoadFailed TrackEndReason = "loadFailed"
	TrackEndReasonStopped    TrackEndReason = "stopped"
	TrackEndReasonReplaced   TrackEndReason = "replaced"
	TrackEndReasonCleanup    TrackEndReason = "cleanup"
)

func (e TrackEndReason) MayStartNext() bool {
	switch e {
	case TrackEndReasonFinished, TrackEndReasonLoadFailed:
		return true
	default:
		return false
	}
}

type TrackExceptionEvent struct {
	Track     Track        `json:"track"`
	Exception Exception    `json:"exception"`
	GuildID   snowflake.ID `json:"guildID"`
}

func (TrackExceptionEvent) Op() Op                     { return OpEvent }
func (TrackExceptionEvent) Type() EventType            { return EventTypeTrackException }
func (e TrackExceptionEvent) GetGuildID() snowflake.ID { return e.GuildID }

type TrackStuckEvent struct {
	Track     Track        `json:"track"`
	Threshold Duration     `json:"thresholdMs"`
	GuildID   snowflake.ID `json:"guildID"`
}

func (TrackStuckEvent) Op() Op                     { return OpEvent }
func (TrackStuckEvent) Type() EventType            { return EventTypeTrackStuck }
func (e TrackStuckEvent) GetGuildID() snowflake.ID { return e.GuildID }

type WebSocketClosedEvent struct {
	Code     int          `json:"code"`
	Reason   string       `json:"reason"`
	ByRemote bool         `json:"byRemote"`
	GuildID  snowflake.ID `json:"guildID"`
}

func (WebSocketClosedEvent) Op() Op                     { return OpEvent }
func (WebSocketClosedEvent) Type() EventType            { return EventTypeWebSocketClosed }
func (e WebSocketClosedEvent) GetGuildID() snowflake.ID { return e.GuildID }

type PlayerPauseEvent struct {
	GuildID snowflake.ID `json:"guildID"`
}

func (PlayerPauseEvent) Op() Op                     { return OpEvent }
func (PlayerPauseEvent) Type() EventType            { return EventTypePlayerPause }
func (e PlayerPauseEvent) GetGuildID() snowflake.ID { return e.GuildID }

type PlayerResumeEvent struct {
	GuildID snowflake.ID `json:"guildID"`
}

func (PlayerResumeEvent) Op() Op                     { return OpEvent }
func (PlayerResumeEvent) Type() EventType            { return EventTypePlayerResume }
func (e PlayerResumeEvent) GetGuildID() snowflake.ID { return e.GuildID }

type UnknownEvent struct {
	EventType EventType       `json:"type"`
	GuildID   snowflake.ID    `json:"guildID"`
	Data      json.RawMessage `json:"-"`
}

func (e *UnknownEvent) UnmarshalJSON(data []byte) error {
	type unknownEvent UnknownEvent
	if err := json.Unmarshal(data, (*unknownEvent)(e)); err != nil {
		return err
	}
	e.Data = data
	return nil
}

func (e UnknownEvent) MarshalJSON() ([]byte, error) {
	return e.Data, nil
}

func (UnknownEvent) Op() Op                     { return OpEvent }
func (e UnknownEvent) Type() EventType          { return e.EventType }
func (e UnknownEvent) GetGuildID() snowflake.ID { return e.GuildID }
