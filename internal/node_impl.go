package internal

import (
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
			n.lavalink.Logger().Infof("type: %d, data: %+v", mt, string(data))
			// TODO: handle stuff
		}
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

	return err
}
