package api

// search prefixes
const (
	YoutubeSearchPrefix      = "ytsearch:"
	YoutubeMusicSearchPrefix = "ytmsearch"
	SoundCloudSearchPrefix   = "scsearch:"
)

type RestClient interface {
	GetYoutubeSearchResult(query string)
	GetYoutubeMusicSearchResult(query string)
	GetSoundcloudSearchResult(query string)
	LoadItemAsync(identifier string, audioLoaderResultHandler AudioLoaderResultHandler)
	LoadItem(identifier string, audioLoaderResultHandler AudioLoaderResultHandler)
}
