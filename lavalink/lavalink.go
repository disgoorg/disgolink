package lavalink

import (
	"context"
	"github.com/DisgoOrg/log"
	"net/http"
	"sync"
	"time"
)

type Lavalink interface {
	Logger() log.Logger

	AddNode(config NodeConfig) Node
	Nodes() []Node
	Node(name string) Node
	BestNode() Node
	BestRestClient() RestClient
	RemoveNode(name string)

	AddPlugins(plugins ...Plugin)
	Plugins() []Plugin
	RemovePlugins(plugins ...Plugin)

	Player(guildID string) Player
	ExistingPlayer(guildID string) Player
	Players() map[string]Player

	UserID() string
	SetUserID(userID string)

	Close(ctx context.Context)

	VoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)
}

func New(opts ...ConfigOpt) Lavalink {
	config := &Config{}
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 20 * time.Second}
	}
	lavalink := &lavalinkImpl{
		config:    *config,
		pluginsMu: &sync.Mutex{},
		nodesMu:   &sync.Mutex{},
		nodes:     map[string]Node{},
		playersMu: &sync.Mutex{},
		players:   map[string]Player{},
	}
	return lavalink
}

var _ Lavalink = (*lavalinkImpl)(nil)

type lavalinkImpl struct {
	config    Config
	pluginsMu sync.Locker

	nodesMu sync.Locker
	nodes   map[string]Node

	playersMu sync.Locker
	players   map[string]Player
}

func (l *lavalinkImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *lavalinkImpl) AddNode(config NodeConfig) Node {
	node := &nodeImpl{
		config:   config,
		lavalink: l,
	}
	node.restClient = newRestClientImpl(node, l.config.HTTPClient)
	go func() {
		delay := 500
		for {
			err := node.Open(context.TODO())
			if err == nil {
				break
			}
			delay += int(float64(delay) * 1.2)
			l.Logger().Errorf("error while connecting to node: %s, waiting %ds, error: %s", node.Name(), delay/1000, err)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}()
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	l.nodes[config.Name] = node
	return node
}

func (l *lavalinkImpl) Nodes() []Node {
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

func (l *lavalinkImpl) Node(name string) Node {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	return l.nodes[name]
}

func (l *lavalinkImpl) BestNode() Node {
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

func (l lavalinkImpl) BestRestClient() RestClient {
	if node := l.BestNode(); node != nil {
		return node.RestClient()
	}
	return nil
}

func (l *lavalinkImpl) RemoveNode(name string) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	delete(l.nodes, name)
}

func (l *lavalinkImpl) AddPlugins(plugins ...Plugin) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for _, plugin := range plugins {
		l.config.Plugins = append(l.config.Plugins, plugin)
	}
}

func (l *lavalinkImpl) Plugins() []Plugin {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	plugins := make([]Plugin, len(l.config.Plugins))
	i := 0
	for _, plugin := range l.config.Plugins {
		plugins[i] = plugin
		i++
	}
	return plugins
}

func (l *lavalinkImpl) RemovePlugins(plugins ...Plugin) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for _, plugin := range plugins {
		for i, p := range l.config.Plugins {
			if p == plugin {
				l.config.Plugins = append(l.config.Plugins[:i], l.config.Plugins[i+1:]...)
				break
			}
		}
	}
}

func (l *lavalinkImpl) Player(guildID string) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.BestNode(), guildID)
	for _, pl := range l.config.Plugins {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNewPlayer(player)
		}
	}
	l.players[guildID] = player
	return player
}

func (l *lavalinkImpl) ExistingPlayer(guildID string) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	return l.players[guildID]
}

func (l *lavalinkImpl) Players() map[string]Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	players := make(map[string]Player, len(l.players))
	for guildID, player := range l.players {
		players[guildID] = player
	}
	return players
}

func (l *lavalinkImpl) UserID() string {
	return l.config.UserID
}

func (l *lavalinkImpl) SetUserID(userID string) {
	l.config.UserID = userID
}

func (l *lavalinkImpl) Close(ctx context.Context) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for _, node := range l.nodes {
		node.Close(ctx)
	}
}

func (l *lavalinkImpl) VoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	player := l.players[voiceServerUpdate.GuildID]
	if player == nil && player.LastSessionID() != nil {
		return
	}
	_ = player.Node().Send(VoiceUpdateCommand{
		GuildID:   voiceServerUpdate.GuildID,
		SessionID: *player.LastSessionID(),
		Event:     voiceServerUpdate,
	})
}

func (l *lavalinkImpl) VoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
	player := l.players[voiceStateUpdate.GuildID]
	if player == nil {
		return
	}
	player.SetChannelID(voiceStateUpdate.ChannelID)
	player.SetLastSessionID(voiceStateUpdate.SessionID)
}
