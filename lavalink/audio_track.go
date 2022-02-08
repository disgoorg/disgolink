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
	Info() AudioTrackInfo
	SetPosition(position time.Duration)
	Clone() AudioTrack
}

type AudioTrackInfo struct {
	Identifier string        `json:"identifier"`
	Author     string        `json:"author"`
	Length     time.Duration `json:"length"`
	IsStream   bool          `json:"isStream"`
	Title      string        `json:"title"`
	URI        *string       `json:"uri"`
	SourceName string        `json:"sourceName"`
	Position   time.Duration `json:"position"`
}

func (i *AudioTrackInfo) UnmarshalJSON(data []byte) error {
	type trackInfo AudioTrackInfo
	var v struct {
		Length   int64 `json:"length"`
		Position int64 `json:"position"`
		trackInfo
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = AudioTrackInfo(v.trackInfo)
	i.Length = time.Duration(v.Length) * time.Millisecond
	i.Position = time.Duration(v.Position) * time.Millisecond

	return nil
}

func (i AudioTrackInfo) MarshalJSON() ([]byte, error) {
	type audioTrackInfo AudioTrackInfo
	return json.Marshal(struct {
		Length   int64 `json:"length"`
		Position int64 `json:"position"`
		audioTrackInfo
	}{
		Length:         i.Length.Milliseconds(),
		Position:       i.Position.Milliseconds(),
		audioTrackInfo: audioTrackInfo(i),
	})
}

func NewAudioTrack(info AudioTrackInfo) AudioTrack {
	return &DefaultAudioTrack{
		AudioTrackInfo: info,
	}
}

type DefaultAudioTrack struct {
	AudioTrackInfo AudioTrackInfo `json:"info"`
}

func (t *DefaultAudioTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		AudioTrack     string         `json:"track"`
		AudioTrackInfo AudioTrackInfo `json:"info"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	t.AudioTrackInfo = v.AudioTrackInfo
	return nil
}

func (t *DefaultAudioTrack) Info() AudioTrackInfo {
	return t.AudioTrackInfo
}

func (t *DefaultAudioTrack) SetPosition(position time.Duration) {
	t.AudioTrackInfo.Position = position
}

func (t *DefaultAudioTrack) Clone() AudioTrack {
	info := t.AudioTrackInfo
	info.Position = 0
	return &DefaultAudioTrack{
		AudioTrackInfo: info,
	}
}
