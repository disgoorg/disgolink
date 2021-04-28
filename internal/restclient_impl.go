package internal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/DisgoOrg/disgolink/api"
)

func newRestClientImpl(node api.Node, httpClient *http.Client) api.RestClient {
	return &RestClientImpl{node: node, httpClient: httpClient}
}

type RestClientImpl struct {
	node       api.Node
	httpClient *http.Client
}

func (c *RestClientImpl) SearchItem(searchType api.SearchType, query string) ([]api.Track, *api.Exception) {
	result, err := c.LoadItem(string(searchType) + query)
	if err != nil {
		return nil, api.NewExceptionFromErr(err)
	}
	if result.Exception != nil {
		return nil, result.Exception
	}

	return api.DefaultTracksToTracks(result.Tracks), nil
}
func (c *RestClientImpl) LoadItemAsync(identifier string, audioLoaderResultHandler api.AudioLoaderResultHandler) {
	go func() {
		result, err := c.LoadItem(identifier)
		if err != nil {
			audioLoaderResultHandler.LoadFailed(api.NewExceptionFromErr(err))
			return
		}

		switch result.LoadType {
		case api.LoadTypeTrackLoaded:
			audioLoaderResultHandler.TrackLoaded(result.Tracks[0])
		case api.LoadTypePlaylistLoaded:
			audioLoaderResultHandler.PlaylistLoaded(api.NewPlaylist(result))
		case api.LoadTypeSearchResult:
			audioLoaderResultHandler.SearchResultLoaded(api.DefaultTracksToTracks(result.Tracks))
		case api.LoadTypeNoMatches:
			audioLoaderResultHandler.NoMatches()
		case api.LoadTypeLoadFailed:
			audioLoaderResultHandler.LoadFailed(result.Exception)
		}
	}()
}

func (c *RestClientImpl) LoadItem(identifier string) (*api.LoadResult, error) {
	var result *api.LoadResult
	err := c.get(c.node.RestURL()+"/loadtracks?identifier="+url.QueryEscape(identifier), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RestClientImpl) get(url string, v interface{}) error {
	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", c.node.Options().Password)
	rq.Header.Set("Content-Type", "application/json")

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return err
	}

	defer rs.Body.Close()

	raw, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, v)
	if err != nil {
		return err
	}
	return nil
}
