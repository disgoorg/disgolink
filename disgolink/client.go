package disgolink

import (
	"context"
	"iter"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

func New(userID snowflake.ID, opts ...ConfigOpt) *Client {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &Client{
		logger:     cfg.Logger,
		httpClient: cfg.HTTPClient,
		userID:     userID,
		nodes:      make(map[string]*Node),
		players:    make(map[snowflake.ID]*Player),
		listeners:  cfg.Listeners,
		plugins:    cfg.Plugins,
	}
}

type Client struct {
	logger     *slog.Logger
	httpClient *http.Client
	userID     snowflake.ID

	nodesMu sync.Mutex
	nodes   map[string]*Node

	playersMu sync.Mutex
	players   map[snowflake.ID]*Player

	listenersMu sync.Mutex
	listeners   []EventListener

	pluginsMu sync.Mutex
	plugins   []Plugin
}

func (c *Client) AddNode(ctx context.Context, config NodeConfig) (*Node, error) {
	node := newNode(c.logger, config, c, c.httpClient)
	if err := node.Open(ctx); err != nil {
		return nil, err
	}

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()
	c.nodes[config.Name] = node
	return node, nil
}

func (c *Client) Nodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		c.nodesMu.Lock()
		defer c.nodesMu.Unlock()

		for _, node := range c.nodes {
			yield(node)
		}
	}
}

func (c *Client) Node(name string) *Node {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	return c.nodes[name]
}

func (c *Client) BestNode() *Node {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	var bestNode *Node
	for _, node := range c.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}

	return bestNode
}

func (c *Client) RemoveNode(name string) {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	if node, ok := c.nodes[name]; ok {
		node.Close()
		delete(c.nodes, name)
	}
}

func (c *Client) Player(guildID snowflake.ID) *Player {
	return c.PlayerOnNode(c.BestNode(), guildID)
}

func (c *Client) PlayerOnNode(node *Node, guildID snowflake.ID) *Player {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	if player, ok := c.players[guildID]; ok {
		return player
	}

	player := newPlayer(c.logger, c, node, guildID)
	for plugin := range c.Plugins() {
		if pl, ok := plugin.(PluginEventHandler); ok {
			pl.OnNewPlayer(player)
		}

	}

	c.players[guildID] = player

	return player
}

func (c *Client) ExistingPlayer(guildID snowflake.ID) *Player {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	return c.players[guildID]
}

func (c *Client) RemovePlayer(guildID snowflake.ID) {
	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	delete(c.players, guildID)
}

func (c *Client) Players() iter.Seq[*Player] {
	return func(yield func(*Player) bool) {
		c.playersMu.Lock()
		defer c.playersMu.Unlock()

		for _, player := range c.players {
			yield(player)
		}
	}
}

func (c *Client) EmitEvent(player *Player, event lavalink.Message) {
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

func (c *Client) AddListeners(listeners ...EventListener) {
	c.listenersMu.Lock()
	defer c.listenersMu.Unlock()

	c.listeners = append(c.listeners, listeners...)
}

func (c *Client) RemoveListeners(listeners ...EventListener) {
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

func (c *Client) AddPlugins(plugins ...Plugin) {
	c.pluginsMu.Lock()
	defer c.pluginsMu.Unlock()

	c.plugins = append(c.plugins, plugins...)
}

func (c *Client) Plugins() iter.Seq[Plugin] {
	return func(yield func(Plugin) bool) {
		c.pluginsMu.Lock()
		defer c.pluginsMu.Unlock()

		for _, plugin := range c.plugins {
			yield(plugin)
		}
	}
}

func (c *Client) RemovePlugins(plugins ...Plugin) {
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

func (c *Client) UserID() snowflake.ID {
	return c.userID
}

func (c *Client) Close() {
	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	for _, node := range c.nodes {
		node.Close()
	}
}

func (c *Client) OnVoiceServerUpdate(ctx context.Context, guildID snowflake.ID, token string, endpoint string) {
	c.Player(guildID).OnVoiceServerUpdate(ctx, token, endpoint)
}

func (c *Client) OnVoiceStateUpdate(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, sessionID string) {
	c.Player(guildID).OnVoiceStateUpdate(ctx, channelID, sessionID)
}
