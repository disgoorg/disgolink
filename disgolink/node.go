package disgolink

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/gorilla/websocket"
)

type Status string

// Indicates how far along the client is to connecting
const (
	StatusConnecting   Status = "CONNECTING"
	StatusConnected    Status = "CONNECTED"
	StatusReconnecting Status = "RECONNECTING"
	StatusDisconnected Status = "DISCONNECTED"
)

var ErrNodeAlreadyConnected = errors.New("node already connected")

var _ Node = (*nodeImpl)(nil)

type Node interface {
	Lavalink() Client
	Config() NodeConfig
	Rest() RestClient

	Stats() lavalink.Stats
	Status() Status
	SessionID() string

	Version(ctx context.Context) (string, error)
	Info(ctx context.Context) (*lavalink.Info, error)
	Update(ctx context.Context, update lavalink.SessionUpdate) error
	LoadTracks(ctx context.Context, identifier string) (*lavalink.LoadResult, error)
	LoadTracksHandler(ctx context.Context, identifier string, handler AudioLoadResultHandler)

	DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error)

	Open(ctx context.Context) error
	Close()
}

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

type nodeImpl struct {
	lavalink Client
	config   NodeConfig
	rest     RestClient

	conn   *websocket.Conn
	connMu sync.Mutex

	status    Status
	stats     lavalink.Stats
	sessionID string
}

func (n *nodeImpl) Lavalink() Client {
	return n.lavalink
}

func (n *nodeImpl) Config() NodeConfig {
	return n.config
}

func (n *nodeImpl) Rest() RestClient {
	return n.rest
}

func (n *nodeImpl) Status() Status {
	return n.status
}

func (n *nodeImpl) Stats() lavalink.Stats {
	return n.stats
}

func (n *nodeImpl) SessionID() string {
	return n.sessionID
}

func (n *nodeImpl) Version(ctx context.Context) (string, error) {
	return n.rest.Version(ctx)
}

func (n *nodeImpl) Info(ctx context.Context) (*lavalink.Info, error) {
	return n.rest.Info(ctx)
}

func (n *nodeImpl) Update(ctx context.Context, update lavalink.SessionUpdate) error {
	session, err := n.rest.UpdateSession(ctx, n.sessionID, update)
	if session != nil && session.Resuming {
		n.config.SessionID = n.sessionID
	}
	return err
}

func (n *nodeImpl) LoadTracks(ctx context.Context, identifier string) (*lavalink.LoadResult, error) {
	return n.rest.LoadTracks(ctx, identifier)
}

func (n *nodeImpl) LoadTracksHandler(ctx context.Context, identifier string, handler AudioLoadResultHandler) {
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

func (n *nodeImpl) syncPlayers(ctx context.Context) error {
	players, err := n.rest.Players(ctx, n.sessionID)
	if err != nil {
		return err
	}

	for _, player := range players {
		p := n.lavalink.PlayerOnNode(n, player.GuildID)
		if p == nil {
			continue
		}
		p.Restore(player)
	}

	return nil
}

func (n *nodeImpl) DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error) {
	return n.rest.DecodeTrack(ctx, encodedTrack)
}

func (n *nodeImpl) DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error) {
	return n.rest.DecodeTracks(ctx, encodedTracks)
}

func (n *nodeImpl) Open(ctx context.Context) error {
	return n.reconnectTry(ctx, 0, false)
}

func (n *nodeImpl) open(ctx context.Context, reconnecting bool) error {
	n.lavalink.Logger().Debug("opening connection to node...")

	n.connMu.Lock()
	defer n.connMu.Unlock()
	if n.conn != nil {
		return ErrNodeAlreadyConnected
	}

	if reconnecting {
		n.status = StatusReconnecting
	} else {
		n.status = StatusConnecting
	}

	header := http.Header{
		"Authorization": []string{n.config.Password},
		"User-Id":       []string{n.lavalink.UserID().String()},
		"Client-Name":   []string{fmt.Sprintf("%s/%s", Name, Version)},
	}
	if n.config.SessionID != "" {
		header.Add("Session-Id", n.config.SessionID)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, n.config.WsURL(), header)
	if err != nil {
		return err
	}

	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	message, err := lavalink.UnmarshalMessage(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ready message. error: %w", err)
	}
	ready, ok := message.(lavalink.ReadyMessage)
	if !ok {
		return fmt.Errorf("expected ready message but got %T", message)
	}

	n.sessionID = ready.SessionID
	if n.config.SessionID != "" {
		if ready.Resumed {
			n.lavalink.Logger().Info("successfully resumed session: ", n.config.SessionID)
			if err = n.syncPlayers(ctx); err != nil {
				n.lavalink.Logger().Warn("failed to sync players: ", err)
			}
		} else {
			n.lavalink.Logger().Warn("failed to resume session: ", n.config.SessionID)
		}
	}
	n.status = StatusConnected

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	n.conn = conn

	go n.listen(conn)

	n.Lavalink().ForPlugins(func(plugin Plugin) {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNodeOpen(n)
		}
	})

	return nil
}

func (n *nodeImpl) Close() {
	n.Lavalink().ForPlugins(func(plugin Plugin) {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNodeClose(n)
		}
	})
	n.status = StatusDisconnected
	if n.conn != nil {
		_ = n.conn.Close()
		n.conn = nil
	}

}

func (n *nodeImpl) reconnectTry(ctx context.Context, try int, reconnecting bool) error {
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

	if err := n.open(ctx, reconnecting); err != nil {
		if err == ErrNodeAlreadyConnected {
			return err
		}
		n.lavalink.Logger().Error("failed to reconnect node. error: ", err)
		n.status = StatusDisconnected
		return n.reconnectTry(ctx, try+1, reconnecting)
	}
	return nil
}

func (n *nodeImpl) reconnect() {
	if err := n.reconnectTry(context.Background(), 0, true); err != nil {
		n.lavalink.Logger().Error("failed to reopen node. error: ", err)
	}
}

func (n *nodeImpl) listen(conn *websocket.Conn) {
	defer n.lavalink.Logger().Debug("exiting listen goroutine...")
loop:
	for {
		_, data, err := conn.ReadMessage()
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
			}

			n.Close()
			if reconnect {
				go n.reconnect()
			}
			break loop
		}

		n.lavalink.Logger().Trace("received message: ", string(data))

		n.Lavalink().ForPlugins(func(plugin Plugin) {
			if pl, ok := plugin.(PluginEventHandler); ok {
				pl.OnNodeMessageIn(n, data)
			}
		})

		m, err := lavalink.UnmarshalMessage(data)
		if err != nil {
			n.lavalink.Logger().Errorf("error while unmarshalling ws data: %s", err)
			return
		}

		switch message := m.(type) {
		case lavalink.UnknownMessage:
			n.Lavalink().ForPlugins(func(plugin Plugin) {
				if pl, ok := plugin.(OpPlugin); ok {
					pl.OnOpInvocation(n, message.Data)
				}
			})

		case lavalink.StatsMessage:
			n.stats = lavalink.Stats(message)

		case lavalink.PlayerUpdateMessage:
			player := n.lavalink.ExistingPlayer(message.GuildID)
			if player == nil {
				continue
			}
			player.OnPlayerUpdate(message.State)

		case lavalink.Event:
			player := n.lavalink.ExistingPlayer(message.GuildID())
			if player == nil {
				continue
			}
			player.OnEvent(message)
		}
	}
}
