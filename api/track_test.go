package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var track = &Track{
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
	trackInfo, err := DecodeString(track.Track)
	assert.NoError(t, err)
	assert.Equal(t, track.Info, trackInfo)
}
