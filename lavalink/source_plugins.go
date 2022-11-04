package lavalink

import (
	"errors"
	"io"

	"github.com/disgoorg/json"
)

func EncodeProbeInfo(probeInfo string, w io.Writer) error {
	return WriteString(w, probeInfo)
}

func DecodeProbeInfo(r io.Reader) (string, error) {
	return ReadString(r)
}

var _ SourcePlugin = (*HTTPSourcePlugin)(nil)

func NewHTTPSourcePlugin() *HTTPSourcePlugin {
	return &HTTPSourcePlugin{}
}

type HTTPSourcePlugin struct{}

func (*HTTPSourcePlugin) SourceName() string {
	return "http"
}

func (p *HTTPSourcePlugin) Encode(track AudioTrack, w io.Writer) error {
	httpTrack, ok := track.(*HTTPAudioTrack)
	if !ok {
		return errors.New("track is not a HTTPAudioTrack")
	}
	return EncodeProbeInfo(httpTrack.ProbeInfo, w)
}

func (p *HTTPSourcePlugin) Decode(info AudioTrackInfo, r io.Reader) (AudioTrack, error) {
	probeInfo, err := DecodeProbeInfo(r)
	if err != nil {
		return nil, err
	}
	return &HTTPAudioTrack{
		AudioTrack: &BasicAudioTrack{
			AudioTrackInfo: info,
		},
		ProbeInfo: probeInfo,
	}, nil
}

var _ AudioTrack = (*HTTPAudioTrack)(nil)

type HTTPAudioTrack struct {
	AudioTrack
	ProbeInfo string `json:"probeInfo"`
}

func (t *HTTPAudioTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		*BasicAudioTrack
		ProbeInfo string `json:"probeInfo"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.AudioTrack = v.BasicAudioTrack
	t.ProbeInfo = v.ProbeInfo
	return nil
}

func (t *HTTPAudioTrack) Clone() AudioTrack {
	return &HTTPAudioTrack{
		AudioTrack: t.AudioTrack.Clone(),
		ProbeInfo:  t.ProbeInfo,
	}
}

var _ SourcePlugin = (*LocalSourcePlugin)(nil)

func NewLocalSourcePlugin() *LocalSourcePlugin {
	return &LocalSourcePlugin{}
}

type LocalSourcePlugin struct{}

func (*LocalSourcePlugin) SourceName() string {
	return "local"
}

func (p *LocalSourcePlugin) Encode(track AudioTrack, w io.Writer) error {
	httpTrack, ok := track.(*LocalAudioTrack)
	if !ok {
		return errors.New("track is not a LocalAudioTrack")
	}
	return EncodeProbeInfo(httpTrack.ProbeInfo, w)
}

func (p *LocalSourcePlugin) Decode(info AudioTrackInfo, r io.Reader) (AudioTrack, error) {
	probeInfo, err := DecodeProbeInfo(r)
	if err != nil {
		return nil, err
	}
	return &LocalAudioTrack{
		AudioTrack: &BasicAudioTrack{
			AudioTrackInfo: info,
		},
		ProbeInfo: probeInfo,
	}, nil
}

var _ AudioTrack = (*LocalAudioTrack)(nil)

type LocalAudioTrack struct {
	AudioTrack
	ProbeInfo string `json:"probeInfo"`
}

func (t *LocalAudioTrack) UnmarshalJSON(data []byte) error {
	var v struct {
		*BasicAudioTrack
		ProbeInfo string `json:"probeInfo"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.AudioTrack = v.BasicAudioTrack
	t.ProbeInfo = v.ProbeInfo
	return nil
}

func (t *LocalAudioTrack) Clone() AudioTrack {
	return &LocalAudioTrack{
		AudioTrack: t.AudioTrack.Clone(),
		ProbeInfo:  t.ProbeInfo,
	}
}
