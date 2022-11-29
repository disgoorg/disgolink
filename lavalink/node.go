package lavalink

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
	Lavalink() Lavalink
	Config() NodeConfig
	Rest() RestClient

	Stats() Stats
	Status() Status
	SessionID() string

	Version(ctx context.Context) (string, error)
	Info(ctx context.Context) (*Info, error)
	Update(ctx context.Context, update SessionUpdate) error
	LoadTracks(ctx context.Context, identifier string, handler AudioLoadResultHandler)

	DecodeTrack(ctx context.Context, encodedTrack string) (*Track, error)
	DecodeTracks(ctx context.Context, encodedTracks []string) ([]Track, error)

	Open(ctx context.Context) error
	Close()
}

type NodeConfig struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	Secure      bool   `json:"secure"`
	ResumingKey string `json:"resumingKey"`

	ReconnectTimeout  time.Duration `json:"-"`
	MaxReconnectTries int           `json:"-"`
}

func (c NodeConfig) RestURL() string {
	scheme := "http"
	if c.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, c.Host, c.Port)
}

func (c NodeConfig) WsURL() string {
	scheme := "ws"
	if c.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, c.Host, c.Port)
}

type nodeImpl struct {
	lavalink Lavalink
	config   NodeConfig
	rest     RestClient

	conn   *websocket.Conn
	connMu sync.Mutex

	status    Status
	stats     Stats
	sessionID string
	mu        sync.Mutex
}

func (n *nodeImpl) Lavalink() Lavalink {
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

func (n *nodeImpl) Stats() Stats {
	return n.stats
}

func (n *nodeImpl) SessionID() string {
	return n.sessionID
}

func (n *nodeImpl) Version(ctx context.Context) (string, error) {
	return n.rest.Version(ctx)
}

func (n *nodeImpl) Info(ctx context.Context) (*Info, error) {
	return n.rest.Info(ctx)
}

func (n *nodeImpl) Update(ctx context.Context, update SessionUpdate) error {
	_, err := n.rest.UpdateSession(ctx, n.sessionID, update)
	return err
}

func (n *nodeImpl) LoadTracks(ctx context.Context, identifier string, handler AudioLoadResultHandler) {
	result, err := n.rest.LoadTracks(ctx, identifier)
	if err != nil {
		handler.LoadFailed(Exception{
			Message:  err.Error(),
			Severity: SeverityFault,
		})
	}

	switch result.LoadType {
	case LoadTypeTrackLoaded:
		handler.TrackLoaded(result.Tracks[0])

	case LoadTypePlaylistLoaded:
		result.Playlist.Tracks = result.Tracks
		handler.PlaylistLoaded(*result.Playlist)

	case LoadTypeSearchResult:
		handler.SearchResultLoaded(result.Tracks)

	case LoadTypeNoMatches:
		handler.NoMatches()

	case LoadTypeLoadFailed:
		handler.LoadFailed(*result.Exception)
	}
}

func (n *nodeImpl) DecodeTrack(ctx context.Context, encodedTrack string) (*Track, error) {
	return n.rest.DecodeTrack(ctx, encodedTrack)
}

func (n *nodeImpl) DecodeTracks(ctx context.Context, encodedTracks []string) ([]Track, error) {
	return n.rest.DecodeTracks(ctx, encodedTracks)
}

func (n *nodeImpl) Open(ctx context.Context) error {
	n.lavalink.Logger().Debug("opening connection to lavalink node...")

	n.connMu.Lock()
	defer n.connMu.Unlock()
	if n.conn != nil {
		return ErrNodeAlreadyConnected
	}

	n.status = StatusConnecting

	header := http.Header{
		"Authorization": []string{n.config.Password},
		"User-Id":       []string{n.lavalink.UserID().String()},
		"Client-Name":   []string{fmt.Sprintf("%s/%s", Name, LibraryVersion)},
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

	message, err := UnmarshalMessage(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ready message. error: %w", err)
	}
	ready, ok := message.(MessageReady)
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

func (n *nodeImpl) reconnectTry(ctx context.Context, try int, delay time.Duration) error {
	if try >= n.config.MaxReconnectTries-1 {
		return fmt.Errorf("failed to reconnect. exceeded max reconnect tries of %d reached", n.config.MaxReconnectTries)
	}
	timer := time.NewTimer(time.Duration(try) * delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	n.lavalink.Logger().Debug("reconnecting node...")
	if err := n.Open(ctx); err != nil {
		if err == ErrNodeAlreadyConnected {
			return err
		}
		n.lavalink.Logger().Error("failed to reconnect node. error: ", err)
		n.status = StatusDisconnected
		return n.reconnectTry(ctx, try+1, delay)
	}
	return nil
}

func (n *nodeImpl) reconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), n.config.ReconnectTimeout)
	defer cancel()

	if err := n.reconnectTry(ctx, 0, time.Second); err != nil {
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
				go n.reconnect()
			}

			break loop
		}

		m, err := UnmarshalMessage(data)
		if err != nil {
			n.lavalink.Logger().Errorf("error while unmarshalling ws data: %s", err)
			return
		}

		switch message := m.(type) {
		case MessageStats:

		case MessagePlayerUpdate:

		case Event:
			n.lavalink.Player(message.GuildID())
		}

	}
}
