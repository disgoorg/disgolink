package lavalink

import (
	"github.com/disgoorg/json"

	"github.com/disgoorg/snowflake/v2"
)

type OpType string

const (
	OpTypeReady        OpType = "ready"
	OpTypeStats        OpType = "stats"
	OpTypePlayerUpdate OpType = "playerUpdate"
	OpTypeEvent        OpType = "event"
)

type EventType string

const (
	EventTypeTrackStart      EventType = "TrackStartEvent"
	EventTypeTrackEnd        EventType = "TrackEndEvent"
	EventTypeTrackException  EventType = "TrackExceptionEvent"
	EventTypeTrackStuck      EventType = "TrackStuckEvent"
	EventTypeWebSocketClosed EventType = "WebSocketClosedEvent"
)

type Message interface {
	Op() OpType
}

type Event interface {
	Message
	Event() EventType
	GuildID() snowflake.ID
	OpEvent()
}

type UnmarshalMessage struct {
	Message
}

func (e *UnmarshalMessage) UnmarshalJSON(data []byte) error {
	var opType struct {
		Op OpType `json:"op"`
	}
	if err := json.Unmarshal(data, &opType); err != nil {
		return err
	}

	var err error

	switch opType.Op {
	case OpTypeReady:
		var v ReadyOp
		err = json.Unmarshal(data, &v)
		e.Message = v

	case OpTypePlayerUpdate:
		var v PlayerUpdateOp
		err = json.Unmarshal(data, &v)
		e.Message = v

	case OpTypeEvent:
		var v UnmarshalEvent
		err = json.Unmarshal(data, &v)
		e.Message = v.Event

	case OpTypeStats:
		var v StatsOp
		err = json.Unmarshal(data, &v)
		e.Message = v

	default:
		var v UnknownOp
		err = json.Unmarshal(data, &v)
		e.Message = v
	}

	return err
}

type ReadyOp struct {
	Resumed   bool   `json:"resumed"`
	SessionID string `json:"sessionId"`
}

func (ReadyOp) Op() OpType { return OpTypeReady }

type PlayerUpdateOp struct {
	GuildID snowflake.ID `json:"guildId"`
	State   PlayerState  `json:"state"`
}

func (PlayerUpdateOp) Op() OpType { return OpTypePlayerUpdate }

type StatsOp struct {
	Stats
}

func (StatsOp) Op() OpType { return OpTypeStats }

type UnknownOp struct {
	op   OpType
	Data []byte
}

func (o *UnknownOp) UnmarshalJSON(data []byte) error {
	var v struct {
		Op OpType `json:"op"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.op = v.Op
	o.Data = data
	return nil
}

func (o UnknownOp) MarshalJSON() ([]byte, error) {
	return o.Data, nil
}

func (o UnknownOp) Op() OpType { return o.op }
