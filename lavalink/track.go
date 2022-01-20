package lavalink

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrEmptyTrack = errors.New("track is empty")
var ErrEmptyTrackInfo = errors.New("trackinfo is empty")

type AudioTrack interface {
	AudioTrackInfo
	Track() string
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
}

func NewTrack(track string) AudioTrack {
	return &DefaultTrack{
		Base64Track: track,
	}
}

func NewTrackByInfo(trackInfo AudioTrackInfo) AudioTrack {
	return &DefaultTrack{
		AudioTrackInfo: trackInfo,
	}
}

type DefaultTrack struct {
	AudioTrackInfo `json:"info"`
	Base64Track    string `json:"track"`
}

func (t *DefaultTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		Base64Track *string           `json:"track"`
		TrackInfo   *defaultTrackInfo `json:"info"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	t.Base64Track = v.Base64Track
	t.AudioTrackInfo = v.TrackInfo
	return nil
}

func (t *DefaultTrack) Track() string {
	return t.Base64Track
}

type defaultTrackInfo struct {
	TrackIdentifier string        `json:"identifier"`
	TrackAuthor     string        `json:"author"`
	TrackLength     time.Duration `json:"length"`
	TrackIsStream   bool          `json:"isStream"`
	TrackTitle      string        `json:"title"`
	TrackURI        *string       `json:"uri"`
	TrackSourceName string        `json:"sourceName"`
	TrackPosition   time.Duration `json:"position"`
}

func (i *defaultTrackInfo) UnmarshalJSON(data []byte) error {
	type aliasDefaultTrackInfo defaultTrackInfo
	var v struct {
		TrackLength   int64 `json:"length"`
		TrackPosition int64 `json:"position"`
		aliasDefaultTrackInfo
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = defaultTrackInfo(v.aliasDefaultTrackInfo)
	i.TrackLength = time.Duration(v.TrackLength) * time.Millisecond
	i.TrackPosition = time.Duration(v.TrackPosition) * time.Millisecond
	return nil
}

func (i *defaultTrackInfo) MarshalJSON() ([]byte, error) {
	if i == nil {
		return nil, ErrEmptyTrackInfo
	}
	type aliasDefaultTrackInfo defaultTrackInfo
	return json.Marshal(struct {
		TrackLength   int64 `json:"length"`
		TrackPosition int64 `json:"position"`
		aliasDefaultTrackInfo
	}{
		TrackLength:           int64(i.TrackLength / time.Millisecond),
		TrackPosition:         int64(i.TrackPosition / time.Millisecond),
		aliasDefaultTrackInfo: aliasDefaultTrackInfo(*i),
	})
}

func (i *defaultTrackInfo) Identifier() string {
	return i.TrackIdentifier
}

func (i *defaultTrackInfo) Author() string {
	return i.TrackAuthor
}

func (i *defaultTrackInfo) Length() time.Duration {
	return i.TrackLength
}

func (i *defaultTrackInfo) IsStream() bool {
	return i.TrackIsStream
}

func (i *defaultTrackInfo) Title() string {
	return i.TrackTitle
}

func (i *defaultTrackInfo) URI() *string {
	return i.TrackURI
}

func (i *defaultTrackInfo) SourceName() string {
	return i.TrackSourceName
}

func (i *defaultTrackInfo) Position() time.Duration {
	return i.TrackPosition
}
