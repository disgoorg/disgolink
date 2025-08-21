package disgolink

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

type NodeConfig struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Password  string `json:"password"`
	Secure    bool   `json:"secure"`
	SessionID string `json:"session_id"`
}

func (c NodeConfig) RestURL() string {
	scheme := "http"
	if c.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Address)
}

func (c NodeConfig) WsURL() string {
	scheme := "ws"
	if c.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s%s", scheme, c.Address, EndpointWebSocket)
}

type Status string

// Indicates how far along the Client is to connecting
const (
	StatusConnecting   Status = "CONNECTING"
	StatusConnected    Status = "CONNECTED"
	StatusDisconnected Status = "DISCONNECTED"
)

var ErrNodeAlreadyConnected = errors.New("node already connected")

func newNode(logger *slog.Logger, config NodeConfig, client *Client, httpClient *http.Client) *Node {
	node := &Node{
		Logger: logger.With(slog.String("name", "disgolink_node"), slog.String("node_name", config.Name)),
		Config: config,
		Client: client,
		status: StatusDisconnected,
	}
	node.Rest = newRestClient(logger, node, httpClient)
	return node
}

type Node struct {
	Logger *slog.Logger
	Client *Client
	Config NodeConfig
	Rest   RestClient

	conn   *websocket.Conn
	connMu sync.Mutex

	statusMu sync.Mutex
	status   Status

	statsMu sync.Mutex
	stats   lavalink.Stats

	SessionID string
}

func (n *Node) Status() Status {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()

	return n.status
}

func (n *Node) Stats() lavalink.Stats {
	n.statsMu.Lock()
	defer n.statsMu.Unlock()

	return n.stats
}

func (n *Node) Version(ctx context.Context) (string, error) {
	return n.Rest.Version(ctx)
}

func (n *Node) Info(ctx context.Context) (*lavalink.Info, error) {
	return n.Rest.Info(ctx)
}

func (n *Node) Update(ctx context.Context, update lavalink.SessionUpdate) error {
	_, err := n.Rest.UpdateSession(ctx, n.SessionID, update)
	return err
}

func (n *Node) LoadTracks(ctx context.Context, identifier string) (*lavalink.LoadResult, error) {
	return n.Rest.LoadTracks(ctx, identifier)
}

func (n *Node) LoadTracksHandler(ctx context.Context, identifier string, handler AudioLoadResultHandler) {
	result, err := n.LoadTracks(ctx, identifier)
	if err != nil {
		handler.LoadFailed(err)
		return
	}

	switch d := result.Data.(type) {
	case lavalink.Track:
		handler.TrackLoaded(d)

	case lavalink.Playlist:
		handler.PlaylistLoaded(d)

	case lavalink.Search:
		handler.SearchResultLoaded(d)

	case lavalink.Empty:
		handler.NoMatches()

	case lavalink.Exception:
		handler.LoadFailed(d)
	}
}

func (n *Node) DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error) {
	return n.Rest.DecodeTrack(ctx, encodedTrack)
}

func (n *Node) DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error) {
	return n.Rest.DecodeTracks(ctx, encodedTracks)
}

func (n *Node) Open(ctx context.Context) error {
	return n.reconnectTry(ctx, 0)
}

func (n *Node) open(ctx context.Context) error {
	n.Logger.Debug("opening connection to node...")

	n.connMu.Lock()
	if n.conn != nil {
		n.connMu.Unlock()
		return ErrNodeAlreadyConnected
	}
	n.statusMu.Lock()
	n.status = StatusConnecting
	n.statusMu.Unlock()

	header := http.Header{
		"Authorization": []string{n.Config.Password},
		"User-Id":       []string{n.Client.userID.String()},
		"Client-Name":   []string{fmt.Sprintf("%s/%s", Name, Version)},
	}

	sessionID := n.SessionID
	if sessionID == "" {
		sessionID = n.Config.SessionID
	}
	if sessionID != "" {
		header.Add("Session-Id", sessionID)
	}

	conn, rs, err := websocket.DefaultDialer.DialContext(ctx, n.Config.WsURL(), header)
	if err != nil {
		var body string
		if rs != nil {
			defer func() {
				_ = rs.Body.Close()
			}()
			rawBody, _ := io.ReadAll(rs.Body)
			body = string(rawBody)
		}

		n.Logger.Error("error connecting to the node", slog.Any("err", err), slog.String("body", body))
		n.connMu.Unlock()
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	n.conn = conn
	n.connMu.Unlock()

	n.statusMu.Lock()
	n.status = StatusConnected
	n.statusMu.Unlock()

	go n.listen(conn)

	return nil
}

func (n *Node) Close() {
	n.connMu.Lock()
	defer n.connMu.Unlock()

	for plugin := range n.Client.Plugins() {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNodeClose(n)
		}
	}
	if n.conn != nil {
		_ = n.conn.Close()
		n.conn = nil
	}
	n.statsMu.Lock()
	n.status = StatusDisconnected
	n.statusMu.Unlock()

}

func (n *Node) reconnectTry(ctx context.Context, try int) error {
	delay := time.Duration(try) * 2 * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	if err := n.open(ctx); err != nil {
		if errors.Is(err, ErrNodeAlreadyConnected) {
			return err
		}

		n.Logger.ErrorContext(ctx, "failed to reconnect node", slog.Any("err", err), slog.Int("try", try), slog.Duration("delay", delay))
		n.statusMu.Lock()
		n.status = StatusDisconnected
		n.statusMu.Unlock()
		return n.reconnectTry(ctx, try+1)
	}
	return nil
}

func (n *Node) reconnect() {
	if err := n.reconnectTry(context.Background(), 0); err != nil {
		n.Logger.Error("failed to reopen node", slog.Any("err", err))
	}
}

func (n *Node) listen(conn *websocket.Conn) {
	defer n.Logger.Debug("exiting listen goroutine")
loop:
	for {
		mt, r, err := conn.NextReader()
		if err != nil {
			n.connMu.Lock()
			sameConnection := n.conn == conn
			n.connMu.Unlock()

			if !sameConnection {
				return
			}

			reconnect := true
			if errors.Is(err, net.ErrClosed) {
				reconnect = false
			} else {
				n.Logger.Error("failed to read next message from node", slog.Any("err", err))
			}

			n.Close()
			if reconnect {
				go n.reconnect()
			}

			break loop
		}

		message, data, err := n.parseMessage(mt, r)
		if err != nil {
			n.Logger.Error("error while parsing gateway message", slog.Any("err", err))
			continue
		}

		for plugin := range n.Client.Plugins() {
			if pl, ok := plugin.(PluginEventHandler); ok {
				pl.OnNodeMessageIn(n, data)
			}
		}

		switch m := message.(type) {
		case lavalink.UnknownMessage:
			for plugin := range n.Client.Plugins() {
				if pl, ok := plugin.(OpPlugin); ok {
					pl.OnOpInvocation(n, m.Data)
				}
			}

		case lavalink.ReadyMessage:
			n.SessionID = m.SessionID
			if m.Resumed {
				n.Logger.Info("successfully resumed session", slog.String("session_id", m.SessionID))
			} else {
				n.Logger.Info("successfully opened session", slog.String("session_id", m.SessionID))
			}

			for plugin := range n.Client.Plugins() {
				if pl, ok := plugin.(PluginEventHandler); ok {
					pl.OnNodeOpen(n)
				}
			}

		case lavalink.StatsMessage:
			n.statsMu.Lock()
			n.stats = m.Stats
			n.statsMu.Unlock()
			n.Client.EmitEvent(nil, m)

		case lavalink.PlayerUpdateMessage:
			player := n.Client.ExistingPlayer(m.GuildID)
			if player == nil {
				continue
			}
			player.OnPlayerUpdate(m.State)
			n.Client.EmitEvent(player, m)

		case lavalink.Event:
			player := n.Client.ExistingPlayer(m.GetGuildID())
			if player == nil {
				continue
			}
			player.OnEvent(m)
			n.Client.EmitEvent(player, m)
		}
	}
}

func (n *Node) parseMessage(mt int, r io.Reader) (lavalink.Message, []byte, error) {
	if mt != websocket.TextMessage {
		return nil, nil, fmt.Errorf("expected text message, got %d", mt)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read message: %w", err)
	}
	n.Logger.Debug("received gateway message", slog.String("data", string(data)))

	message, err := lavalink.UnmarshalMessage(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return message, data, nil
}
