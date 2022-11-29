package lavalink

import (
	"fmt"

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
		var m MessageReady
		err = json.Unmarshal(data, &m)
		message = m
	case OpStats:
		var m MessageStats
		err = json.Unmarshal(data, &m)
		message = m
	case OpPlayerUpdate:
		var m MessagePlayerUpdate
		err = json.Unmarshal(data, &m)
		message = m
	case OpEvent:
		switch v.Event {
		case EventTypeTrackStart:
			var m EventTrackStart
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackEnd:
			var m EventTrackEnd
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackException:
			var m EventTrackException
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeTrackStuck:
			var m EventTrackStuck
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypeWebSocketClosed:
			var m EventWebSocketClosed
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypePlayerPause:
			var m EventPlayerPause
			err = json.Unmarshal(data, &m)
			message = m
		case EventTypePlayerResume:
			var m EventPlayerResume
			err = json.Unmarshal(data, &m)
			message = m
		default:
			err = fmt.Errorf("unknown event type: %s", v.Event)
		}
	default:
		err = fmt.Errorf("unknown op: %s", v.Op)
	}
	if err != nil {
		return nil, err
	}
	return message, nil
}

type Message interface {
	Op() Op
}

type MessageReady struct {
	Resumed   bool   `json:"resumed"`
	SessionID string `json:"sessionId"`
}

func (MessageReady) Op() Op { return OpReady }

type MessagePlayerUpdate struct {
	State   PlayerState  `json:"state"`
	GuildID snowflake.ID `json:"guildId"`
}

func (MessagePlayerUpdate) Op() Op { return OpPlayerUpdate }

type MessageStats Stats

func (MessageStats) Op() Op { return OpStats }

type Event interface {
	Type() EventType
	GuildID() snowflake.ID
}

type EventTrackStart struct {
	EncodedTrack string       `json:"encodedTrack"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (EventTrackStart) Op() Op                  { return OpEvent }
func (EventTrackStart) Type() EventType         { return EventTypeTrackStart }
func (e EventTrackStart) GuildID() snowflake.ID { return e.GuildID_ }

type EventTrackEnd struct {
	EncodedTrack string         `json:"encodedTrack"`
	Reason       TrackEndReason `json:"reason"`
	GuildID_     snowflake.ID   `json:"guildId"`
}

func (EventTrackEnd) Op() Op                  { return OpEvent }
func (EventTrackEnd) Type() EventType         { return EventTypeTrackStart }
func (e EventTrackEnd) GuildID() snowflake.ID { return e.GuildID_ }

type EventTrackException struct {
	EncodedTrack string       `json:"encodedTrack"`
	Exception    Exception    `json:"exception"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (EventTrackException) Op() Op                  { return OpEvent }
func (EventTrackException) Type() EventType         { return EventTypeTrackException }
func (e EventTrackException) GuildID() snowflake.ID { return e.GuildID_ }

type EventTrackStuck struct {
	EncodedTrack string       `json:"encodedTrack"`
	ThresholdMs  int          `json:"thresholdMs"`
	GuildID_     snowflake.ID `json:"guildId"`
}

func (EventTrackStuck) Op() Op                  { return OpEvent }
func (EventTrackStuck) Type() EventType         { return EventTypeTrackStuck }
func (e EventTrackStuck) GuildID() snowflake.ID { return e.GuildID_ }

type EventWebSocketClosed struct {
	Code     int          `json:"code"`
	Reason   string       `json:"reason"`
	ByRemote bool         `json:"byRemote"`
	GuildID_ snowflake.ID `json:"guildId"`
}

func (EventWebSocketClosed) Op() Op                  { return OpEvent }
func (EventWebSocketClosed) Type() EventType         { return EventTypeWebSocketClosed }
func (e EventWebSocketClosed) GuildID() snowflake.ID { return e.GuildID_ }

type EventPlayerPause struct {
	GuildID_ snowflake.ID `json:"guildId"`
}

func (EventPlayerPause) Op() Op                  { return OpEvent }
func (EventPlayerPause) Type() EventType         { return EventTypePlayerPause }
func (e EventPlayerPause) GuildID() snowflake.ID { return e.GuildID_ }

type EventPlayerResume struct {
	GuildID_ snowflake.ID `json:"guildId"`
}

func (EventPlayerResume) Op() Op                  { return OpEvent }
func (EventPlayerResume) Type() EventType         { return EventTypePlayerResume }
func (e EventPlayerResume) GuildID() snowflake.ID { return e.GuildID_ }
