package lavalink

import "github.com/disgoorg/json"

func DefaultPlayerUpdate() *PlayerUpdate {
	return &PlayerUpdate{}
}

type PlayerUpdateTrack struct {
	Encoded    *json.Nullable[string] `json:"encoded,omitempty"`
	Identifier *string                `json:"identifier,omitempty"`
	UserData   any                    `json:"userData,omitempty"`
}

type PlayerUpdate struct {
	Track     *PlayerUpdateTrack `json:"track,omitempty"`
	Position  *Duration          `json:"position,omitempty"`
	EndTime   *Duration          `json:"endTime,omitempty"`
	Volume    *int               `json:"volume,omitempty"`
	Paused    *bool              `json:"paused,omitempty"`
	Voice     *VoiceState        `json:"voice,omitempty"`
	Filters   *Filters           `json:"filters,omitempty"`
	NoReplace bool               `json:"-"`
}

type PlayerUpdateOpt func(update *PlayerUpdate)

func (u *PlayerUpdate) Apply(opts []PlayerUpdateOpt) {
	for _, opt := range opts {
		opt(u)
	}
}

func WithNoReplace(noReplace bool) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.NoReplace = noReplace
	}
}

func WithTrack(track Track) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		WithEncodedTrack(track.Encoded)(update)
		WithTrackUserData(track.UserData)(update)
	}
}

func WithEncodedTrack(encodedTrack string) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		if update.Track == nil {
			update.Track = &PlayerUpdateTrack{}
		}
		update.Track.Encoded = json.NewNullablePtr(encodedTrack)
	}
}

func WithNullTrack() PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		if update.Track == nil {
			update.Track = &PlayerUpdateTrack{}
		}
		update.Track.Encoded = json.NullPtr[string]()
	}
}

func WithTrackIdentifier(identifier string) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		if update.Track == nil {
			update.Track = &PlayerUpdateTrack{}
		}
		update.Track.Identifier = &identifier
	}
}

func WithTrackUserData(userData any) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		if update.Track == nil {
			update.Track = &PlayerUpdateTrack{}
		}
		update.Track.UserData = userData
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
