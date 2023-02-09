package lavalink

type SearchType string

// search prefixes
const (
	SearchTypeYouTube      SearchType = "ytsearch"
	SearchTypeYouTubeMusic SearchType = "ytmsearch"
	SearchTypeSoundCloud   SearchType = "scsearch"
)

func (t SearchType) Apply(searchString string) string {
	return string(t) + ":" + searchString
}
