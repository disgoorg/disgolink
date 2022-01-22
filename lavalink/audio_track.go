package lavalink

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrEmptyTrack     = errors.New("track is empty")
	ErrEmptyTrackInfo = errors.New("trackinfo is empty")
)

const (
	trackInfoVersioned int   = 1
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
	setPosition(position time.Duration)
}

func NewAudioTrack(track string) AudioTrack {
	return &DefaultAudioTrack{
		Base64Track: track,
	}
}

func NewAudioTrackByInfo(trackInfo AudioTrackInfo) AudioTrack {
	return &DefaultAudioTrack{
		AudioTrackInfo: trackInfo,
	}
}

type DefaultAudioTrack struct {
	Base64Track    string `json:"track"`
	AudioTrackInfo `json:"info"`
}

func (t *DefaultAudioTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		Base64Track string                 `json:"track"`
		TrackInfo   *DefaultAudioTrackInfo `json:"info"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	t.Base64Track = v.Base64Track
	t.AudioTrackInfo = v.TrackInfo
	return nil
}

func (t *DefaultAudioTrack) Track() string {
	if t.Base64Track == "" {
		if err := t.EncodeInfo(); err != nil {
			return ""
		}
	}
	return t.Base64Track
}

func (t *DefaultAudioTrack) Info() AudioTrackInfo {
	if t.AudioTrackInfo == nil {
		if err := t.DecodeInfo(); err != nil {
			return nil
		}
	}
	return t.AudioTrackInfo
}

func (t *DefaultAudioTrack) EncodeInfo() error {
	if t.AudioTrackInfo == nil {
		return ErrEmptyTrackInfo
	}

	var err error
	t.Base64Track, err = EncodeToString(t, nil)
	return err
}

func (t *DefaultAudioTrack) DecodeInfo() error {
	if t.Base64Track == "" {
		return ErrEmptyTrack
	}

	track, err := DecodeString(t.Base64Track, nil)
	if err != nil {
		return err
	}
	*t = *track
	return nil
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
	type defaultTrackInfo DefaultAudioTrackInfo

	return json.Marshal(struct {
		TrackLength   int64 `json:"length"`
		TrackPosition int64 `json:"position"`
		defaultTrackInfo
	}{
		TrackLength:      int64(i.TrackLength / time.Millisecond),
		TrackPosition:    int64(i.TrackPosition / time.Millisecond),
		defaultTrackInfo: defaultTrackInfo(i),
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

func (i *DefaultAudioTrackInfo) setPosition(position time.Duration) {
	i.TrackPosition = position
}
