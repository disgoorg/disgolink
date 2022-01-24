package lavalink

import (
	"encoding/json"
	"time"
)

const (
	trackInfoVersioned int32 = 1
	trackInfoVersion   int32 = 2
)

type AudioTrack interface {
	Track() string
	Info() AudioTrackInfo
}

type AudioTrackInfo interface {
	Identifier() string
	Author() string
	Length() time.Duration
	IsStream() bool
	Title() string
	URI() *string
	SourceName() string
	Position() time.Duration
	SetPosition(position time.Duration)
}

func NewAudioTrack(track string, info AudioTrackInfo) AudioTrack {
	return &DefaultAudioTrack{
		AudioTrack:     track,
		AudioTrackInfo: info,
	}
}

type DefaultAudioTrack struct {
	AudioTrack     string         `json:"track"`
	AudioTrackInfo AudioTrackInfo `json:"info"`
}

func (t *DefaultAudioTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		AudioTrack     string                `json:"track"`
		AudioTrackInfo DefaultAudioTrackInfo `json:"info"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	t.AudioTrack = v.AudioTrack
	t.AudioTrackInfo = &v.AudioTrackInfo
	return nil
}

func (t DefaultAudioTrack) Track() string {
	return t.AudioTrack
}

func (t DefaultAudioTrack) Info() AudioTrackInfo {
	return t.AudioTrackInfo
}

type DefaultAudioTrackInfo struct {
	TrackIdentifier string        `json:"identifier"`
	TrackAuthor     string        `json:"author"`
	TrackLength     time.Duration `json:"length"`
	TrackIsStream   bool          `json:"isStream"`
	TrackTitle      string        `json:"title"`
	TrackURI        *string       `json:"uri"`
	TrackSourceName string        `json:"sourceName"`
	TrackPosition   time.Duration `json:"position"`
}

func (i *DefaultAudioTrackInfo) UnmarshalJSON(data []byte) error {
	type defaultTrackInfo DefaultAudioTrackInfo
	var v struct {
		TrackLength   int64 `json:"length"`
		TrackPosition int64 `json:"position"`
		defaultTrackInfo
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = DefaultAudioTrackInfo(v.defaultTrackInfo)
	i.TrackLength = time.Duration(v.TrackLength) * time.Millisecond
	i.TrackPosition = time.Duration(v.TrackPosition) * time.Millisecond

	return nil
}

func (i DefaultAudioTrackInfo) MarshalJSON() ([]byte, error) {
	type defaultAudioTrackInfo DefaultAudioTrackInfo
	return json.Marshal(struct {
		TrackLength   int64 `json:"length"`
		TrackPosition int64 `json:"position"`
		defaultAudioTrackInfo
	}{
		TrackLength:           i.TrackLength.Milliseconds(),
		TrackPosition:         i.TrackPosition.Milliseconds(),
		defaultAudioTrackInfo: defaultAudioTrackInfo(i),
	})
}

func (i DefaultAudioTrackInfo) Identifier() string {
	return i.TrackIdentifier
}

func (i DefaultAudioTrackInfo) Author() string {
	return i.TrackAuthor
}

func (i DefaultAudioTrackInfo) Length() time.Duration {
	return i.TrackLength
}

func (i DefaultAudioTrackInfo) IsStream() bool {
	return i.TrackIsStream
}

func (i DefaultAudioTrackInfo) Title() string {
	return i.TrackTitle
}

func (i DefaultAudioTrackInfo) URI() *string {
	return i.TrackURI
}

func (i DefaultAudioTrackInfo) SourceName() string {
	return i.TrackSourceName
}

func (i DefaultAudioTrackInfo) Position() time.Duration {
	return i.TrackPosition
}

func (i *DefaultAudioTrackInfo) SetPosition(position time.Duration) {
	i.TrackPosition = position
}
