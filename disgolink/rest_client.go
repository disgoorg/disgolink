package disgolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Endpoint string

func (e Endpoint) Format(a ...any) string {
	return fmt.Sprintf(string(e), a...)
}

var (
	EndpointBase    Endpoint = "/v4"
	EndpointVersion Endpoint = "/version"
	EndpointInfo             = EndpointBase + "/info"
	EndpointStats            = EndpointBase + "/stats"

	EndpointUpdateSession = EndpointBase + "/sessions/%s"
	EndpointPlayers       = EndpointBase + "/sessions/%s/players"
	EndpointPlayer        = EndpointBase + "/sessions/%s/players/%s"
	EndpointUpdatePlayer  = EndpointBase + "/sessions/%s/players/%s?noReplace=%t"
	EndpointDestroyPlayer = EndpointBase + "/sessions/%s/players/%s"

	EndpointLoadTracks   = EndpointBase + "/loadtracks?identifier=%s"
	EndpointDecodeTrack  = EndpointBase + "/decodetrack?track=%s"
	EndpointDecodeTracks = EndpointBase + "/decodetracks"

	EndpointWebSocket = EndpointBase + "/websocket"
)

type RestClient interface {
	// Do executes a http.Request and replaces the host and scheme with the node's config. It also sets the Authorization header to the node's password. It returns the http.Response or an error
	Do(rq *http.Request) (*http.Response, error)

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

type restClientImpl struct {
	logger     *slog.Logger
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
	err = c.doJSON(ctx, http.MethodPatch, EndpointUpdateSession.Format(sessionID), sessionUpdate, &session)
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

func (c *restClientImpl) Do(rq *http.Request) (*http.Response, error) {
	rq.Header.Set("Authorization", c.node.Config().Password)
	rq.URL.Host = c.node.Config().Address
	if c.node.Config().Secure {
		rq.URL.Scheme = "https"
	} else {
		rq.URL.Scheme = "http"
	}
	return c.httpClient.Do(rq)
}

func (c *restClientImpl) do(ctx context.Context, method string, path string, rqBody []byte) (int, []byte, error) {
	rq, err := http.NewRequestWithContext(ctx, method, c.node.Config().RestURL()+path, bytes.NewReader(rqBody))
	if err != nil {
		return 0, nil, err
	}
	rq.Header.Set("Authorization", c.node.Config().Password)
	if len(rqBody) > 0 {
		rq.Header.Set("Content-Type", "application/json")
	}

	c.logger.DebugContext(ctx, "sending request", slog.String("method", method), slog.String("path", path), slog.String("body", fmt.Sprintf("%v", rqBody)))

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return 0, nil, err
	}

	defer rs.Body.Close()
	rawBody, err := io.ReadAll(rs.Body)
	c.logger.DebugContext(ctx, "received response", slog.String("path", path), slog.Int("status_code", rs.StatusCode), slog.String("body", string(rawBody)))
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
	var rawRqBody []byte
	if rqBody != nil {
		var err error
		rawRqBody, err = json.Marshal(rqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}
	statusCode, rawBody, err := c.do(ctx, method, path, rawRqBody)
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
