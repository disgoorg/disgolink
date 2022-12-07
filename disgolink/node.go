package disgolink

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/disgolink/v2/lavalink"
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
	LoadTracks(ctx context.Context, identifier string, handler AudioLoadResultHandler)

	DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error)

	Open(ctx context.Context) error
	Close()
}

type NodeConfig struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	Secure      bool   `json:"secure"`
	ResumingKey string `json:"resumingKey"`
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
	_, err := n.rest.UpdateSession(ctx, n.sessionID, update)
	return err
}

func (n *nodeImpl) LoadTracks(ctx context.Context, identifier string, handler AudioLoadResultHandler) {
	result, err := n.rest.LoadTracks(ctx, identifier)
	if err != nil {
		handler.LoadFailed(err)
	}

	switch result.LoadType {
	case lavalink.LoadTypeTrackLoaded:
		handler.TrackLoaded(result.Tracks[0])

	case lavalink.LoadTypePlaylistLoaded:
		handler.PlaylistLoaded(lavalink.Playlist{
			Info:       result.PlaylistInfo,
			PluginInfo: result.PluginInfo,
			Tracks:     result.Tracks,
		})

	case lavalink.LoadTypeSearchResult:
		handler.SearchResultLoaded(result.Tracks)

	case lavalink.LoadTypeNoMatches:
		handler.NoMatches()

	case lavalink.LoadTypeLoadFailed:
		handler.LoadFailed(result.Exception)
	}
}

func (n *nodeImpl) DecodeTrack(ctx context.Context, encodedTrack string) (*lavalink.Track, error) {
	return n.rest.DecodeTrack(ctx, encodedTrack)
}

func (n *nodeImpl) DecodeTracks(ctx context.Context, encodedTracks []string) ([]lavalink.Track, error) {
	return n.rest.DecodeTracks(ctx, encodedTracks)
}

func (n *nodeImpl) Open(ctx context.Context) error {
	return n.reconnectTry(ctx, 0)
}

func (n *nodeImpl) open(ctx context.Context) error {
	n.lavalink.Logger().Debug("opening connection to node...")

	n.connMu.Lock()
	defer n.connMu.Unlock()
	if n.conn != nil {
		return ErrNodeAlreadyConnected
	}

	n.status = StatusConnecting

	header := http.Header{
		"Authorization": []string{n.config.Password},
		"User-Id":       []string{n.lavalink.UserID().String()},
		"Client-Name":   []string{fmt.Sprintf("%s/%s", Name, Version)},
	}
	if n.config.ResumingKey != "" {
		header.Add("Resume-Key", n.config.ResumingKey)
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
	if n.config.ResumingKey != "" {
		if ready.Resumed {
			n.lavalink.Logger().Info("successfully resumed session with key: %s", n.config.ResumingKey)
		} else {
			n.lavalink.Logger().Warn("failed to resume session with key: ", n.config.ResumingKey)
		}
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	n.conn = conn

	go n.listen(conn)

	return nil
}

func (n *nodeImpl) Close() {
	n.status = StatusDisconnected
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			n.lavalink.Logger().Errorf("error while closing wsconn: %s", err)
		}
	}
}

func (n *nodeImpl) reconnectTry(ctx context.Context, try int) error {
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
		if err == ErrNodeAlreadyConnected {
			return err
		}
		n.lavalink.Logger().Error("failed to reconnect node. error: ", err)
		n.status = StatusDisconnected
		return n.reconnectTry(ctx, try+1)
	}
	return nil
}

func (n *nodeImpl) reconnect() {
	if err := n.reconnectTry(context.Background(), 0); err != nil {
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

			if !errors.Is(err, net.ErrClosed) {
				n.Close()
				go n.reconnect()
			}

			break loop
		}

		m, err := lavalink.UnmarshalMessage(data)
		if err != nil {
			n.lavalink.Logger().Errorf("error while unmarshalling ws data: %s", err)
			return
		}

		switch message := m.(type) {
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
