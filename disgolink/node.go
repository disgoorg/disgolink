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

const maximumConnectDelay = 60 * time.Second

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
		logger: logger.With(slog.String("name", "disgolink_node"), slog.String("node_name", config.Name)),
		Config: config,
		Client: client,
		status: StatusDisconnected,
	}
	node.Rest = newRestClient(logger, node, httpClient)
	return node
}

type Node struct {
	logger *slog.Logger
	Client *Client
	Config NodeConfig
	Rest   *RestClient

	conn   *websocket.Conn
	connMu sync.Mutex

	statusMu sync.Mutex
	status   Status

	statsMu sync.Mutex
	stats   lavalink.Stats
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

func (n *Node) Open(ctx context.Context) error {
	return n.doReconnect(ctx)
}

func (n *Node) open(ctx context.Context) error {
	n.logger.DebugContext(ctx, "opening node connection")

	n.connMu.Lock()
	defer n.connMu.Unlock()
	if n.conn != nil {
		return ErrNodeAlreadyConnected
	}
	n.statusMu.Lock()
	if n.status != StatusDisconnected {
		n.statusMu.Unlock()
		return ErrNodeAlreadyConnected
	}
	n.status = StatusConnecting
	n.statusMu.Unlock()

	header := http.Header{
		"Authorization": []string{n.Config.Password},
		"User-Id":       []string{n.Client.userID.String()},
		"Client-Name":   []string{fmt.Sprintf("%s/%s", Name, Version)},
	}

	sessionID := n.Config.SessionID
	if sessionID == "" {
		sessionID = n.Config.SessionID
	}
	if sessionID != "" {
		header.Add("Session-Id", sessionID)
	}

	conn, rs, err := websocket.DefaultDialer.DialContext(ctx, n.Config.WsURL(), header)
	if err != nil {
		var body []byte
		if rs != nil {
			data, readErr := io.ReadAll(rs.Body)
			if readErr != nil {
				n.logger.ErrorContext(ctx, "error while reading response body", slog.Any("err", readErr))
			}
			body = data
		}
		n.logger.ErrorContext(ctx, "error connecting to the node", slog.Any("err", err), slog.String("body", string(body)))
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	n.conn = conn

	n.statusMu.Lock()
	n.status = StatusConnected
	n.statusMu.Unlock()

	go n.listen(conn)

	return nil
}

func (n *Node) Close() {
	n.connMu.Lock()
	defer n.connMu.Unlock()

	n.statusMu.Lock()
	if n.status != StatusDisconnected {
		n.statusMu.Unlock()
		return
	}
	n.status = StatusDisconnected
	n.statusMu.Unlock()

	for plugin := range n.Client.Plugins() {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNodeClose(n)
		}
	}
	if n.conn != nil {
		_ = n.conn.Close()
		n.conn = nil
	}
}

func (n *Node) doReconnect(ctx context.Context) error {
	var (
		try              int
		backoffIncrement int
	)

	for {
		// Exponentially backoff up to a limit of 10s
		delay := time.Duration(1<<backoffIncrement) * time.Second
		if delay > maximumConnectDelay {
			delay = maximumConnectDelay
		} else {
			backoffIncrement++
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}

		err := n.open(ctx)
		if err == nil {
			// Successfully connected, our job here is done
			return nil
		}

		if errors.Is(err, ErrNodeAlreadyConnected) {
			return err
		}

		n.logger.ErrorContext(ctx, "failed to reconnect node", slog.Any("err", err), slog.Int("try", try), slog.Duration("delay", delay))
		n.statusMu.Lock()
		n.status = StatusDisconnected
		n.statusMu.Unlock()

		try++
	}
}

func (n *Node) reconnect() {
	if err := n.doReconnect(context.Background()); err != nil {
		n.logger.Error("failed to reopen node", slog.Any("err", err))
	}
}

func (n *Node) listen(conn *websocket.Conn) {
	defer n.logger.Debug("exiting listen goroutine")

	for {
		mt, r, err := conn.NextReader()
		if err != nil {
			n.connMu.Lock()
			sameConn := n.conn == conn
			n.connMu.Unlock()

			if !sameConn {
				return
			}

			reconnect := true
			if errors.Is(err, net.ErrClosed) {
				reconnect = false
			} else {
				n.logger.Warn("failed to read next message from node", slog.Any("err", err))
			}

			n.Close()
			if reconnect {
				go n.reconnect()
			}

			return
		}

		message, data, err := n.parseMessage(mt, r)
		if err != nil {
			n.logger.Error("error while parsing gateway message", slog.Any("err", err))
			continue
		}

		for plugin := range n.Client.Plugins() {
			if pl, ok := plugin.(OpPlugin); ok && pl.Op() == message.Op() {
				pl.OnOpInvocation(n, data)
			}
			if pl, ok := plugin.(OpPlugins); ok {
				for _, pls := range pl.OpPlugins() {
					if pls.Op() == message.Op() {
						pls.OnOpInvocation(n, data)
					}
				}
			}
		}

		switch m := message.(type) {
		case lavalink.ReadyMessage:
			n.Config.SessionID = m.SessionID
			if m.Resumed {
				n.logger.Info("successfully resumed session", slog.String("session_id", m.SessionID))
			} else {
				n.logger.Info("successfully opened session", slog.String("session_id", m.SessionID))
			}

			for plugin := range n.Client.Plugins() {
				if pl, ok := plugin.(PluginEventHandler); ok {
					pl.OnNodeOpen(n)
				}
			}
			n.Client.EmitEvent(&ReadyEvent{
				GenericEvent: newGenericEvent(n),
				ReadyMessage: m,
			})

		case lavalink.StatsMessage:
			n.statsMu.Lock()
			n.stats = m.Stats
			n.statsMu.Unlock()
			n.Client.EmitEvent(&StatsEvent{
				GenericEvent: newGenericEvent(n),
				StatsMessage: m,
			})

		case lavalink.PlayerUpdateMessage:
			player := n.Client.ExistingPlayer(m.GuildID)
			if player == nil {
				continue
			}
			player.OnPlayerUpdate(m.State)
			n.Client.EmitEvent(&PlayerUpdateEvent{
				GenericEvent:        newGenericEvent(n),
				PlayerUpdateMessage: m,
				Player:              player,
			})

		case lavalink.Event:
			player := n.Client.ExistingPlayer(m.GetGuildID())
			if player == nil {
				continue
			}
			player.OnEvent(m)

			switch e := m.(type) {
			case lavalink.TrackStartEvent:
				n.Client.EmitEvent(&PlayerTrackStartEvent{
					GenericEvent:    newGenericEvent(n),
					TrackStartEvent: e,
					Player:          player,
				})

			case lavalink.TrackEndEvent:
				n.Client.EmitEvent(&PlayerTrackEndEvent{
					GenericEvent:  newGenericEvent(n),
					TrackEndEvent: e,
					Player:        player,
				})

			case lavalink.TrackExceptionEvent:
				n.Client.EmitEvent(&PlayerTrackExceptionEvent{
					GenericEvent:        newGenericEvent(n),
					TrackExceptionEvent: e,
					Player:              player,
				})

			case lavalink.TrackStuckEvent:
				n.Client.EmitEvent(&PlayerTrackStuckEvent{
					GenericEvent:    newGenericEvent(n),
					TrackStuckEvent: e,
					Player:          player,
				})

			case lavalink.WebSocketClosedEvent:
				n.Client.EmitEvent(&PlayerWebSocketClosedEvent{
					GenericEvent:         newGenericEvent(n),
					WebSocketClosedEvent: e,
					Player:               player,
				})

			case lavalink.UnknownEvent:
				n.Client.EmitEvent(&UnknownPlayerEvent{
					GenericEvent: newGenericEvent(n),
					UnknownEvent: e,
					Player:       player,
				})
			}

		case lavalink.UnknownMessage:
			n.Client.EmitEvent(&UnknownEvent{
				GenericEvent:   newGenericEvent(n),
				UnknownMessage: m,
			})
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
	n.logger.Debug("received gateway message", slog.String("data", string(data)))

	message, err := lavalink.UnmarshalMessage(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return message, data, nil
}
