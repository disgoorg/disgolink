package lavalink

import "encoding/json"

type UnmarshalOpEvent struct {
	OpEvent
}

func (e *UnmarshalOpEvent) UnmarshalJSON(data []byte) error {
	var eType struct {
		Type EventType `json:"type"`
	}
	if err := json.Unmarshal(data, &eType); err != nil {
		return err
	}

	var (
		opEvent OpEvent
		err     error
	)

	switch eType.Type {
	case EventTypeTrackStart:
		var v TrackStartEvent
		err = json.Unmarshal(data, &v)
		opEvent = v

	case EventTypeTrackEnd:
		var v TrackEndEvent
		err = json.Unmarshal(data, &v)
		opEvent = v

	case EventTypeTrackException:
		var v TrackExceptionEvent
		err = json.Unmarshal(data, &v)
		opEvent = v

	case EventTypeTrackStuck:
		var v TrackStuckEvent
		err = json.Unmarshal(data, &v)
		opEvent = v

	case EventTypeWebSocketClosed:
		var v WebsocketClosedEvent
		err = json.Unmarshal(data, &v)
		opEvent = v

	default:
		return nil
	}

	if err != nil {
		return err
	}

	e.OpEvent = opEvent
	return nil
}

type TrackStartEvent struct {
	GID   string `json:"guildId"`
	Track string `json:"track"`
}

func (TrackStartEvent) Event() EventType  { return EventTypeTrackStart }
func (TrackStartEvent) Op() OpType        { return OpTypeEvent }
func (e TrackStartEvent) GuildID() string { return e.GID }
func (TrackStartEvent) OpEvent()          {}

type TrackEndEvent struct {
	GID    string         `json:"guildId"`
	Track  string         `json:"track"`
	Reason TrackEndReason `json:"reason"`
}

func (TrackEndEvent) Event() EventType  { return EventTypeTrackStart }
func (TrackEndEvent) Op() OpType        { return OpTypeEvent }
func (e TrackEndEvent) GuildID() string { return e.GID }
func (TrackEndEvent) OpEvent()          {}

type TrackExceptionEvent struct {
	GID       string            `json:"guildId"`
	Track     string            `json:"track"`
	Exception FriendlyException `json:"exception"`
}

func (TrackExceptionEvent) Event() EventType  { return EventTypeTrackStart }
func (TrackExceptionEvent) Op() OpType        { return OpTypeEvent }
func (e TrackExceptionEvent) GuildID() string { return e.GID }
func (TrackExceptionEvent) OpEvent()          {}

type TrackStuckEvent struct {
	GID         string `json:"guildId"`
	Track       string `json:"track"`
	ThresholdMs int    `json:"threasholdMs"`
}

func (TrackStuckEvent) Event() EventType  { return EventTypeTrackStuck }
func (TrackStuckEvent) Op() OpType        { return OpTypeEvent }
func (e TrackStuckEvent) GuildID() string { return e.GID }
func (TrackStuckEvent) OpEvent()          {}

type WebsocketClosedEvent struct {
	GID      string `json:"guildId"`
	Code     int    `json:"code"`
	Reason   string `json:"reason"`
	ByRemote bool   `json:"byRemote"`
}

func (WebsocketClosedEvent) Event() EventType  { return EventTypeWebSocketClosed }
func (WebsocketClosedEvent) Op() OpType        { return OpTypeEvent }
func (e WebsocketClosedEvent) GuildID() string { return e.GID }
func (WebsocketClosedEvent) OpEvent()          {}

type UnknownEvent struct {
	event   EventType
	guildID string
	Data    []byte `json:"-"`
}

func (o *UnknownEvent) UnmarshalJSON(data []byte) error {
	var v struct {
		Event   EventType `json:"type"`
		GuildID string    `json:"guildId"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.event = v.Event
	o.guildID = v.GuildID
	o.Data = data
	return nil
}

func (e UnknownEvent) Event() EventType { return e.event }
func (UnknownEvent) Op() OpType         { return OpTypeEvent }
func (e UnknownEvent) GuildID() string  { return e.guildID }
func (UnknownEvent) OpEvent()           {}
