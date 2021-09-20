package disgolink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var uri = "https://www.youtube.com/watch?v=jdWhJcrrjQs"

var encodedTrack = "QAAAfwIAFkFyY2hpdGVjdHMgLSAiQW5pbWFscyIAD0VwaXRhcGggUmVjb3JkcwAAAAAAA70IAAtqZFdoSmNycmpRcwABACtodHRwczovL3d3dy55b3V0dWJlLmNvbS93YXRjaD92PWpkV2hKY3JyalFzAAd5b3V0dWJlAAAAAAADuiQ="
var testTrack = &DefaultTrack{
	Base64Track: &encodedTrack,
	TrackInfo: &DefaultTrackInfo{
		TrackIdentifier: "jdWhJcrrjQs",
		TrackIsSeekable: true,
		TrackAuthor:     "Epitaph Records",
		TrackLength:     245000,
		TrackIsStream:   false,
		TrackPosition:   244260,
		TrackTitle:      "Architects - \"Animals\"",
		TrackURI:        &uri,
		TrackSourceName: "youtube",
	},
}

func TestDecodeString(t *testing.T) {
	trackInfo, err := DecodeString(*testTrack.Track())
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info, trackInfo)
}

func TestEncodeTrackString(t *testing.T) {
	track, err := EncodeToString(testTrack.Info())
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track, track)
}

func TestEncodeDecodeString(t *testing.T) {
	trackInfo, err := DecodeString(*testTrack.Track())
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info, trackInfo)

	track, err := EncodeToString(trackInfo)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track, track)
}
