package lavalink

import (
	"context"
	"fmt"
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
	status NodeStatus
	stats  protocol.Stats
	mu     sync.Mutex
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

	for _, pl := range n.Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNodeOpen(n)
		}
	}

	return err
}

func (n *nodeImpl) listen() {

}

func (n *nodeImpl) Close() {

	for _, pl := range n.Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNodeDestroy(n)
		}
	}
	n.status = Disconnected
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			n.lavalink.Logger().Errorf("error while closing wsconn: %s", err)
		}
	}
}
