package lavalink

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOp_UnmarshalJSON(t *testing.T) {
	data := []struct {
		payload   []byte
		opType    OpType
		eventType EventType
	}{
		{
			payload: []byte(`{"op":"playerUpdate","state":{"connected":true,"time":1645969143529},"guildId":"817327181659111454"}`),
			opType:  OpTypePlayerUpdate,
		},
		{
			payload: []byte(`{"playingPlayers":2,"op":"stats","memory":{"reservable":67108864,"used":30463224,"free":36645640,"allocated":67108864},"players":3,"cpu":{"cores":2,"systemLoad":0.49279138440159803,"lavalinkLoad":0.01450580232092837},"uptime":341757213}`),
			opType:  OpTypeStats,
		},
		{
			payload:   []byte(`{"op":"event","type":"TrackStartEvent","guildId":"817327181659111454"}`),
			opType:    OpTypeEvent,
			eventType: EventTypeTrackStart,
		},
		{
			payload:   []byte(`{"op":"event","reason":"FINISHED","type":"TrackEndEvent","track":"QAAAewIAFlNwaXJpdGJveCAtIFN1biBLaWxsZXIAC3Jpc2VyZWNvcmRzAAAAAAADeqAAC2lYWTkzYU1MRzdzAAEAK2h0dHBzOi8vd3d3LnlvdXR1YmUuY29tL3dhdGNoP3Y9aVhZOTNhTUxHN3MAB3lvdXR1YmUAAAAAAAN5OA==","guildId":"817327181659111454"}`),
			opType:    OpTypeEvent,
			eventType: EventTypeTrackEnd,
		},
		{
			payload:   []byte(`{"op":"event","type":"TrackExceptionEvent","guildId":"817327181659111454"}`),
			opType:    OpTypeEvent,
			eventType: EventTypeTrackException,
		},
		{
			payload:   []byte(`{"op":"event","type":"TrackStuckEvent","guildId":"817327181659111454"}`),
			opType:    OpTypeEvent,
			eventType: EventTypeTrackStuck,
		},
		{
			payload:   []byte(`{"op":"event","type":"WebSocketClosedEvent","guildId":"817327181659111454"}`),
			opType:    OpTypeEvent,
			eventType: EventTypeWebSocketClosed,
		},
	}

	for _, d := range data {
		var op UnmarshalOp
		err := json.Unmarshal(d.payload, &op)
		assert.NoError(t, err)
		assert.Equal(t, op.Op.Op(), d.opType)
		if eventOp, ok := op.Op.(OpEvent); ok {
			assert.Equal(t, eventOp.Event(), d.eventType)
		}
	}
}
