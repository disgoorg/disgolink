package lavalink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	uri          = "https://www.youtube.com/watch?v=jdWhJcrrjQs"
	testTrackStr = "QAAAfwIAFkFyY2hpdGVjdHMgLSAiQW5pbWFscyIAD0VwaXRhcGggUmVjb3JkcwAAAAAAA70IAAtqZFdoSmNycmpRcwABACtodHRwczovL3d3dy55b3V0dWJlLmNvbS93YXRjaD92PWpkV2hKY3JyalFzAAd5b3V0dWJlAAAAAAAAAAA="
	testTrack    = &BasicAudioTrack{
		AudioTrackInfo: AudioTrackInfo{
			Identifier: "jdWhJcrrjQs",
			Author:     "Epitaph Records",
			Length:     245000,
			IsStream:   false,
			Title:      `Architects - "Animals"`,
			URI:        &uri,
			SourceName: "youtube",
			Position:   0,
		},
	}
)

func TestDecodeString(t *testing.T) {
	track, err := DecodeString(testTrackStr, nil)

	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info(), track.Info())
}

func TestEncodeTrackString(t *testing.T) {
	track, err := EncodeToString(testTrack, nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrackStr, track)
}

func TestEncodeDecodeString(t *testing.T) {
	audioTrack, err := DecodeString(testTrackStr, nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info(), audioTrack.Info())

	track, err := EncodeToString(audioTrack, nil)
	assert.NoError(t, err)
	assert.Equal(t, testTrackStr, track)
}
