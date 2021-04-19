package api

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTrack = &Track{
	Track: "QAAAfwIAFkFyY2hpdGVjdHMgLSAiQW5pbWFscyIAD0VwaXRhcGggUmVjb3JkcwAAAAAAA70IAAtqZFdoSmNycmpRcwABACtodHRwczovL3d3dy55b3V0dWJlLmNvbS93YXRjaD92PWpkV2hKY3JyalFzAAd5b3V0dWJlAAAAAAADuiQ=",
	Info: &TrackInfo{
		Identifier: "jdWhJcrrjQs",
		IsSeekable: true,
		Author:     "Epitaph Records",
		Length:     245000,
		IsStream:   false,
		Position:   0,
		Title:      "Architects - \"Animals\"",
		URI:        "https://www.youtube.com/watch?v=jdWhJcrrjQs",
		SourceName: "youtube",
	},
}

func TestDecodeString(t *testing.T) {
	trackInfo, err := DecodeString(testTrack.Track)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Info, trackInfo)
}

func TestEncodeTrackInfo(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(testTrack.Track)
	log.Println("expected: ", data)
	track, err := EncodeToString(testTrack.Info)
	assert.NoError(t, err)
	assert.Equal(t, testTrack.Track, track)
}
