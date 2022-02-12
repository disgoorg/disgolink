package lavalink

import (
	"encoding/json"
)

const (
	trackInfoVersioned int32 = 1
	trackInfoVersion   int32 = 2
)

type AudioTrack interface {
	Info() AudioTrackInfo
	SetPosition(position Duration)
	Clone() AudioTrack
}

type AudioTrackInfo struct {
	Identifier string   `json:"identifier"`
	Author     string   `json:"author"`
	Length     Duration `json:"length"`
	IsStream   bool     `json:"isStream"`
	Title      string   `json:"title"`
	URI        *string  `json:"uri"`
	SourceName string   `json:"sourceName"`
	Position   Duration `json:"position"`
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

func (t *DefaultAudioTrack) SetPosition(position Duration) {
	t.AudioTrackInfo.Position = position
}

func (t *DefaultAudioTrack) Clone() AudioTrack {
	info := t.AudioTrackInfo
	info.Position = 0
	return &DefaultAudioTrack{
		AudioTrackInfo: info,
	}
}
