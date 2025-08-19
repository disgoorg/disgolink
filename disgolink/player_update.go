package disgolink

import (
	"github.com/disgoorg/omit"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

func defaultPlayerUpdate() PlayerUpdate {
	return PlayerUpdate{}
}

type PlayerUpdateTrack struct {
	Encoded    omit.Omit[*string] `json:"encoded,omitzero"`
	Identifier *string            `json:"identifier,omitzero"`
	UserData   any                `json:"userData,omitzero"`
}

type PlayerUpdate struct {
	Track     *PlayerUpdateTrack   `json:"track,omitzero"`
	Position  *lavalink.Duration   `json:"position,omitzero"`
	EndTime   *lavalink.Duration   `json:"endTime,omitzero"`
	Volume    *int                 `json:"volume,omitzero"`
	Paused    *bool                `json:"paused,omitzero"`
	Voice     *lavalink.VoiceState `json:"voice,omitzero"`
	Filters   *lavalink.Filters    `json:"filters,omitzero"`
	NoReplace bool                 `json:"-"`
}

type PlayerUpdateOpt func(update *PlayerUpdate)

func (u *PlayerUpdate) apply(opts []PlayerUpdateOpt) {
	for _, opt := range opts {
		opt(u)
	}
}

func WithNoReplace(noReplace bool) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.NoReplace = noReplace
	}
}

func WithTrack(track lavalink.Track) PlayerUpdateOpt {
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
		update.Track.Encoded = omit.NewPtr(encodedTrack)
	}
}

func WithNullTrack() PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		if update.Track == nil {
			update.Track = &PlayerUpdateTrack{}
		}
		update.Track.Encoded = omit.NewNilPtr[string]()
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

func WithPosition(position lavalink.Duration) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Position = &position
	}
}

func WithEndTime(endTime lavalink.Duration) PlayerUpdateOpt {
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

func WithVoice(voice lavalink.VoiceState) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Voice = &voice
	}
}

func WithFilters(filters lavalink.Filters) PlayerUpdateOpt {
	return func(update *PlayerUpdate) {
		update.Filters = &filters
	}
}
