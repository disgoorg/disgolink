package api

type SearchType string

// search prefixes
const (
	SearchTypeYoutube      SearchType = "ytsearch:"
	SearchTypeYoutubeMusic SearchType = "ytmsearch:"
	SearchTypeSoundCloud   SearchType = "scsearch:"
)

func (t SearchType) Apply(searchString string) string {
	return string(t) + searchString
}

type RestClient interface {
	SearchItem(searchType SearchType, query string) ([]Track, *Exception)
	LoadItem(identifier string) *LoadResult
	LoadItemHandler(identifier string, audioLoaderResultHandler AudioLoaderResultHandler)
}
