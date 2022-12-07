package disgolink

import (
	"context"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type Client interface {
	Logger() log.Logger

	AddNode(ctx context.Context, config NodeConfig) (Node, error)
	ForNodes(nodeFunc func(node Node))
	Node(name string) Node
	BestNode() Node
	RemoveNode(name string)

	Player(guildID snowflake.ID) Player
	PlayerOnNode(name string, guildID snowflake.ID) Player
	ExistingPlayer(guildID snowflake.ID) Player
	RemovePlayer(guildID snowflake.ID)
	ForPlayers(playerFunc func(player Player))

	EmitEvent(player Player, event lavalink.Event)
	AddListeners(listeners ...EventListener)
	RemoveListeners(listeners ...EventListener)

	UserID() snowflake.ID
	Close()

	OnVoiceServerUpdate(guildID snowflake.ID, token string, endpoint string)
	OnVoiceStateUpdate(guildID snowflake.ID, channelID *snowflake.ID, sessionID string)
}

func New(userID snowflake.ID, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 20 * time.Second}
	}
	return &clientImpl{
		logger:     config.Logger,
		httpClient: config.HTTPClient,
		userID:     userID,
		nodes:      map[string]Node{},
		players:    map[snowflake.ID]Player{},
		listeners:  config.Listeners,
	}
}

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	logger     log.Logger
	httpClient *http.Client
	userID     snowflake.ID

	nodesMu sync.Mutex
	nodes   map[string]Node

	playersMu sync.Mutex
	players   map[snowflake.ID]Player

	listenersMu sync.Mutex
	listeners   []EventListener
}

func (l *clientImpl) Logger() log.Logger {
	return l.logger
}

func (l *clientImpl) AddNode(ctx context.Context, config NodeConfig) (Node, error) {
	node := &nodeImpl{
		config:   config,
		lavalink: l,
		status:   StatusDisconnected,
	}
	node.rest = newRestClientImpl(node, l.httpClient)
	if err := node.Open(ctx); err != nil {
		return nil, err
	}

	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	l.nodes[config.Name] = node
	return node, nil
}

func (l *clientImpl) ForNodes(nodeFunc func(node Node)) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for i := range l.nodes {
		nodeFunc(l.nodes[i])
	}
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

func (l *clientImpl) EmitEvent(player Player, event lavalink.Event) {
	l.listenersMu.Lock()
	defer l.listenersMu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			l.Logger().Errorf("recovered from panic in event listener: %#v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	for _, listener := range l.listeners {
		listener.OnEvent(player, event)
	}
}
func (l *clientImpl) AddListeners(listeners ...EventListener) {
	l.listenersMu.Lock()
	defer l.listenersMu.Unlock()
	l.listeners = append(l.listeners, listeners...)
}
func (l *clientImpl) RemoveListeners(listeners ...EventListener) {
	l.listenersMu.Lock()
	defer l.listenersMu.Unlock()
	for _, listener := range listeners {
		for i, ln := range l.listeners {
			if ln == listener {
				l.listeners = append(l.listeners[:i], l.listeners[i+1:]...)
			}
		}
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
