package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/DisgoOrg/disgolink/api"
	"github.com/gorilla/websocket"
)

type NodeImpl struct {
	options    *api.NodeOptions
	lavalink   api.Lavalink
	conn       *websocket.Conn
	quit       chan interface{}
	status     api.NodeStatus
	stats      *api.Stats
	available  bool
	restClient api.RestClient
}

func (n NodeImpl) RestURL() string {
	scheme := "http"
	if n.options.Secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s:%d", scheme, n.options.Host, n.options.Port)
}

func (n *NodeImpl) Lavalink() api.Lavalink {
	return n.lavalink
}

func (n *NodeImpl) Options() *api.NodeOptions {
	return n.options
}

func (n *NodeImpl) RestClient() api.RestClient {
	return n.restClient
}

func (n *NodeImpl) Name() string {
	return n.options.Name
}

func (n *NodeImpl) Send(d interface{}) {
	err := n.conn.WriteJSON(d)
	if err != nil {
		log.Println(err)
	}
}

func (n *NodeImpl) Status() api.NodeStatus {
	return n.status
}

func (n *NodeImpl) reconnect(delay time.Duration) {
	go func() {
		time.Sleep(delay)

		if n.Status() == api.Connecting || n.Status() == api.Reconnecting {
			n.lavalink.Logger().Error("tried to reconnect gateway while connecting/reconnecting")
			return
		}
		n.lavalink.Logger().Info("reconnecting gateway...")
		if err := n.Open(); err != nil {
			n.lavalink.Logger().Errorf("failed to reconnect gateway: %s", err)
			n.status = api.Disconnected
			n.reconnect(delay * 2)
		}
	}()
}

func (n *NodeImpl) Close() {
	n.status = api.Disconnected
	if n.quit != nil {
		n.lavalink.Logger().Info("closing ws goroutines...")
		close(n.quit)
		n.quit = nil
		n.lavalink.Logger().Info("closed ws goroutines")
	}
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			n.lavalink.Logger().Errorf("error while closing wsconn: %s", err)
		}
	}
}

func (n *NodeImpl) listen() {
	defer func() {
		n.lavalink.Logger().Info("shut down listen goroutine")
	}()
	for {
		select {
		case <-n.quit:
			n.lavalink.Logger().Infof("existed listen routine")
			return
		default:
			if n.conn == nil {
				return
			}
			mt, data, err := n.conn.ReadMessage()
			if err != nil {
				n.lavalink.Logger().Errorf("error while reading from ws. error: %s", err)
				n.Close()
				n.reconnect(1 * time.Second)
				return
			}
			op, err := n.getOp(mt, data)
			if err != nil {
				n.lavalink.Logger().Errorf("error getting op from websocket message: %s", err)
				continue
			}
			switch *op {
			case api.OpPlayerUpdate:
				n.onPlayerUpdate(data)
			case api.OpEvent:
				n.onTrackEvent(data)
			case api.OpStats:
			default:
				n.lavalink.Logger().Warnf("unexpected op received: %s", op)
			}
		}
	}
}

func (n *NodeImpl) getOp(mt int, data []byte) (*api.Op, error) {
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("recieved unexpected mt type: %d", mt)
	}

	var op api.GenericOp
	if err := json.Unmarshal(data, &op); err != nil {
		return nil, err
	}
	return &op.Op, nil
}

func (n *NodeImpl) onPlayerUpdate(data []byte) {
	var playerUpdate api.PlayerUpdateEvent
	err := json.Unmarshal(data, &playerUpdate)
	if err != nil {
		n.lavalink.Logger().Errorf("error unmarshalling PlayerUpdateEvent: %s", err)
		return
	}
	if player := n.lavalink.ExistingPlayer(playerUpdate.GuildID); player != nil {
		player.PlayerUpdate(playerUpdate.State)
		player.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnPlayerUpdate(player, playerUpdate.State)
		})
	}
}

func (n *NodeImpl) onTrackEvent(data []byte) {
	var event api.GenericPlayerEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		n.lavalink.Logger().Errorf("error unmarshalling GenericPlayerEvent: %s", err)
		return
	}
	p := n.lavalink.ExistingPlayer(event.GuildID)
	if p == nil {
		return
	}

	switch event.Type {
	case api.WebsocketEventTrackStart:
		var trackStartEvent api.TrackStartEvent
		if err = json.Unmarshal(data, &trackStartEvent); err != nil {
			n.lavalink.Logger().Errorf("error unmarshalling TrackStartEvent: %s", err)
			return
		}
		track := trackStartEvent.Track()
		p.SetTrack(track)
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnTrackStart(p, track)
		})
	case api.WebsocketEventTrackEnd:
		var trackEndEvent api.TrackEndEvent
		if err = json.Unmarshal(data, &trackEndEvent); err != nil {
			n.lavalink.Logger().Errorf("error unmarshalling TrackEndEvent: %s", err)
			return
		}
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnTrackEnd(p, trackEndEvent.Track(), trackEndEvent.EndReason)
		})
	case api.WebsocketEventTrackException:
		var trackExceptionEvent api.TrackExceptionEvent
		if err = json.Unmarshal(data, &trackExceptionEvent); err != nil {
			n.lavalink.Logger().Errorf("error unmarshalling TrackExceptionEvent: %s", err)
			return
		}
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnTrackException(p, trackExceptionEvent.Track(), trackExceptionEvent.Exception)
		})
	case api.WebsocketEventTrackStuck:
		var trackStuckEvent api.TrackStuckEvent
		if err = json.Unmarshal(data, &trackStuckEvent); err != nil {
			n.lavalink.Logger().Errorf("error unmarshalling TrackStuckEvent: %s", err)
			return
		}
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnTrackStuck(p, trackStuckEvent.Track(), trackStuckEvent.ThresholdMs)
		})
	case api.WebSocketEventClosed:
		var websocketClosedEvent api.WebSocketClosedEvent
		if err = json.Unmarshal(data, &websocketClosedEvent); err != nil {
			n.lavalink.Logger().Errorf("error unmarshalling WebSocketClosedEvent: %s", err)
			return
		}
		p.EmitEvent(func(listener api.PlayerEventListener) {
			listener.OnWebSocketClosed(p, websocketClosedEvent.Code, websocketClosedEvent.Reason, websocketClosedEvent.ByRemote)
		})
	default:
		n.lavalink.Logger().Warnf("unexpected event received: %s", string(data))
		return
	}
}

func (n *NodeImpl) Open() error {
	scheme := "ws"
	if n.options.Secure {
		scheme += "s"
	}
	header := http.Header{}
	header.Add("Authorization", n.options.Password)
	header.Add("User-Id", n.lavalink.UserID())
	header.Add("Client-Name", n.lavalink.ClientName())
	u := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%v:%v", n.options.Host, n.options.Port),
	}

	var err error
	n.conn, _, err = websocket.DefaultDialer.Dial(u.String(), header)

	go n.listen()

	return err
}
