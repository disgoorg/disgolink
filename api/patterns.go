package api

import "regexp"

// URLPattern is a general url pattern
var URLPattern = regexp.MustCompile("^(https?|ftp|file)://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

// SpotifyURLPattern is a spotify url pattern with regions to get track/album/playlist
var SpotifyURLPattern = regexp.MustCompile("^(https?://)?(www\\.)?open\\.spotify\\.com/(track|album|playlist)/([a-zA-Z0-9-_]+)(\\?si=[a-zA-Z0-9-_]+)?")

var SpotifyURIPattern = regexp.MustCompile("^spotify:(track|album|playlist):([a-zA-Z0-9-_]+)")