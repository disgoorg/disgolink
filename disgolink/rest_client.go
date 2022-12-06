package disgolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Endpoint string

func (e Endpoint) Format(a ...any) string {
	return fmt.Sprintf(string(e), a...)
}

var (
	EndpointVersion Endpoint = "/version"
	EndpointInfo    Endpoint = "/v3/info"
	EndpointStats   Endpoint = "/v3/stats"

	EndpointUpdateSession Endpoint = "/v3/sessions/%s"
	EndpointPlayers       Endpoint = "/v3/sessions/%s/players"
	EndpointPlayer        Endpoint = "/v3/sessions/%s/players/%s"
	EndpointUpdatePlayer  Endpoint = "/v3/sessions/%s/players/%s?noReplace=%t"
	EndpointDestroyPlayer Endpoint = "/v3/sessions/%s/players/%s"

	EndpointLoadTracks   Endpoint = "/v3/loadtracks?identifier=%s"
	EndpointDecodeTrack  Endpoint = "/v3/decodetrack?track=%s"
	EndpointDecodeTracks Endpoint = "/v3/decodetracks"

	EndpointWebSocket Endpoint = "/v3/websocket"
)

type RestClient interface {
	Version(ctx context.Context) (string, error)
	Info(ctx context.Context) (*lavalink.Info, error)
	Stats(ctx context.Context) (*lavalink.Stats, error)

	UpdateSession(ctx context.Context, sessionID string, sessionUpdate lavalink.SessionUpdate) (*lavalink.Session, error)

	Players(ctx context.Context, sessionID string) ([]lavalink.Player, error)
	Player(ctx context.Context, sessionID string, guildID snowflake.ID) (*lavalink.Player, error)
	UpdatePlayer(ctx context.Context, sessionID string, guildID snowflake.ID, playerUpdate lavalink.PlayerUpdate) (*lavalink.Player, error)
	DestroyPlayer(ctx context.Context, sessionID string, guildID snowflake.ID) error

	LoadTracks(ctx context.Context, identifier string) (*lavalink.LoadResult, error)
	DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error)
}

func newRestClientImpl(node Node, httpClient *http.Client) RestClient {
	return &restClientImpl{node: node, httpClient: httpClient}
}

type restClientImpl struct {
	node       Node
	httpClient *http.Client
}

func (c *restClientImpl) Version(ctx context.Context) (string, error) {
	_, rawBody, err := c.do(ctx, http.MethodGet, string(EndpointVersion), nil)
	if err != nil {
		return "", err
	}
	return string(rawBody), nil
}

func (c *restClientImpl) Info(ctx context.Context) (info *lavalink.Info, err error) {
	err = c.doJSON(ctx, http.MethodGet, string(EndpointInfo), nil, &info)
	return
}

func (c *restClientImpl) Stats(ctx context.Context) (stats *lavalink.Stats, err error) {
	err = c.doJSON(ctx, http.MethodGet, string(EndpointStats), nil, &stats)
	return
}

func (c *restClientImpl) UpdateSession(ctx context.Context, sessionID string, sessionUpdate lavalink.SessionUpdate) (session *lavalink.Session, err error) {
	err = c.doJSON(ctx, http.MethodPost, EndpointUpdateSession.Format(sessionID), sessionUpdate, &session)
	return
}

func (c *restClientImpl) Players(ctx context.Context, sessionID string) (players []lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointPlayers.Format(sessionID), nil, &players)
	return
}

func (c *restClientImpl) Player(ctx context.Context, sessionID string, guildID snowflake.ID) (player *lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointPlayer.Format(sessionID, guildID), nil, &player)
	return
}

func (c *restClientImpl) UpdatePlayer(ctx context.Context, sessionID string, guildID snowflake.ID, playerUpdate lavalink.PlayerUpdate) (player *lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodPatch, EndpointUpdatePlayer.Format(sessionID, guildID, playerUpdate.NoReplace), playerUpdate, &player)
	return
}

func (c *restClientImpl) DestroyPlayer(ctx context.Context, sessionID string, guildID snowflake.ID) error {
	_, _, err := c.do(ctx, http.MethodDelete, EndpointDestroyPlayer.Format(sessionID, guildID), nil)
	return err
}

func (c *restClientImpl) LoadTracks(ctx context.Context, identifier string) (result *lavalink.LoadResult, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointLoadTracks.Format(url.QueryEscape(identifier)), nil, &result)
	return
}

func (c *restClientImpl) DecodeTrack(ctx context.Context, encodedTrack string) (track *lavalink.Track, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointDecodeTrack.Format(url.QueryEscape(encodedTrack)), nil, &track)
	return
}

func (c *restClientImpl) DecodeTracks(ctx context.Context, encodedTracks []string) (tracks []lavalink.Track, err error) {
	err = c.doJSON(ctx, http.MethodPost, string(EndpointDecodeTracks), encodedTracks, &tracks)
	return
}

func (c *restClientImpl) do(ctx context.Context, method string, path string, rqBody io.Reader) (int, []byte, error) {
	rq, err := http.NewRequestWithContext(ctx, method, c.node.Config().RestURL()+path, rqBody)
	if err != nil {
		return 0, nil, err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)
	if rqBody != nil {
		rq.Header.Set("Content-Type", "application/json")
	}

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
		var lavalinkErr lavalink.Error
		if err = json.Unmarshal(rawBody, &lavalinkErr); err != nil {
			return rs.StatusCode, rawBody, fmt.Errorf("error while unmarshalling disgolink error: %w", err)
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
