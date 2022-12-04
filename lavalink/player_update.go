package lavalink

import "github.com/disgoorg/json"

func DefaultPlayerUpdate() *PlayerUpdate {
	return &PlayerUpdate{}
}

type PlayerUpdate struct {
	EncodedTrack *json.Nullable[string] `json:"encodedTrack,omitempty"`
	Identifier   *string                `json:"identifier,omitempty"`
	Position     *Duration              `json:"position,omitempty"`
	EndTime      *Duration              `json:"endTime,omitempty"`
	Volume       *int                   `json:"volume,omitempty"`
	Paused       *bool                  `json:"paused,omitempty"`
	Voice        *VoiceState            `json:"voice,omitempty"`
	Filters      *Filters               `json:"filters,omitempty"`
	NoReplace    bool                   `json:"-"`
}

type PlayerUpdateOpt func(update *PlayerUpdate)

func (u *PlayerUpdate) Apply(opts []PlayerUpdateOpt) {
	for _, opt := range opts {
		opt(u)
	}
}

func WithTrack(track Track) PlayerUpdateOpt {
	return WithEncodedTrack(track.Encoded)
}

func WithEncodedTrack(encodedTrack string) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.EncodedTrack = json.NewNullablePtr(encodedTrack)
	}
}

func WithNullTrack() PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.EncodedTrack = json.NullPtr[string]()
	}
}

func WithNoReplace(noReplace bool) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.NoReplace = noReplace
	}
}

func WithIdentifier(identifier string) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Identifier = &identifier
	}
}

func WithPosition(position Duration) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Position = &position
	}
}

func WithEndTime(endTime Duration) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.EndTime = &endTime
	}
}

func WithVolume(volume int) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Volume = &volume
	}
}

func WithPaused(paused bool) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Paused = &paused
	}
}

func WithVoice(voice VoiceState) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Voice = &voice
	}
}

func WithFilters(filters Filters) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Filters = &filters
	}
}
