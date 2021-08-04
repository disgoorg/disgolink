package api

import "regexp"

// URLPattern is a general url pattern
var URLPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

// SpotifyURLPattern is a spotify url pattern with regions to get track/album/playlist
var SpotifyURLPattern = regexp.MustCompile("^(https?://)?(www\\.)?open\\.spotify\\.com/(user/(?P<user>[a-zA-Z0-9-_]+)/)?(?P<type>track|album|playlist)/(?P<identifier>[a-zA-Z0-9-_]+)?.+")
