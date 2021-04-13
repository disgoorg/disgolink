package internal

import "github.com/DisgoOrg/disgolink/api"

type RestClientImpl struct {
}

func (c *RestClientImpl) GetYoutubeSearchResult(query string) {

}
func (c *RestClientImpl) GetYoutubeMusicSearchResult(query string) {

}
func (c *RestClientImpl) GetSoundcloudSearchResult(query string) {

}
func (c *RestClientImpl) LoadItemAsync(identifier string, audioLoaderResultHandler api.AudioLoaderResultHandler) {

}

func (c *RestClientImpl) LoadItem(identifier string) (*api.LoadResult, error) {
	return nil, nil
}
