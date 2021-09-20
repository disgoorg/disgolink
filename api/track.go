package api

import "errors"

var ErrEmptyTrack = errors.New("track is empty")
var ErrEmptyTrackInfo = errors.New("trackinfo is empty")

type Track interface {
	Track() string
	Info() TrackInfo
}

type TrackInfo interface {
	Identifier() string
	IsSeekable() bool
	Author() string
	Length() int
	IsStream() bool
	Position() int
	Title() string
	URI() *string
	SourceName() string
}
