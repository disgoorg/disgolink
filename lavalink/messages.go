package lavalink

import (
	"github.com/disgoorg/json"
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
	EventTypePlayerPause     EventType = "PlayerPauseEvent"  // not actually sent by lavalink
	EventTypePlayerResume    EventType = "PlayerResumeEvent" // not actually sent by lavalink
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

type StatsMessage Stats

func (StatsMessage) Op() Op { return OpStats }

type UnknownMessage struct {
	Op_  Op              `json:"op"`
	Data json.RawMessage `json:"-"`
}

func (e *UnknownMessage) UnmarshalJSON(data []byte) error {
	type unknownMessage UnknownMessage
	if err := json.Unmarshal(data, (*unknownMessage)(e)); err != nil {
		return err
	}
	e.Data = data
	return nil
}

func (e UnknownMessage) MarshalJSON() ([]byte, error) {
	return e.Data, nil
}

func (m UnknownMessage) Op() Op { return m.Op_ }

type Event interface {
	Type() EventType
	GuildID() snowflake.ID
}

type TrackStartEvent struct {
	EncodedTrack string       `json:"encodedTrack"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (TrackStartEvent) Op() Op                  { return OpEvent }
func (TrackStartEvent) Type() EventType         { return EventTypeTrackStart }
func (e TrackStartEvent) GuildID() snowflake.ID { return e.GuildID_ }

type TrackEndEvent struct {
	EncodedTrack string         `json:"encodedTrack"`
	Reason       TrackEndReason `json:"reason"`
	GuildID_     snowflake.ID   `json:"guildId"`
}

func (TrackEndEvent) Op() Op                  { return OpEvent }
func (TrackEndEvent) Type() EventType         { return EventTypeTrackStart }
func (e TrackEndEvent) GuildID() snowflake.ID { return e.GuildID_ }

type TrackExceptionEvent struct {
	EncodedTrack string       `json:"encodedTrack"`
	Exception    Exception    `json:"exception"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (TrackExceptionEvent) Op() Op                  { return OpEvent }
func (TrackExceptionEvent) Type() EventType         { return EventTypeTrackException }
func (e TrackExceptionEvent) GuildID() snowflake.ID { return e.GuildID_ }

type TrackStuckEvent struct {
	EncodedTrack string       `json:"encodedTrack"`
	ThresholdMs  int          `json:"thresholdMs"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (TrackStuckEvent) Op() Op                  { return OpEvent }
func (TrackStuckEvent) Type() EventType         { return EventTypeTrackStuck }
func (e TrackStuckEvent) GuildID() snowflake.ID { return e.GuildID_ }

type WebSocketClosedEvent struct {
	Code     int          `json:"code"`
	Reason   string       `json:"reason"`
	ByRemote bool         `json:"byRemote"`
	GuildID_ snowflake.ID `json:"guildId"`
}

func (WebSocketClosedEvent) Op() Op                  { return OpEvent }
func (WebSocketClosedEvent) Type() EventType         { return EventTypeWebSocketClosed }
func (e WebSocketClosedEvent) GuildID() snowflake.ID { return e.GuildID_ }

type PlayerPauseEvent struct {
	GuildID_ snowflake.ID `json:"guildId"`
}

func (PlayerPauseEvent) Op() Op                  { return OpEvent }
func (PlayerPauseEvent) Type() EventType         { return EventTypePlayerPause }
func (e PlayerPauseEvent) GuildID() snowflake.ID { return e.GuildID_ }

type PlayerResumeEvent struct {
	GuildID_ snowflake.ID `json:"guildId"`
}

func (PlayerResumeEvent) Op() Op                  { return OpEvent }
func (PlayerResumeEvent) Type() EventType         { return EventTypePlayerResume }
func (e PlayerResumeEvent) GuildID() snowflake.ID { return e.GuildID_ }

type UnknownEvent struct {
	Type_    EventType       `json:"type"`
	GuildID_ snowflake.ID    `json:"guildId"`
	Data     json.RawMessage `json:"-"`
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

func (UnknownEvent) Op() Op                  { return OpEvent }
func (e UnknownEvent) Type() EventType       { return e.Type_ }
func (e UnknownEvent) GuildID() snowflake.ID { return e.GuildID_ }
