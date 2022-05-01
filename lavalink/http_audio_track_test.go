package lavalink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	httpURI          = "https://p.scdn.co/mp3-preview/ee121ca281c629bb4e99c33d877fe98fbb752289?cid=774b29d4f13844c495f206cafdad9c86"
	testHTTPTrackStr = "QAABFwIADVVua25vd24gdGl0bGUADlVua25vd24gYXJ0aXN0AAAAAAAAdSQAa2h0dHBzOi8vcC5zY2RuLmNvL21wMy1wcmV2aWV3L2VlMTIxY2EyODFjNjI5YmI0ZTk5YzMzZDg3N2ZlOThmYmI3NTIyODk/Y2lkPTc3NGIyOWQ0ZjEzODQ0YzQ5NWYyMDZjYWZkYWQ5Yzg2AAEAa2h0dHBzOi8vcC5zY2RuLmNvL21wMy1wcmV2aWV3L2VlMTIxY2EyODFjNjI5YmI0ZTk5YzMzZDg3N2ZlOThmYmI3NTIyODk/Y2lkPTc3NGIyOWQ0ZjEzODQ0YzQ5NWYyMDZjYWZkYWQ5Yzg2AARodHRwAANtcDMAAAAAAAAAAA=="
	testHTTPTrack    = &HTTPAudioTrack{
		AudioTrack: &BasicAudioTrack{
			AudioTrackInfo: AudioTrackInfo{
				Identifier: httpURI,
				Author:     "Unknown artist",
				Length:     29988,
				IsStream:   false,
				Title:      "Unknown title",
				URI:        &httpURI,
				SourceName: "http",
				Position:   0,
			},
		},
		ProbeInfo: "mp3",
	}
)

func TestDecodeHTTPTrack(t *testing.T) {
	track, err := DecodeString(testHTTPTrackStr, NewHTTPSourcePlugin().Decode)

	assert.NoError(t, err)
	assert.Equal(t, testHTTPTrack.Info(), track.Info())
}

func TestEncodeHTTPTrack(t *testing.T) {
	track, err := EncodeToString(testHTTPTrack, NewHTTPSourcePlugin().Encode)
	assert.NoError(t, err)
	assert.Equal(t, testHTTPTrackStr, track)
}

func TestEncodeDecodeHTTPTrack(t *testing.T) {
	audioTrack, err := DecodeString(testHTTPTrackStr, NewHTTPSourcePlugin().Decode)
	assert.NoError(t, err)
	assert.Equal(t, testHTTPTrack.Info(), audioTrack.Info())

	track, err := EncodeToString(audioTrack, NewHTTPSourcePlugin().Encode)
	assert.NoError(t, err)
	assert.Equal(t, testHTTPTrackStr, track)
}
