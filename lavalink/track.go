package lavalink

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrEmptyTrack = errors.New("track is empty")
var ErrEmptyTrackInfo = errors.New("trackinfo is empty")

type Track interface {
	Track() string
	TrackInfo
}

type TrackInfo interface {
	Identifier() string
	Author() string
	Length() time.Duration
	IsStream() bool
	Title() string
	URI() *string
	SourceName() string
}

func NewTrack(track string) Track {
	return &DefaultTrack{
		Base64Track: track,
	}
}

func NewTrackByInfo(trackInfo TrackInfo) Track {
	return &DefaultTrack{
		TrackInfo: trackInfo,
	}
}

type DefaultTrack struct {
	Base64Track string `json:"track"`
	TrackInfo   `json:"info"`
}

func (t *DefaultTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		Base64Track string           `json:"track"`
		TrackInfo   DefaultTrackInfo `json:"info"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	t.Base64Track = v.Base64Track
	t.TrackInfo = v.TrackInfo
	return nil
}

func (t *DefaultTrack) Track() string {
	if t.Base64Track == "" {
		if err := t.EncodeInfo(); err != nil {
			return ""
		}
	}
	return t.Base64Track
}

func (t *DefaultTrack) Info() TrackInfo {
	if t.TrackInfo == nil {
		if err := t.DecodeInfo(); err != nil {
			return nil
		}
	}
	return t.TrackInfo
}

func (t *DefaultTrack) EncodeInfo() error {
	if t.TrackInfo == nil {
		return ErrEmptyTrackInfo
	}
	var err error
	t.Base64Track, err = EncodeToString(t.TrackInfo)
	return err
}

func (t *DefaultTrack) DecodeInfo() error {
	if t.Base64Track == "" {
		return ErrEmptyTrack
	}
	var err error
	t.TrackInfo, err = DecodeString(t.Base64Track)
	return err
}

type DefaultTrackInfo struct {
	TrackIdentifier string        `json:"identifier"`
	TrackAuthor     string        `json:"author"`
	TrackLength     time.Duration `json:"length"`
	TrackIsStream   bool          `json:"isStream"`
	TrackTitle      string        `json:"title"`
	TrackURI        *string       `json:"uri"`
	TrackSourceName string        `json:"sourceName"`
}

func (i *DefaultTrackInfo) UnmarshalJSON(data []byte) error {
	type defaultTrackInfo DefaultTrackInfo
	var v struct {
		TrackLength int64 `json:"length"`
		defaultTrackInfo
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = DefaultTrackInfo(v.defaultTrackInfo)
	i.TrackLength = time.Duration(v.TrackLength) * time.Millisecond
	return nil
}

func (i DefaultTrackInfo) MarshalJSON() ([]byte, error) {
	type defaultTrackInfo DefaultTrackInfo
	return json.Marshal(struct {
		TrackLength int64 `json:"length"`
		defaultTrackInfo
	}{
		TrackLength:      int64(i.TrackLength / time.Millisecond),
		defaultTrackInfo: defaultTrackInfo(i),
	})
}

func (i DefaultTrackInfo) Identifier() string {
	return i.TrackIdentifier
}

func (i DefaultTrackInfo) Author() string {
	return i.TrackAuthor
}

func (i DefaultTrackInfo) Length() time.Duration {
	return i.TrackLength
}

func (i DefaultTrackInfo) IsStream() bool {
	return i.TrackIsStream
}

func (i DefaultTrackInfo) Title() string {
	return i.TrackTitle
}

func (i DefaultTrackInfo) URI() *string {
	return i.TrackURI
}

func (i DefaultTrackInfo) SourceName() string {
	return i.TrackSourceName
}
