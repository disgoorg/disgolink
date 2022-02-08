package lavalink

import (
	"context"
	"encoding/json"
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
	Plugins(ctx context.Context) ([]Plugin, error)
	LoadItem(ctx context.Context, identifier string) (*LoadResult, error)
	LoadItemHandler(ctx context.Context, identifier string, audioLoaderResultHandler AudioLoadResultHandler) error
	Version(ctx context.Context) (string, error)
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) Plugins(ctx context.Context) (plugins []Plugin, err error) {
	err = c.get(ctx, "/plugins", &plugins)
	if err != nil {
		return nil, err
	}
	return
}

func (c *restClientImpl) LoadItem(ctx context.Context, identifier string) (*LoadResult, error) {
	var result LoadResult
	err := c.get(ctx, "/loadtracks?identifier="+url.QueryEscape(identifier), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *restClientImpl) LoadItemHandler(ctx context.Context, identifier string, audioLoaderResultHandler AudioLoadResultHandler) error {
	result, err := c.LoadItem(ctx, identifier)
	if err != nil {
		return err
	}

	tracks, err := c.parseLoadResultTracks(result.Tracks)
	if err != nil {
		return err
	}

	switch result.LoadType {
	case LoadTypeTrackLoaded:
		audioLoaderResultHandler.TrackLoaded(tracks[0])

	case LoadTypePlaylistLoaded:
		audioLoaderResultHandler.PlaylistLoaded(NewAudioPlaylist(*result.PlaylistInfo, tracks))

	case LoadTypeSearchResult:
		audioLoaderResultHandler.SearchResultLoaded(tracks)

	case LoadTypeNoMatches:
		audioLoaderResultHandler.NoMatches()

	case LoadTypeLoadFailed:
		audioLoaderResultHandler.LoadFailed(*result.Exception)
	}
	return nil
}

func (c *restClientImpl) Version(ctx context.Context) (string, error) {
	path := "/version"
	rq, err := http.NewRequestWithContext(ctx, "GET", c.node.RestURL()+path, nil)
	if err != nil {
		return "", err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return "", err
	}
	defer rs.Body.Close()

	raw, err := ioutil.ReadAll(rs.Body)
	c.node.Lavalink().Logger().Debugf("response from %s, code %d, body: %s", c.node.RestURL()+path, rs.StatusCode, string(raw))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (c *restClientImpl) parseLoadResultTracks(loadResultTracks []LoadResultAudioTrack) ([]AudioTrack, error) {
	var tracks []AudioTrack
	for _, loadResultTrack := range loadResultTracks {
		track, err := c.node.Lavalink().DecodeTrack(loadResultTrack.Track)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (c *restClientImpl) get(ctx context.Context, path string, v interface{}) error {
	rq, err := http.NewRequestWithContext(ctx, "GET", c.node.RestURL()+path, nil)
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
	c.node.Lavalink().Logger().Debugf("response from %s, code %d, body: %s", c.node.RestURL()+path, rs.StatusCode, string(raw))
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, v)
	if err != nil {
		return err
	}
	return nil
}
