package lavalink

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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
	Version(ctx context.Context) (string, error)
	Plugins(ctx context.Context) ([]Plugin, error)
	LoadItem(ctx context.Context, identifier string) (*LoadResult, error)
	LoadItemHandler(ctx context.Context, identifier string, audioLoaderResultHandler AudioLoadResultHandler) error
	DecodeTrack(ctx context.Context, track string) (*AudioTrackInfo, error)
	DecodeTracks(ctx context.Context, tracks []string) ([]RestAudioTrack, error)
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) Version(ctx context.Context) (string, error) {
	rawBody, err := c.get(ctx, "/version")
	if err != nil {
		return "", err
	}
	return string(rawBody), nil
}

func (c *restClientImpl) Plugins(ctx context.Context) (plugins []Plugin, err error) {
	err = c.getJSON(ctx, "/plugins", &plugins)
	if err != nil {
		return nil, err
	}
	return
}

func (c *restClientImpl) LoadItem(ctx context.Context, identifier string) (*LoadResult, error) {
	var result LoadResult
	err := c.getJSON(ctx, "/loadtracks?identifier="+url.QueryEscape(identifier), &result)
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

	tracks, err := c.parseRestAudioTracks(result.Tracks)
	if err != nil {
		return err
	}

	switch result.LoadType {
	case LoadTypeTrackLoaded:
		audioLoaderResultHandler.TrackLoaded(tracks[0])

	case LoadTypePlaylistLoaded:
		audioLoaderResultHandler.PlaylistLoaded(NewAudioPlaylist(result.PlaylistInfo.Name, result.PlaylistInfo.SelectedTrack, tracks))

	case LoadTypeSearchResult:
		audioLoaderResultHandler.SearchResultLoaded(tracks)

	case LoadTypeNoMatches:
		audioLoaderResultHandler.NoMatches()

	case LoadTypeLoadFailed:
		audioLoaderResultHandler.LoadFailed(*result.Exception)
	}
	return nil
}

func (c *restClientImpl) DecodeTrack(ctx context.Context, track string) (*AudioTrackInfo, error) {
	var info AudioTrackInfo
	err := c.getJSON(ctx, "/decodetrack?track="+url.QueryEscape(track), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *restClientImpl) DecodeTracks(ctx context.Context, tracks []string) ([]RestAudioTrack, error) {
	var audioTracks []RestAudioTrack
	err := c.postJSON(ctx, "/decodetracks", tracks, &audioTracks)
	if err != nil {
		return nil, err
	}
	return audioTracks, nil
}

func (c *restClientImpl) parseRestAudioTracks(loadResultTracks []RestAudioTrack) ([]AudioTrack, error) {
	tracks := make([]AudioTrack, len(loadResultTracks))
	for i := range loadResultTracks {
		decodedTrack, err := c.node.Lavalink().DecodeTrack(loadResultTracks[i].Track)
		if err != nil {
			return nil, err
		}
		tracks[i] = decodedTrack
	}
	return tracks, nil
}

func (c *restClientImpl) do(ctx context.Context, method string, path string, body io.Reader) ([]byte, error) {
	rq, err := http.NewRequestWithContext(ctx, method, c.node.RestURL()+path, body)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return nil, errors.New(rs.Status)
	}
	rawBody, _ := io.ReadAll(rs.Body)
	c.node.Lavalink().Logger().Tracef("response from %s, code %d, body: %s", c.node.RestURL()+path, rs.StatusCode, string(rawBody))
	return rawBody, nil
}

func (c *restClientImpl) get(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *restClientImpl) getJSON(ctx context.Context, path string, v interface{}) error {
	rawBody, err := c.get(ctx, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(rawBody, v)
}

func (c *restClientImpl) postJSON(ctx context.Context, path string, b interface{}, v interface{}) error {
	rqBody, err := json.Marshal(b)
	if err != nil {
		return err
	}
	rsBody, err := c.do(ctx, http.MethodPost, path, bytes.NewReader(rqBody))
	if err != nil {
		return err
	}
	return json.Unmarshal(rsBody, v)
}
