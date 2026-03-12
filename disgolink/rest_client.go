package disgolink

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v4/lavalink"
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

func newRestClient(logger *slog.Logger, node *Node, httpClient *http.Client) *RestClient {
	return &RestClient{
		logger:     logger.With(slog.String("name", "disgolink_rest_client"), slog.String("node_name", node.Config.Name)),
		node:       node,
		httpClient: httpClient,
	}
}

type RestClient struct {
	logger     *slog.Logger
	node       *Node
	httpClient *http.Client
}

func (c *RestClient) Version(ctx context.Context) (string, error) {
	_, rawBody, err := c.do(ctx, http.MethodGet, string(EndpointVersion), nil)
	if err != nil {
		return "", err
	}
	return string(rawBody), nil
}

func (c *RestClient) Info(ctx context.Context) (info *lavalink.Info, err error) {
	err = c.doJSON(ctx, http.MethodGet, string(EndpointInfo), nil, &info)
	return
}

func (c *RestClient) Stats(ctx context.Context) (stats *lavalink.Stats, err error) {
	err = c.doJSON(ctx, http.MethodGet, string(EndpointStats), nil, &stats)
	return
}

func (c *RestClient) UpdateSession(ctx context.Context, sessionUpdate lavalink.SessionUpdate) (session *lavalink.Session, err error) {
	err = c.doJSON(ctx, http.MethodPatch, EndpointUpdateSession.Format(c.node.Config.SessionID), sessionUpdate, &session)
	return
}

func (c *RestClient) Players(ctx context.Context) (players []lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointPlayers.Format(c.node.Config.SessionID), nil, &players)
	return
}

func (c *RestClient) Player(ctx context.Context, guildID snowflake.ID) (player *lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointPlayer.Format(c.node.Config.SessionID, guildID), nil, &player)
	return
}

func (c *RestClient) UpdatePlayer(ctx context.Context, guildID snowflake.ID, playerUpdate lavalink.PlayerUpdate) (player *lavalink.Player, err error) {
	err = c.doJSON(ctx, http.MethodPatch, EndpointUpdatePlayer.Format(c.node.Config.SessionID, guildID, playerUpdate.NoReplace), playerUpdate, &player)
	return
}

func (c *RestClient) DestroyPlayer(ctx context.Context, guildID snowflake.ID) error {
	_, _, err := c.do(ctx, http.MethodDelete, EndpointDestroyPlayer.Format(c.node.Config.SessionID, guildID), nil)
	return err
}

func (c *RestClient) LoadTracks(ctx context.Context, identifier string) (result *lavalink.LoadResult, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointLoadTracks.Format(url.QueryEscape(identifier)), nil, &result)
	return
}

func (c *RestClient) LoadTracksHandler(ctx context.Context, identifier string, handler TrackLoadingResultHandler) {
	result, err := c.LoadTracks(ctx, identifier)
	if err != nil {
		handler.OnError(err)
		return
	}

	switch d := result.Data.(type) {
	case lavalink.Track:
		handler.OnTrack(d)

	case lavalink.Playlist:
		handler.OnPlaylist(d)

	case lavalink.Search:
		handler.OnSearch(d)

	case lavalink.Empty:
		handler.OnEmpty()

	case lavalink.Exception:
		handler.OnError(d)
	}
}

func (c *RestClient) DecodeTrack(ctx context.Context, encodedTrack string) (track *lavalink.Track, err error) {
	err = c.doJSON(ctx, http.MethodGet, EndpointDecodeTrack.Format(url.QueryEscape(encodedTrack)), nil, &track)
	return
}

func (c *RestClient) DecodeTracks(ctx context.Context, encodedTracks []string) (tracks []lavalink.Track, err error) {
	err = c.doJSON(ctx, http.MethodPost, string(EndpointDecodeTracks), encodedTracks, &tracks)
	return
}

func (c *RestClient) Do(rq *http.Request) (*http.Response, error) {
	rq.Header.Set("Authorization", c.node.Config.Password)
	rq.URL.Host = c.node.Config.Address
	if c.node.Config.Secure {
		rq.URL.Scheme = "https"
	} else {
		rq.URL.Scheme = "http"
	}
	return c.httpClient.Do(rq)
}

func (c *RestClient) do(ctx context.Context, method string, path string, rqBody []byte) (int, []byte, error) {
	rq, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(rqBody))
	if err != nil {
		return 0, nil, err
	}

	if len(rqBody) > 0 {
		rq.Header.Set("Content-Type", "application/json")
	}

	c.logger.DebugContext(ctx, "sending request", slog.String("method", method), slog.String("path", path), slog.String("body", string(rqBody)))

	rs, err := c.Do(rq)
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

func (c *RestClient) doJSON(ctx context.Context, method string, path string, rqBody any, rsBody any) error {
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
