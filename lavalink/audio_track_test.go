package lavalink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	uri       = "https://www.youtube.com/watch?v=jdWhJcrrjQs"
	testTrack = &DefaultAudioTrack{
		AudioTrack: "QAAAfwIAFkFyY2hpdGVjdHMgLSAiQW5pbWFscyIAD0VwaXRhcGggUmVjb3JkcwAAAAAAA70IAAtqZFdoSmNycmpRcwABACtodHRwczovL3d3dy55b3V0dWJlLmNvbS93YXRjaD92PWpkV2hKY3JyalFzAAd5b3V0dWJlAAAAAAAAAAA=",
		AudioTrackInfo: &DefaultAudioTrackInfo{
			TrackIdentifier: "jdWhJcrrjQs",
			TrackAuthor:     "Epitaph Records",
			TrackLength:     245000 * time.Millisecond,
			TrackIsStream:   false,
			TrackTitle:      `Architects - "Animals"`,
			TrackURI:        &uri,
			TrackSourceName: "youtube",
			TrackPosition:   0,
		},
	}
)

func TestDecodeString(t *testing.T) {
	track, err := DecodeString(testTrack.Track(), nil)

	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info(), track.Info())
}

func TestEncodeTrackString(t *testing.T) {
	track, err := EncodeToString(testTrack, nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track(), track)
}

func TestEncodeDecodeString(t *testing.T) {
	audioTrack, err := DecodeString(testTrack.Track(), nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info(), audioTrack.Info())

	track, err := EncodeToString(audioTrack, nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track(), track)
}
