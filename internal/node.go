package internal

import (
	"fmt"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type NodeOptions struct {
	Host     string
	Port     int
	Password string
	Secure   bool
}

type Node struct {
	options *NodeOptions
	manager *Manager
	conn    *websocket.Conn
	Stats   *api.Stats
}

func CreateNode(manager Manager, options NodeOptions) *Node {
	return &Node{
		options: &options,
		manager: &manager,
	}
}

func (n *Node) Send(d interface{}) error {
	return n.conn.WriteJSON(d)
}

func (n *Node) Open() error {
	scheme := "ws"
	if n.options.Secure {
		scheme += "s"
	}
	header := http.Header{}
	header.Add("Authorization", n.options.Password)
	header.Add("User-Id", *n.manager.Options.ClientID)
	header.Add("Client-Name", *n.manager.Options.ClientName)
	u := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%v:%v", n.options.Host, n.options.Port),
	}

	var err error
	n.conn, _, err = websocket.DefaultDialer.Dial(u.String(), header)

	return err
}
