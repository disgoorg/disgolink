package lavalink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

type RestClient interface {
	Version(ctx context.Context) (string, error)
	Info(ctx context.Context) (*Info, error)
	Stats(ctx context.Context) (*Stats, error)

	UpdateSession(ctx context.Context, sessionID string, sessionUpdate SessionUpdate) (*Session, error)

	Players(ctx context.Context, sessionID string) ([]Player, error)
	Player(ctx context.Context, sessionID string, guildID snowflake.ID) (*Player, error)
	UpdatePlayer(ctx context.Context, sessionID string, guildID snowflake.ID, playerUpdate PlayerUpdate) (*Player, error)
	DestroyPlayer(ctx context.Context, sessionID string, guildID snowflake.ID) error

	LoadTracks(ctx context.Context, identifier string) (*LoadResult, error)
	DecodeTrack(ctx context.Context, encodedTrack string) (*Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]Track, error)
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) Version(ctx context.Context) (string, error) {
	_, rawBody, err := c.do(ctx, http.MethodGet, "/version", nil)
	if err != nil {
		return "", err
	}
	return string(rawBody), nil
}

func (c *restClientImpl) Info(ctx context.Context) (info *Info, err error) {
	err = c.doJSON(ctx, http.MethodGet, "/v3/info", nil, &info)
	return
}

func (c *restClientImpl) Stats(ctx context.Context) (stats *Stats, err error) {
	err = c.doJSON(ctx, http.MethodGet, "/v3/stats", nil, &stats)
	return
}

func (c *restClientImpl) UpdateSession(ctx context.Context, sessionID string, sessionUpdate SessionUpdate) (session *Session, err error) {
	err = c.doJSON(ctx, http.MethodPost, "/v3/sessions/"+sessionID, sessionUpdate, &session)
	return
}

func (c *restClientImpl) Players(ctx context.Context, sessionID string) (players []Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, "/v3/sessions/"+sessionID+"/players", nil, &players)
	return
}

func (c *restClientImpl) Player(ctx context.Context, sessionID string, guildID snowflake.ID) (player *Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, "/v3/sessions/"+sessionID+"/players/"+guildID.String(), nil, &player)
	return
}

func (c *restClientImpl) UpdatePlayer(ctx context.Context, sessionID string, guildID snowflake.ID, playerUpdate PlayerUpdate) (player *Player, err error) {
	err = c.doJSON(ctx, http.MethodPost, "/v3/sessions/"+sessionID+"/players/"+guildID.String(), playerUpdate, &player)
	return
}

func (c *restClientImpl) DestroyPlayer(ctx context.Context, sessionID string, guildID snowflake.ID) error {
	_, _, err := c.do(ctx, http.MethodDelete, "/v3/sessions/"+sessionID+"/players/"+guildID.String(), nil)
	return err
}

func (c *restClientImpl) LoadTracks(ctx context.Context, identifier string) (*LoadResult, error) {
	var result LoadResult
	err := c.doJSON(ctx, http.MethodGet, "/v3/loadtracks?identifier="+url.QueryEscape(identifier), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *restClientImpl) DecodeTrack(ctx context.Context, encodedTrack string) (track *Track, err error) {
	err = c.doJSON(ctx, http.MethodGet, "/v3/decodetrack?track="+url.QueryEscape(encodedTrack), nil, &track)
	return
}

func (c *restClientImpl) DecodeTracks(ctx context.Context, encodedTracks []string) (tracks []Track, err error) {
	err = c.doJSON(ctx, http.MethodPost, "/v3/decodetracks", encodedTracks, &tracks)
	return
}

func (c *restClientImpl) do(ctx context.Context, method string, path string, rqBody io.Reader) (int, []byte, error) {
	rq, err := http.NewRequestWithContext(ctx, method, c.node.Config().RestURL()+path, rqBody)
	if err != nil {
		return 0, nil, err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return 0, nil, err
	}

	defer rs.Body.Close()
	rawBody, err := io.ReadAll(rs.Body)
	c.node.Lavalink().Logger().Tracef("response from %s, code %d, body: %s", c.node.Config().RestURL()+path, rs.StatusCode, string(rawBody))
	if err != nil {
		return rs.StatusCode, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if rs.StatusCode >= http.StatusBadRequest {
		var lavalinkErr Error
		if err = json.Unmarshal(rawBody, &lavalinkErr); err != nil {
			return rs.StatusCode, rawBody, fmt.Errorf("error while unmarshalling lavalink error: %w", err)
		}
		return rs.StatusCode, nil, lavalinkErr
	}

	return rs.StatusCode, rawBody, nil
}

func (c *restClientImpl) doJSON(ctx context.Context, method string, path string, rqBody any, rsBody any) error {
	var rqBodyReader io.Reader
	if rqBody != nil {
		var err error
		rawRqBody, err := json.Marshal(rqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		rqBodyReader = bytes.NewReader(rawRqBody)
	}
	statusCode, rawBody, err := c.do(ctx, method, path, rqBodyReader)
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		if err = json.Unmarshal(rawBody, rsBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}
	return json.Unmarshal(rawBody, rsBody)
}
