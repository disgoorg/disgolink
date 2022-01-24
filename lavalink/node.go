package lavalink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	ConfigureResuming(key string, timeout time.Duration) error

	Open(ctx context.Context) error
	Close()

	Name() string
	RestClient() RestClient
	RestURL() string
	Config() NodeConfig
	Stats() Stats
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
	statusMu   sync.Locker
	stats      Stats
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
		return errors.Errorf("node is not %s", n.statusMu)
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
		return errors.Wrap(err, "error while sending to lavalink websocket")
	}
	return nil
}

func (n *nodeImpl) ConfigureResuming(key string, timeout time.Duration) error {
	return n.Send(ConfigureResumingCommand{
		Key:     key,
		Timeout: timeout,
	})
}

func (n *nodeImpl) Status() NodeStatus {
	return n.status
}

func (n *nodeImpl) Stats() Stats {
	return n.stats
}

func (n *nodeImpl) reconnect() error {
	n.statusMu.Lock()
	n.status = Reconnecting
	defer n.statusMu.Unlock()

	if err := n.open(context.TODO(), 0); err != nil {
		n.status = Disconnected
		return err
	}
	n.status = Connected
	return nil
}

func (n *nodeImpl) listen() {
	defer func() {
		n.lavalink.Logger().Info("shut down listen goroutine")
	}()
	for {
		if n.conn == nil {
			return
		}
		_, data, err := n.conn.ReadMessage()
		if err != nil {
			n.lavalink.Logger().Errorf("error while reading from lavalink websocket. error: %s", err)
			n.Close()
			if err := n.reconnect(); err != nil {
				n.lavalink.Logger().Errorf("error while reconnecting to lavalink websocket. error: %s", err)
			}
			return
		}

		n.lavalink.Logger().Debugf("received: %s", string(data))

		for _, pl := range n.Lavalink().Plugins() {
			if plugin, ok := pl.(PluginEventHandler); ok {
				plugin.OnNodeMessageIn(n, data)
			}
		}

		var v UnmarshalOp
		if err = json.Unmarshal(data, &v); err != nil {
			n.lavalink.Logger().Errorf("error while unmarshalling op. error: %s", err)
			continue
		}

		switch op := v.Op.(type) {
		case UnknownOp:
			for _, pl := range n.Lavalink().Plugins() {
				if plugin, ok := pl.(OpExtension); ok {
					plugin.OnOpInvocation(n, op.Data)
				}
				if plugin, ok := pl.(OpExtensions); ok {
					for _, ext := range plugin.OpExtensions() {
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
		player.PlayerUpdate(playerUpdate.State)
		player.EmitEvent(func(l interface{}) {
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
		track, err := n.lavalink.DecodeTrack(e.Track())
		if err != nil {
			n.lavalink.Logger().Errorf("error while decoding track: %s", err)
			return
		}
		switch ee := e.(type) {
		case TrackStartEvent:
			player.SetTrack(track)
			player.EmitEvent(func(l interface{}) {
				if listener := l.(PlayerEventListener); listener != nil {
					listener.OnTrackStart(player, track)
				}
			})

		case TrackEndEvent:
			player.EmitEvent(func(l interface{}) {
				if listener := l.(PlayerEventListener); listener != nil {
					listener.OnTrackEnd(player, track, ee.Reason)
				}
			})

		case TrackExceptionEvent:
			player.EmitEvent(func(l interface{}) {
				if listener := l.(PlayerEventListener); listener != nil {
					listener.OnTrackException(player, track, ee.Exception)
				}
			})

		case TrackStuckEvent:
			player.EmitEvent(func(l interface{}) {
				if listener := l.(PlayerEventListener); listener != nil {
					listener.OnTrackStuck(player, track, ee.ThresholdMs)
				}
			})
		}

	case WebsocketClosedEvent:
		player.EmitEvent(func(l interface{}) {
			if listener := l.(PlayerEventListener); listener != nil {
				listener.OnWebSocketClosed(player, e.Code, e.Reason, e.ByRemote)
			}
		})

	case UnknownEvent:
		for _, pl := range n.Lavalink().Plugins() {
			if plugin, ok := pl.(EventExtension); ok {
				plugin.OnEventInvocation(n, e.Data)
			}
			if plugin, ok := pl.(EventExtensions); ok {
				for _, ext := range plugin.EventExtensions() {
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
	n.stats = stats.Stats
}

func (n *nodeImpl) open(ctx context.Context, delay time.Duration) error {
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
		n.lavalink.Logger().Warnf("error while connecting to lavalink websocket, retrying in %f seconds: %s", delay.Seconds(), err)
		if delay > 0 {
			time.Sleep(delay)
		} else {
			delay = 1 * time.Second
		}
		if delay < 30*time.Second {
			delay *= 2
		}
		return n.open(ctx, delay)
	}
	if n.config.ResumingKey != "" {
		if rs.Header.Get("Session-Resumed") != "true" {
			n.lavalink.Logger().Warnf("failed to resume session with key %s", n.config.ResumingKey)
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
	n.status = Connecting
	defer n.statusMu.Unlock()

	if err := n.open(ctx, 0); err != nil {
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
