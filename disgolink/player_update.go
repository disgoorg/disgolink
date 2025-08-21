package disgolink

import (
	"github.com/disgoorg/omit"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

func defaultPlayerUpdate() lavalink.PlayerUpdate {
	return lavalink.PlayerUpdate{}
}

type PlayerUpdateOpt func(update *lavalink.PlayerUpdate)

func playerUpdateApply(u *lavalink.PlayerUpdate, opts []PlayerUpdateOpt) {
	for _, opt := range opts {
		opt(u)
	}
}

func WithNoReplace(noReplace bool) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.NoReplace = noReplace
	}
}

func WithTrack(track lavalink.Track) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		WithEncodedTrack(track.Encoded)(update)
		WithTrackUserData(track.UserData)(update)
	}
}

func WithEncodedTrack(encodedTrack string) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		if update.Track == nil {
			update.Track = &lavalink.PlayerUpdateTrack{}
		}
		update.Track.Encoded = omit.NewPtr(encodedTrack)
	}
}

func WithNullTrack() PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		if update.Track == nil {
			update.Track = &lavalink.PlayerUpdateTrack{}
		}
		update.Track.Encoded = omit.NewNilPtr[string]()
	}
}

func WithTrackIdentifier(identifier string) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		if update.Track == nil {
			update.Track = &lavalink.PlayerUpdateTrack{}
		}
		update.Track.Identifier = &identifier
	}
}

func WithTrackUserData(userData any) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		if update.Track == nil {
			update.Track = &lavalink.PlayerUpdateTrack{}
		}
		update.Track.UserData = userData
	}
}

func WithPosition(position lavalink.Duration) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.Position = &position
	}
}

func WithEndTime(endTime lavalink.Duration) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.EndTime = &endTime
	}
}

func WithVolume(volume int) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.Volume = &volume
	}
}

func WithPaused(paused bool) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.Paused = &paused
	}
}

func WithVoice(voice lavalink.VoiceState) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.Voice = &voice
	}
}

func WithFilters(filters lavalink.Filters) PlayerUpdateOpt {
	return func(update *lavalink.PlayerUpdate) {
		update.Filters = &filters
	}
}
