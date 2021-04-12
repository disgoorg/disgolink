package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/DisgoOrg/disgolink/api"
	"github.com/gorilla/websocket"
)

type NodeImpl struct {
	api.NodeOptions
	lavalink  api.Lavalink
	conn      *websocket.Conn
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
