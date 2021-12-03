package lavalink

import "github.com/DisgoOrg/disgo/json"

type DefaultTrack struct {
	Base64Track *string   `json:"track"`
	TrackInfo   TrackInfo `json:"info"`
}

func (t *DefaultTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		Base64Track *string           `json:"track"`
		TrackInfo   *DefaultTrackInfo `json:"info"`
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
	if t.Base64Track == nil {
		if err := t.EncodeInfo(); err != nil {
			return ""
		}
	}
	return *t.Base64Track
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
	track, err := EncodeToString(t.TrackInfo)
	if err != nil {
		return err
	}
	t.Base64Track = &track
	return nil
}

func (t *DefaultTrack) DecodeInfo() error {
	if t.Base64Track == nil {
		return ErrEmptyTrack
	}
	var err error
	t.TrackInfo, err = DecodeString(*t.Base64Track)
	if err != nil {
		return err
	}
	return nil
}

type DefaultTrackInfo struct {
	TrackIdentifier string  `json:"identifier"`
	TrackIsSeekable bool    `json:"isSeekable"`
	TrackAuthor     string  `json:"author"`
	TrackLength     int     `json:"length"`
	TrackIsStream   bool    `json:"isStream"`
	TrackPosition   int     `json:"position"`
	TrackTitle      string  `json:"title"`
	TrackURI        *string `json:"uri"`
	TrackSourceName string  `json:"sourceName"`
}

func (i *DefaultTrackInfo) Identifier() string {
	return i.TrackIdentifier
}

func (i *DefaultTrackInfo) IsSeekable() bool {
	return i.TrackIsSeekable
}

func (i *DefaultTrackInfo) Author() string {
	return i.TrackAuthor
}

func (i *DefaultTrackInfo) Length() int {
	return i.TrackLength
}

func (i *DefaultTrackInfo) IsStream() bool {
	return i.TrackIsStream
}

func (i *DefaultTrackInfo) Position() int {
	return i.TrackPosition
}

func (i *DefaultTrackInfo) Title() string {
	return i.TrackTitle
}

func (i *DefaultTrackInfo) URI() *string {
	return i.TrackURI
}

func (i *DefaultTrackInfo) SourceName() string {
	return i.TrackSourceName
}
