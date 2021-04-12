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
	api.NodeOptions
	lavalink  api.Lavalink
	conn      *websocket.Conn
	quit      chan interface{}
	status    api.NodeStatus
	stats     *api.Stats
	available bool
}

func (n *NodeImpl) Name() string {
	return n.NodeOptions.Name
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

func (n *NodeImpl) name() {
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
	if n.Secure {
		scheme += "s"
	}
	header := http.Header{}
	header.Add("Authorization", n.Password)
	header.Add("User-Id", n.lavalink.UserID().String())
	header.Add("Client-Name", n.lavalink.ClientName())
	u := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%v:%v", n.Host, n.Port),
	}

	var err error
	n.conn, _, err = websocket.DefaultDialer.Dial(u.String(), header)

	return err
}
