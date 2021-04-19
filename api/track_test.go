package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var uri = "https://www.youtube.com/watch?v=jdWhJcrrjQs"

var testTrack = &Track{
	Track: "QAAAfwIAFkFyY2hpdGVjdHMgLSAiQW5pbWFscyIAD0VwaXRhcGggUmVjb3JkcwAAAAAAA70IAAtqZFdoSmNycmpRcwABACtodHRwczovL3d3dy55b3V0dWJlLmNvbS93YXRjaD92PWpkV2hKY3JyalFzAAd5b3V0dWJlAAAAAAADuiQ=",
	Info: &TrackInfo{
		Identifier: "jdWhJcrrjQs",
		IsSeekable: true,
		Author:     "Epitaph Records",
		Length:     245000,
		IsStream:   false,
		Position:   244260,
		Title:      "Architects - \"Animals\"",
		URI:        &uri,
		SourceName: "youtube",
	},
}

func TestDecodeString(t *testing.T) {
	trackInfo, err := DecodeString(testTrack.Track)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info, trackInfo)
}

func TestEncodeTrackString(t *testing.T) {
	track, err := EncodeToString(testTrack.Info)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track, track)
}

func TestEncodeDecodeString(t *testing.T) {
	trackInfo, err := DecodeString(testTrack.Track)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info, trackInfo)

	track, err := EncodeToString(trackInfo)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track, track)
}


