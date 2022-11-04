package lavalink

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
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

type Info struct {
	Version        VersionInfo `json:"version"`
	BuildTime      time.Time   `json:"buildTime"`
	Git            Git         `json:"git"`
	JVM            string      `json:"jvm"`
	Lavaplayer     string      `json:"lavaplayer"`
	SourceManagers []string    `json:"sourceManagers"`
	Plugins        []Plugin    `json:"plugins"`
}

func (i Info) MarshalJSON() ([]byte, error) {
	type info Info
	return json.Marshal(struct {
		BuildTime int64 `json:"buildTime"`
		info
	}{
		BuildTime: i.BuildTime.UnixMilli(),
		info:      info(i),
	})
}

func (i *Info) UnmarshalJSON(data []byte) error {
	type info Info
	var v struct {
		BuildTime int64 `json:"buildTime"`
		info
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Info(v.info)
	i.BuildTime = time.UnixMilli(v.BuildTime)
	return nil
}

type VersionInfo struct {
	Semver     string `json:"semver"`
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"preRelease"`
}

type Git struct {
	Branch     string    `json:"branch"`
	Commit     string    `json:"commit"`
	CommitTime time.Time `json:"commitTime"`
}

func (g Git) MarshalJSON() ([]byte, error) {
	type git Git
	return json.Marshal(struct {
		CommitTime int64 `json:"commitTime"`
		git
	}{
		CommitTime: g.CommitTime.UnixMilli(),
		git:        git(g),
	})
}

func (g *Git) UnmarshalJSON(data []byte) error {
	type git Git
	var v struct {
		CommitTime int64 `json:"commitTime"`
		git
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*g = Git(v.git)
	g.CommitTime = time.UnixMilli(v.CommitTime)
	return nil
}

type PlayerUpdate struct {
	EncodedTrack *json.Nullable[string] `json:"encodedTrack,omitempty"`
	Identifier   *string                `json:"identifier,omitempty"`
	Position     *Duration              `json:"position,omitempty"`
	EndTime      *Duration              `json:"endTime,omitempty"`
	Volume       *int                   `json:"volume,omitempty"`
	Paused       *bool                  `json:"paused,omitempty"`
	Filters      *Filters               `json:"filters,omitempty"`
	Voice        *VoiceState            `json:"voice,omitempty"`
}

type VoiceState struct {
	Token     string `json:"token"`
	Endpoint  string `json:"endpoint"`
	SessionID string `json:"sessionId"`
	Connected bool   `json:"connected,omitempty"`
	Ping      int    `json:"ping,omitempty"`
}

type SessionUpdate struct {
	Key     *json.Nullable[*string] `json:"key,omitempty"`
	Timeout *int                    `json:"timeout,omitempty"`
}

type RestClient interface {
	Version(ctx context.Context) (*string, error)
	Info(ctx context.Context) (*Info, error)
	LoadItem(ctx context.Context, identifier string) (*LoadResult, error)
	LoadItemHandler(ctx context.Context, identifier string, audioLoaderResultHandler AudioLoadResultHandler) error
	DecodeTrack(ctx context.Context, encodedTrack string) (*Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]Track, error)

	GetPlayer(ctx context.Context, guildID snowflake.ID) (Player, error)
	UpdatePlayer(ctx context.Context, guildID snowflake.ID, update PlayerUpdate, noReplace bool) error
	DestroyPlayer(ctx context.Context, guildID snowflake.ID) error
	UpdateSession(ctx context.Context, sessionUpdate SessionUpdate) error
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) Version(ctx context.Context) (*string, error) {
	rawBody, err := c.get(ctx, "/version")
	if err != nil {
		return nil, err
	}
	version := string(rawBody)
	return &version, nil
}

func (c *restClientImpl) Info(ctx context.Context) (info *Info, err error) {
	err = c.getJSON(ctx, "/v3/info", &info)
	if err != nil {
		return nil, err
	}
	return
}

func (c *restClientImpl) LoadItem(ctx context.Context, identifier string) (*LoadResult, error) {
	var result LoadResult
	err := c.getJSON(ctx, "/v3/loadtracks?identifier="+url.QueryEscape(identifier), &result)
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

func (c *restClientImpl) DecodeTrack(ctx context.Context, encodedTrack string) (info *Track, err error) {
	err = c.getJSON(ctx, "/v3/decodetrack?encodedTrack="+url.QueryEscape(encodedTrack), &info)
	return
}

func (c *restClientImpl) DecodeTracks(ctx context.Context, encodedTracks []string) (audioTracks []Track, err error) {
	err = c.postJSON(ctx, "/v3/decodetracks", encodedTracks, &audioTracks)
	return
}

func (c *restClientImpl) GetPlayer(ctx context.Context, guildID snowflake.ID) (player Player, err error) {
	var defaultPlayer DefaultPlayer
	err = c.getJSON(ctx, fmt.Sprintf("/players/%d", guildID), &defaultPlayer)
	if err == nil {
		player = &defaultPlayer
	}
	return
}

func (c *restClientImpl) UpdatePlayer(ctx context.Context, guildID snowflake.ID, update PlayerUpdate, noReplace bool) error {
	return c.patchJSON(ctx, fmt.Sprintf("/v3/sessions/%s/players/%d?noReplace=%t", c.node.SessionID(), guildID, noReplace), update, nil)
}

func (c *restClientImpl) DestroyPlayer(ctx context.Context, guildID snowflake.ID) error {
	return c.delete(ctx, fmt.Sprintf("/v3/sessions/%s/players/%d", c.node.SessionID(), guildID))
}

func (c *restClientImpl) UpdateSession(ctx context.Context, sessionUpdate SessionUpdate) error {
	return c.patchJSON(ctx, fmt.Sprintf("/v3/sessions/%s", c.node.SessionID()), sessionUpdate, nil)
}

func (c *restClientImpl) parseRestAudioTracks(loadResultTracks []Track) ([]AudioTrack, error) {
	tracks := make([]AudioTrack, len(loadResultTracks))
	for i := range loadResultTracks {
		decodedTrack, err := c.node.Lavalink().DecodeTrack(loadResultTracks[i].Encoded)
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
	rq.Header.Set("Session-Id", c.node.SessionID())
	rq.Header.Set("Content-Type", "application/json")

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK && rs.StatusCode != http.StatusNoContent {
		return nil, errors.New(rs.Status)
	}
	rawBody, _ := io.ReadAll(rs.Body)
	//c.node.Lavalink().Logger().Tracef("response from %s, code %d, body: %s", c.node.RestURL()+path, rs.StatusCode, string(rawBody))
	return rawBody, nil
}

func (c *restClientImpl) get(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *restClientImpl) delete(ctx context.Context, path string) error {
	_, err := c.do(ctx, http.MethodDelete, path, nil)
	return err
}

func (c *restClientImpl) getJSON(ctx context.Context, path string, v any) error {
	rsBody, err := c.get(ctx, path)
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(rsBody, v)
	}
	return nil
}

func (c *restClientImpl) postJSON(ctx context.Context, path string, b any, v any) error {
	rqBody, err := json.Marshal(b)
	if err != nil {
		return err
	}
	rsBody, err := c.do(ctx, http.MethodPost, path, bytes.NewReader(rqBody))
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(rsBody, v)
	}
	return nil
}

func (c *restClientImpl) patchJSON(ctx context.Context, path string, b any, v any) error {
	rqBody, err := json.Marshal(b)
	if err != nil {
		return err
	}

	c.node.Lavalink().Logger().Tracef("request to %s, body: %s", c.node.RestURL()+path, string(rqBody))

	rsBody, err := c.do(ctx, http.MethodPatch, path, bytes.NewReader(rqBody))
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(rsBody, v)
	}
	return nil
}
