package lavalink

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type OpType string

const (
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
	GuildID() string
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

	var (
		op  Op
		err error
	)

	switch opType.Op {
	case OpTypePlayerUpdate:
		var v PlayerUpdateOp
		err = json.Unmarshal(data, &v)
		op = v

	case OpTypeEvent:
		var v UnmarshalOpEvent
		err = json.Unmarshal(data, &v)
		op = v.OpEvent

	case OpTypeStats:
		var v StatsOp
		err = json.Unmarshal(data, &v)
		op = v

	default:
		return errors.Errorf("unknown op type %s received", opType.Op)
	}

	if err != nil {
		return err
	}

	e.Op = op
	return nil
}

type PlayerUpdateOp struct {
	GuildID string      `json:"guildId"`
	State   PlayerState `json:"state"`
}

func (PlayerUpdateOp) Op() OpType { return OpTypePlayerUpdate }

type StatsOp struct {
	Stats
}

func (StatsOp) Op() OpType { return OpTypeStats }
