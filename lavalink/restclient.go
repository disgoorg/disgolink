package lavalink

import (
	"github.com/DisgoOrg/disgo/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchType string

// search prefixes
const (
	SearchTypeYoutube      SearchType = "ytsearch"
	SearchTypeYoutubeMusic SearchType = "ytmsearch"
	SearchTypeSoundCloud   SearchType = "scsearch"
)

func (t SearchType) Apply(searchString string) string {
	return string(t) + ":" + searchString
}

type RestClient interface {
	SearchItem(searchType SearchType, query string) ([]Track, *Exception)
	LoadItem(identifier string) LoadResult
	LoadItemHandler(identifier string, audioLoaderResultHandler AudioLoaderResultHandler)
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) SearchItem(searchType SearchType, query string) ([]Track, *Exception) {
	result := c.LoadItem(searchType.Apply(query))
	if result.Exception != nil {
		return nil, result.Exception
	}

	return result.Tracks, nil
}

func (c *restClientImpl) LoadItem(identifier string) LoadResult {
	var result LoadResult
	err := c.get(c.node.RestURL()+"/loadtracks?identifier="+url.QueryEscape(identifier), &result)
	if err != nil {
		return LoadResult{LoadType: LoadTypeLoadFailed, Exception: NewExceptionFromErr(err)}
	}
	return result
}

func (c *restClientImpl) LoadItemHandler(identifier string, audioLoaderResultHandler AudioLoaderResultHandler) {
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

func (c *restClientImpl) get(url string, v interface{}) error {
	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)
	rq.Header.Set("Content-Type", "application/json")

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return err
	}

	defer rs.Body.Close()

	raw, err := ioutil.ReadAll(rs.Body)
	c.node.Lavalink().Logger().Debugf("response from %s, code %d, body: %s", url, rs.StatusCode, string(raw))
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, v)
	if err != nil {
		return err
	}
	return nil
}
