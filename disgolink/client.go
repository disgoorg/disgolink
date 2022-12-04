package disgolink

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var ErrNoUserID = errors.New("no user id has been configured")

type Client interface {
	Logger() log.Logger

	AddNode(ctx context.Context, config NodeConfig) (Node, error)
	Nodes() []Node
	Node(name string) Node
	BestNode() Node
	RemoveNode(name string)

	Player(guildID snowflake.ID) Player
	PlayerOnNode(name string, guildID snowflake.ID) Player
	ExistingPlayer(guildID snowflake.ID) Player
	RemovePlayer(guildID snowflake.ID)
	ForPlayers(playerFunc func(player Player))

	UserID() snowflake.ID
	Close()

	OnVoiceServerUpdate(guildID snowflake.ID, token string, endpoint string)
	OnVoiceStateUpdate(guildID snowflake.ID, channelID *snowflake.ID, sessionID string)
}

func New(userID snowflake.ID, opts ...ConfigOpt) (Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	if userID == 0 {
		return nil, ErrNoUserID
	}

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 20 * time.Second}
	}
	return &clientImpl{
		config:  *config,
		userID:  userID,
		nodes:   map[string]Node{},
		players: map[snowflake.ID]Player{},
	}, nil
}

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	config Config
	userID snowflake.ID

	nodesMu sync.Mutex
	nodes   map[string]Node

	playersMu sync.Mutex
	players   map[snowflake.ID]Player
}

func (l *clientImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *clientImpl) AddNode(ctx context.Context, config NodeConfig) (Node, error) {
	node := &nodeImpl{
		config:   config,
		lavalink: l,
		status:   StatusDisconnected,
	}
	node.rest = newRestClientImpl(node, l.config.HTTPClient)
	if err := node.Open(ctx); err != nil {
		return nil, err
	}

	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	l.nodes[config.Name] = node
	return node, nil
}

func (l *clientImpl) Nodes() []Node {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	nodes := make([]Node, len(l.nodes))
	i := 0
	for _, node := range l.nodes {
		nodes[i] = node
		i++
	}
	return nodes
}

func (l *clientImpl) Node(name string) Node {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	return l.nodes[name]
}

func (l *clientImpl) BestNode() Node {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	var bestNode Node
	for _, node := range l.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}
	return bestNode
}

func (l *clientImpl) RemoveNode(name string) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	if node, ok := l.nodes[name]; ok {
		node.Close()
		delete(l.nodes, name)
	}
}

func (l *clientImpl) Player(guildID snowflake.ID) Player {
	return l.PlayerOnNode("", guildID)
}

func (l *clientImpl) PlayerOnNode(name string, guildID snowflake.ID) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	if player, ok := l.players[guildID]; ok {
		return player
	}
	node := l.Node(name)
	if node == nil {
		node = l.BestNode()
	}

	player := NewPlayer(l, node, guildID)
	l.players[guildID] = player
	return player
}

func (l *clientImpl) ExistingPlayer(guildID snowflake.ID) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	return l.players[guildID]
}

func (l *clientImpl) RemovePlayer(guildID snowflake.ID) {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	delete(l.players, guildID)
}

func (l *clientImpl) ForPlayers(playerFunc func(player Player)) {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	for _, player := range l.players {
		playerFunc(player)
	}
}

func (l *clientImpl) UserID() snowflake.ID {
	return l.userID
}

func (l *clientImpl) Close() {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for _, node := range l.nodes {
		node.Close()
	}
}

func (l *clientImpl) OnVoiceServerUpdate(guildID snowflake.ID, token string, endpoint string) {
	player := l.ExistingPlayer(guildID)
	if player == nil {
		return
	}
	player.OnVoiceServerUpdate(token, endpoint)
}

func (l *clientImpl) OnVoiceStateUpdate(guildID snowflake.ID, channelID *snowflake.ID, sessionID string) {
	player := l.ExistingPlayer(guildID)
	if player == nil {
		return
	}
	player.OnVoiceStateUpdate(channelID, sessionID)
}
