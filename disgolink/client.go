package disgolink

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Client interface {
	AddNode(ctx context.Context, config NodeConfig) (Node, error)
	ForNodes(nodeFunc func(node Node))
	Node(name string) Node
	BestNode() Node
	RemoveNode(name string)

	Player(guildID snowflake.ID) Player
	PlayerOnNode(node Node, guildID snowflake.ID) Player
	ExistingPlayer(guildID snowflake.ID) Player
	RemovePlayer(guildID snowflake.ID)
	ForPlayers(playerFunc func(player Player))

	EmitEvent(player Player, event lavalink.Message)
	AddListeners(listeners ...EventListener)
	RemoveListeners(listeners ...EventListener)

	AddPlugins(plugins ...Plugin)
	ForPlugins(pluginFunc func(plugin Plugin))
	RemovePlugins(plugins ...Plugin)

	UserID() snowflake.ID
	Close()

	OnVoiceServerUpdate(ctx context.Context, guildID snowflake.ID, token string, endpoint string)
	OnVoiceStateUpdate(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, sessionID string)
}

func New(userID snowflake.ID, opts ...ConfigOpt) Client {
	cfg := DefaultConfig()
	cfg.Apply(opts)
	cfg.Logger = cfg.Logger.With(slog.String("name", "disgolink_client"))

	return &clientImpl{
		logger:     cfg.Logger,
		httpClient: cfg.HTTPClient,
		userID:     userID,
		nodes:      map[string]Node{},
		players:    map[snowflake.ID]Player{},
		listeners:  cfg.Listeners,
		plugins:    cfg.Plugins,
	}
}

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	logger     *slog.Logger
	httpClient *http.Client
	userID     snowflake.ID

	nodesMu sync.Mutex
	nodes   map[string]Node

	playersMu sync.Mutex
	players   map[snowflake.ID]Player

	listenersMu sync.Mutex
	listeners   []EventListener

	pluginsMu sync.Mutex
	plugins   []Plugin
}

func (c *clientImpl) AddNode(ctx context.Context, config NodeConfig) (Node, error) {
	node := &nodeImpl{
		logger:   c.logger.With(slog.String("name", "disgolink_node"), slog.String("node_name", config.Name)),
		config:   config,
		lavalink: c,
		status:   StatusDisconnected,
	}
	node.rest = &restClientImpl{
		logger:     c.logger.With(slog.String("name", "disgolink_rest_client"), slog.String("node_name", config.Name)),
		node:       node,
		httpClient: c.httpClient,
	}
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
	return c.PlayerOnNode(c.BestNode(), guildID)
}

func (c *clientImpl) PlayerOnNode(node Node, guildID snowflake.ID) Player {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()
	if player, ok := c.players[guildID]; ok {
		return player
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

func (c *clientImpl) EmitEvent(player Player, event lavalink.Message) {
	c.listenersMu.Lock()
	defer c.listenersMu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("recovered from panic in event listener", slog.Any("r", r), slog.String("stack", string(debug.Stack())))
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

func (c *clientImpl) UserID() snowflake.ID {
	return c.userID
}

func (c *clientImpl) Close() {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	for _, node := range c.nodes {
		node.Close()
	}
}

func (c *clientImpl) OnVoiceServerUpdate(ctx context.Context, guildID snowflake.ID, token string, endpoint string) {
	c.Player(guildID).OnVoiceServerUpdate(ctx, token, endpoint)
}

func (c *clientImpl) OnVoiceStateUpdate(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, sessionID string) {
	c.Player(guildID).OnVoiceStateUpdate(ctx, channelID, sessionID)
}
