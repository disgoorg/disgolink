package disgolink

import (
	"context"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/disgoorg/disgolink/v2/lavalink"

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

	AddPlugins(plugins ...Plugin)
	ForPlugins(pluginFunc func(plugin Plugin))
	RemovePlugins(plugins ...Plugin)

	Close()

	OnVoiceServerUpdate(ctx context.Context, guildID snowflake.ID, token string, endpoint string)
	OnVoiceStateUpdate(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, sessionID string)
}

func New(opts ...ConfigOpt) Client {
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
		nodes:      map[string]Node{},
		players:    map[snowflake.ID]Player{},
		listeners:  config.Listeners,
		plugins:    config.Plugins,
	}
}

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	logger     log.Logger
	httpClient *http.Client

	nodesMu sync.Mutex
	nodes   map[string]Node

	playersMu sync.Mutex
	players   map[snowflake.ID]Player

	listenersMu sync.Mutex
	listeners   []EventListener

	pluginsMu sync.Mutex
	plugins   []Plugin
}

func (c *clientImpl) Logger() log.Logger {
	return c.logger
}

func (c *clientImpl) AddNode(ctx context.Context, config NodeConfig) (Node, error) {
	node := &nodeImpl{
		config:   config,
		lavalink: c,
		status:   StatusDisconnected,
	}
	node.rest = newRestClientImpl(node, c.httpClient)
	if err := node.Open(ctx); err != nil {
		return nil, err
	}

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	c.nodes[config.Name] = node
	return node, nil
}

func (c *clientImpl) ForNodes(nodeFunc func(node Node)) {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	for i := range c.nodes {
		nodeFunc(c.nodes[i])
	}
}

func (c *clientImpl) Node(name string) Node {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	return c.nodes[name]
}

func (c *clientImpl) BestNode() Node {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	var bestNode Node
	for _, node := range c.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}
	return bestNode
}

func (c *clientImpl) RemoveNode(name string) {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	if node, ok := c.nodes[name]; ok {
		node.Close()
		delete(c.nodes, name)
	}
}

func (c *clientImpl) Player(guildID snowflake.ID) Player {
	return c.PlayerOnNode("", guildID)
}

func (c *clientImpl) PlayerOnNode(name string, guildID snowflake.ID) Player {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()
	if player, ok := c.players[guildID]; ok {
		return player
	}
	node := c.Node(name)
	if node == nil {
		node = c.BestNode()
	}

	player := NewPlayer(c, node, guildID)
	c.ForPlugins(func(plugin Plugin) {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNewPlayer(player)
		}
	})
	c.players[guildID] = player
	return player
}

func (c *clientImpl) ExistingPlayer(guildID snowflake.ID) Player {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()
	return c.players[guildID]
}

func (c *clientImpl) RemovePlayer(guildID snowflake.ID) {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()
	delete(c.players, guildID)
}

func (c *clientImpl) ForPlayers(playerFunc func(player Player)) {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()
	for _, player := range c.players {
		playerFunc(player)
	}
}

func (c *clientImpl) EmitEvent(player Player, event lavalink.Event) {
	c.listenersMu.Lock()
	defer c.listenersMu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			c.Logger().Errorf("recovered from panic in event listener: %#v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	for _, listener := range c.listeners {
		listener.OnEvent(player, event)
	}
}

func (c *clientImpl) AddListeners(listeners ...EventListener) {
	c.listenersMu.Lock()
	defer c.listenersMu.Unlock()
	c.listeners = append(c.listeners, listeners...)
}

func (c *clientImpl) RemoveListeners(listeners ...EventListener) {
	c.listenersMu.Lock()
	defer c.listenersMu.Unlock()
	for _, listener := range listeners {
		for i, ln := range c.listeners {
			if ln == listener {
				c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
			}
		}
	}
}

func (c *clientImpl) AddPlugins(plugins ...Plugin) {
	c.pluginsMu.Lock()
	defer c.pluginsMu.Unlock()
	c.plugins = append(c.plugins, plugins...)
}

func (c *clientImpl) ForPlugins(pluginFunc func(plugin Plugin)) {
	c.pluginsMu.Lock()
	defer c.pluginsMu.Unlock()
	for _, plugin := range c.plugins {
		pluginFunc(plugin)
	}
}

func (c *clientImpl) RemovePlugins(plugins ...Plugin) {
	c.pluginsMu.Lock()
	defer c.pluginsMu.Unlock()
	for _, plugin := range plugins {
		for i, pl := range c.plugins {
			if pl == plugin {
				c.plugins = append(c.plugins[:i], c.plugins[i+1:]...)
			}
		}
	}
}

func (c *clientImpl) Close() {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	for _, node := range c.nodes {
		node.Close()
	}
}

func (c *clientImpl) OnVoiceServerUpdate(ctx context.Context, guildID snowflake.ID, token string, endpoint string) {
	player := c.ExistingPlayer(guildID)
	if player == nil {
		return
	}
	player.OnVoiceServerUpdate(ctx, token, endpoint)
}

func (c *clientImpl) OnVoiceStateUpdate(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, sessionID string) {
	player := c.ExistingPlayer(guildID)
	if player == nil {
		return
	}
	player.OnVoiceStateUpdate(ctx, channelID, sessionID)
}
