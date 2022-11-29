package lavalink

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/disgoorg/disgolink/lavalink/protocol"
	"github.com/gorilla/websocket"
)

type NodeStatus string

// Indicates how far along the client is to connecting
const (
	Connecting   NodeStatus = "CONNECTING"
	Connected    NodeStatus = "CONNECTED"
	Reconnecting NodeStatus = "RECONNECTING"
	Disconnected NodeStatus = "DISCONNECTED"
)

type Node interface {
	Lavalink() Lavalink
	Config() NodeConfig
	Rest() RestClient

	Stats() protocol.Stats
	Status() NodeStatus
	SessionID() string

	Open(ctx context.Context) error
	Close() error
}

type NodeConfig struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	Secure      bool   `json:"secure"`
	ResumingKey string `json:"resumingKey"`
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

	status    NodeStatus
	stats     protocol.Stats
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

func (n *nodeImpl) Status() NodeStatus {
	return n.status
}

func (n *nodeImpl) Stats() protocol.Stats {
	return n.stats
}

func (n *nodeImpl) SessionID() string {
	return n.sessionID
}

func (n *nodeImpl) Open(ctx context.Context) error {
	n.status = Connecting
	if err := n.open(ctx); err != nil {
		n.status = Disconnected
		return err
	}
	n.status = Connected
	return nil
}

func (n *nodeImpl) open(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	scheme := "ws"
	if n.config.Secure {
		scheme += "s"
	}
	header := http.Header{}
	header.Add("Authorization", n.config.Password)
	header.Add("User-Id", n.lavalink.UserID().String())
	header.Add("Client-Name", fmt.Sprintf("%s/%s", Name, Version))
	if n.config.ResumingKey != "" {
		header.Add("Resume-Key", n.config.ResumingKey)
	}

	var (
		err error
		rs  *http.Response
	)
	n.conn, rs, err = websocket.DefaultDialer.DialContext(ctx, fmt.Sprintf("%s://%s:%s", scheme, n.config.Host, n.config.Port), header)
	if err != nil {
		return err
	}
	if n.config.ResumingKey != "" {
		if rs.Header.Get("Session-Resumed") == "true" {
			n.lavalink.Logger().Info("successfully resumed session with key: %s", n.config.ResumingKey)
		} else {
			n.lavalink.Logger().Warn("failed to resume session with key: ", n.config.ResumingKey)
		}
	}

	go n.listen()

	return err
}

func (n *nodeImpl) listen(conn *websocket.Conn) {
	defer n.lavalink.Logger().Debug("exiting listen goroutine...")
loop:
	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			n.connMu.Lock()
			sameConnection := n.conn == conn
			n.connMu.Unlock()

			if !sameConnection {
				return
			}

			break loop
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			n.lavalink.Logger().Errorf("error while reading ws data: %s", err)
			continue
		}

		m, err := protocol.UnmarshalMessage(data)
		if err != nil {
			n.lavalink.Logger().Errorf("error while unmarshalling ws data: %s", err)
			return
		}

		switch message := m.(type) {
		case protocol.MessageStats:

		case protocol.MessagePlayerUpdate:

		case protocol.Event:
			n.lavalink.Player(message.GuildID())
		}

	}
}

func (n *nodeImpl) Close() {
	n.status = Disconnected
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			n.lavalink.Logger().Errorf("error while closing wsconn: %s", err)
		}
	}
}
