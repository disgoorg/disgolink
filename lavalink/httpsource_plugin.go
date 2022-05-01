package lavalink

import (
	"errors"
	"io"
)

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
	return WriteString(w, httpTrack.ProbeInfo)
}

func (p *HTTPSourcePlugin) Decode(info AudioTrackInfo, r io.Reader) (AudioTrack, error) {
	probeInfo, err := ReadString(r)
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

func (t *HTTPAudioTrack) Clone() AudioTrack {
	return &HTTPAudioTrack{
		AudioTrack: t.AudioTrack.Clone(),
		ProbeInfo:  t.ProbeInfo,
	}
}
