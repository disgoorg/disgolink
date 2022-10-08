package lavalink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
	Send(cmd OpCommand) error
	ConfigureResuming(key string, timeoutSeconds int) error

	Open(ctx context.Context) error
	Close()

	Name() string
	RestClient() RestClient
	RestURL() string
	Config() NodeConfig
	Stats() *Stats
	Status() NodeStatus
}

type NodeConfig struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	Secure      bool   `json:"secure"`
	ResumingKey string `json:"resumingKey"`
}

type nodeImpl struct {
	config     NodeConfig
	lavalink   Lavalink
	conn       *websocket.Conn
	status     NodeStatus
	statusMu   sync.Mutex
	stats      *Stats
	available  bool
	restClient RestClient
}

func (n *nodeImpl) RestURL() string {
	scheme := "http"
	if n.config.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, n.config.Host, n.config.Port)
}

func (n *nodeImpl) Lavalink() Lavalink {
	return n.lavalink
}

func (n *nodeImpl) Config() NodeConfig {
	return n.config
}

func (n *nodeImpl) RestClient() RestClient {
	return n.restClient
}

func (n *nodeImpl) Name() string {
	return n.config.Name
}

func (n *nodeImpl) Send(cmd OpCommand) error {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()

	if n.status != Connected {
		return fmt.Errorf("node is %s and cannot send a cmd to the node", n.status)
	}
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	for _, pl := range n.Lavalink().Plugins() {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNodeMessageOut(n, data)
		}
	}
	if err = n.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("error while sending to lavalink websocket: %w", err)
	}

	return nil
}

func (n *nodeImpl) ConfigureResuming(key string, timeoutSeconds int) error {
	return n.Send(ConfigureResumingCommand{
		Key:     key,
		Timeout: timeoutSeconds,
	})
}

func (n *nodeImpl) Status() NodeStatus {
	return n.status
}

func (n *nodeImpl) Stats() *Stats {
	return n.stats
}

func (n *nodeImpl) reconnect(ctx context.Context) {
	if err := n.reconnectTry(ctx, 0, time.Second); err != nil {
		n.lavalink.Logger().Error("failed to reconnect to node: ", err)
	}
}

func (n *nodeImpl) reconnectTry(ctx context.Context, try int, delay time.Duration) error {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.status = Reconnecting

	timer := time.NewTimer(time.Duration(try) * delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	n.lavalink.Logger().Debug("reconnecting gateway...")
	if err := n.open(ctx); err != nil {
		n.lavalink.Logger().Error("failed to reconnect node. error: ", err)
		n.status = Disconnected
		return n.reconnectTry(ctx, try+1, delay)
	}
	n.status = Connected
	return nil
}

func (n *nodeImpl) listen() {
	defer n.lavalink.Logger().Debug("shutting down listen goroutine")
loop:
	for {
		if n.conn == nil {
			return
		}
		_, data, err := n.conn.ReadMessage()
		if err != nil {
			n.Close()
			if !errors.Is(err, net.ErrClosed) {
				go n.reconnect(context.TODO())
			}
			break loop
		}

		n.lavalink.Logger().Trace("received: ", string(data))

		for _, pl := range n.Lavalink().Plugins() {
			if plugin, ok := pl.(PluginEventHandler); ok {
				plugin.OnNodeMessageIn(n, data)
			}
		}

		var v UnmarshalOp
		if err = json.Unmarshal(data, &v); err != nil {
			n.lavalink.Logger().Error("error while unmarshalling op. error: ", err)
			continue
		}

		switch op := v.Op.(type) {
		case UnknownOp:
			for _, pl := range n.Lavalink().Plugins() {
				if plugin, ok := pl.(OpPlugin); ok {
					plugin.OnOpInvocation(n, op.Data)
				}
				if plugin, ok := pl.(OpPlugins); ok {
					for _, ext := range plugin.OpPlugins() {
						if ext.Op() == op.Op() {
							ext.OnOpInvocation(n, op.Data)
						}
					}
				}
			}

		case PlayerUpdateOp:
			n.onPlayerUpdate(op)

		case OpEvent:
			n.onEvent(op)

		case StatsOp:
			n.onStatsEvent(op)

		default:
			n.lavalink.Logger().Warnf("unexpected op received: %T, data: ", op, string(data))
		}
	}
}

func (n *nodeImpl) onPlayerUpdate(playerUpdate PlayerUpdateOp) {
	if player := n.lavalink.ExistingPlayer(playerUpdate.GuildID); player != nil {
		player.OnPlayerUpdate(playerUpdate.State)
		player.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnPlayerUpdate(player, playerUpdate.State)
			}
		})
		return
	}
	n.lavalink.Logger().Warnf("player update received for unknown player: %s", playerUpdate.GuildID)
}

func (n *nodeImpl) onEvent(event OpEvent) {
	player := n.lavalink.ExistingPlayer(event.GuildID())
	if player == nil {
		return
	}

	switch e := event.(type) {
	case TrackEvent:
		player.OnEvent(e)

	case WebsocketClosedEvent:
		player.EmitEvent(func(l any) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnWebSocketClosed(player, e.Code, e.Reason, e.ByRemote)
			}
		})

	case UnknownEvent:
		for _, pl := range n.Lavalink().Plugins() {
			if plugin, ok := pl.(EventPlugin); ok {
				plugin.OnEventInvocation(n, e.Data)
			}
			if plugin, ok := pl.(EventPlugins); ok {
				for _, ext := range plugin.EventPlugins() {
					if ext.Event() == e.Event() {
						ext.OnEventInvocation(n, e.Data)
					}
				}
			}
		}

	default:
		n.lavalink.Logger().Warnf("unexpected event received: %T, data: ", event)
		return
	}
}

func (n *nodeImpl) onStatsEvent(stats StatsOp) {
	n.stats = &stats.Stats
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

func (n *nodeImpl) Open(ctx context.Context) error {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()

	n.status = Connecting
	if err := n.open(ctx); err != nil {
		n.status = Disconnected
		return err
	}
	n.status = Connected
	return nil
}

func (n *nodeImpl) Close() {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()

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
