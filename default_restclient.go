package disgolink

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func newDefaultRestClient(node Node, httpClient *http.Client) RestClient {
	return &defaultRestClient{node: node, httpClient: httpClient}
}

type defaultRestClient struct {
	node       Node
	httpClient *http.Client
}

func (c *defaultRestClient) SearchItem(searchType SearchType, query string) ([]Track, *Exception) {
	result := c.LoadItem(searchType.Apply(query))
	if result.Exception != nil {
		return nil, result.Exception
	}

	return result.Tracks, nil
}

func (c *defaultRestClient) LoadItem(identifier string) LoadResult {
	var result LoadResult
	err := c.get(c.node.RestURL()+"/loadtracks?identifier="+url.QueryEscape(identifier), &result)
	if err != nil {
		return LoadResult{LoadType: LoadTypeLoadFailed, Exception: NewExceptionFromErr(err)}
	}
	return result
}

func (c *defaultRestClient) LoadItemHandler(identifier string, audioLoaderResultHandler AudioLoaderResultHandler) {
	result := c.LoadItem(identifier)

	switch result.LoadType {
	case LoadTypeTrackLoaded:
		audioLoaderResultHandler.TrackLoaded(result.Tracks[0])
	case LoadTypePlaylistLoaded:
		audioLoaderResultHandler.PlaylistLoaded(NewPlaylist(result))
	case LoadTypeSearchResult:
		audioLoaderResultHandler.SearchResultLoaded(result.Tracks)
	case LoadTypeNoMatches:
		audioLoaderResultHandler.NoMatches()
	case LoadTypeLoadFailed:
		audioLoaderResultHandler.LoadFailed(result.Exception)
	}
}

func (c *defaultRestClient) get(url string, v interface{}) error {
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
