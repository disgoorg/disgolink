package lavalink

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
)

type OpType string

const (
	OpTypeReady             OpType = "ready"
	OpTypePlay              OpType = "play"
	OpTypeStop              OpType = "stop"
	OpTypePause             OpType = "pause"
	OpTypeSeek              OpType = "seek"
	OpTypeVolume            OpType = "volume"
	OpTypeDestroy           OpType = "destroy"
	OpTypeStats             OpType = "stats"
	OpTypeVoiceUpdate       OpType = "voiceUpdate"
	OpTypePlayerUpdate      OpType = "playerUpdate"
	OpTypeEvent             OpType = "event"
	OpTypeConfigureResuming OpType = "configureResuming"
	OpTypeFilters           OpType = "filters"
)

type EventType string

const (
	EventTypeTrackStart      EventType = "TrackStartEvent"
	EventTypeTrackEnd        EventType = "TrackEndEvent"
	EventTypeTrackException  EventType = "TrackExceptionEvent"
	EventTypeTrackStuck      EventType = "TrackStuckEvent"
	EventTypeWebSocketClosed EventType = "WebSocketClosedEvent"
)

type Op interface {
	Op() OpType
}

type OpCommand interface {
	json.Marshaler
	Op
	OpCommand()
}

type OpEvent interface {
	Op
	Event() EventType
	GuildID() snowflake.ID
	OpEvent()
}

type UnmarshalOp struct {
	Op
}

func (e *UnmarshalOp) UnmarshalJSON(data []byte) error {
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
		e.Op = v

	case OpTypePlayerUpdate:
		var v PlayerUpdateOp
		err = json.Unmarshal(data, &v)
		e.Op = v

	case OpTypeEvent:
		var v UnmarshalOpEvent
		err = json.Unmarshal(data, &v)
		e.Op = v.OpEvent

	case OpTypeStats:
		var v StatsOp
		err = json.Unmarshal(data, &v)
		e.Op = v

	default:
		var v UnknownOp
		err = json.Unmarshal(data, &v)
		e.Op = v
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
