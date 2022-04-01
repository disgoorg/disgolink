package lavalink

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake"
)

var ErrNoUserID = errors.New("no user id has been configured")

type Lavalink interface {
	Logger() log.Logger

	AddNode(ctx context.Context, config NodeConfig) (Node, error)
	Nodes() []Node
	Node(name string) Node
	BestNode() Node
	BestRestClient() RestClient
	RemoveNode(name string)

	AddPlugins(plugins ...interface{})
	Plugins() []interface{}
	RemovePlugins(plugins ...interface{})

	EncodeTrack(track AudioTrack) (string, error)
	DecodeTrack(track string) (AudioTrack, error)

	Player(guildID snowflake.Snowflake) Player
	PlayerOnNode(name string, guildID snowflake.Snowflake) Player
	RestorePlayer(restoreState PlayerRestoreState) (Player, error)
	ExistingPlayer(guildID snowflake.Snowflake) Player
	RemovePlayer(guildID snowflake.Snowflake)
	Players() map[snowflake.Snowflake]Player

	UserID() snowflake.Snowflake
	SetUserID(userID snowflake.Snowflake)

	Close()

	OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)
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
	return &lavalinkImpl{
		config:  *config,
		nodes:   map[string]Node{},
		players: map[snowflake.Snowflake]Player{},
	}
}

var _ Lavalink = (*lavalinkImpl)(nil)

type lavalinkImpl struct {
	config    Config
	pluginsMu sync.Mutex

	nodesMu sync.Mutex
	nodes   map[string]Node

	playersMu sync.Mutex
	players   map[snowflake.Snowflake]Player
}

func (l *lavalinkImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *lavalinkImpl) AddNode(ctx context.Context, config NodeConfig) (Node, error) {
	if l.UserID() == "" {
		return nil, ErrNoUserID
	}
	node := &nodeImpl{
		config:   config,
		lavalink: l,
		status:   Disconnected,
	}
	node.restClient = newRestClientImpl(node, l.config.HTTPClient)
	if err := node.Open(ctx); err != nil {
		return nil, err
	}

	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	l.nodes[config.Name] = node
	return node, nil
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
		if bestNode == nil || (node.Stats() != nil && bestNode.Stats() != nil && node.Stats().Better(*bestNode.Stats())) {
			bestNode = node
		}
	}
	return bestNode
}

func (l *lavalinkImpl) BestRestClient() RestClient {
	if node := l.BestNode(); node != nil {
		return node.RestClient()
	}
	return nil
}

func (l *lavalinkImpl) RemoveNode(name string) {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	if node, ok := l.nodes[name]; ok {
		node.Close()
		delete(l.nodes, name)
	}
}

func (l *lavalinkImpl) AddPlugins(plugins ...interface{}) {
	l.pluginsMu.Lock()
	defer l.pluginsMu.Unlock()
	for _, plugin := range plugins {
		l.config.Plugins = append(l.config.Plugins, plugin)
	}
}

func (l *lavalinkImpl) Plugins() []interface{} {
	l.pluginsMu.Lock()
	defer l.pluginsMu.Unlock()
	plugins := make([]interface{}, len(l.config.Plugins))
	i := 0
	for _, plugin := range l.config.Plugins {
		plugins[i] = plugin
		i++
	}
	return plugins
}

func (l *lavalinkImpl) RemovePlugins(plugins ...interface{}) {
	l.pluginsMu.Lock()
	defer l.pluginsMu.Unlock()
	for _, plugin := range plugins {
		for i, p := range l.config.Plugins {
			if p == plugin {
				l.config.Plugins = append(l.config.Plugins[:i], l.config.Plugins[i+1:]...)
				break
			}
		}
	}
}

func (l *lavalinkImpl) EncodeTrack(track AudioTrack) (string, error) {
	return EncodeToString(track, func(track AudioTrack, w io.Writer) error {
		for _, pl := range l.Plugins() {
			if plugin, ok := pl.(SourceExtension); ok {
				if plugin.SourceName() == track.Info().SourceName {
					return plugin.Encode(track, w)
				}
			}
		}
		return nil
	})
}

func (l *lavalinkImpl) DecodeTrack(str string) (AudioTrack, error) {
	return DecodeString(str, func(info AudioTrackInfo, r io.Reader) (AudioTrack, error) {
		for _, pl := range l.Plugins() {
			if plugin, ok := pl.(SourceExtension); ok {
				if plugin.SourceName() == info.SourceName {
					return plugin.Decode(info, r)
				}
			}
		}
		return nil, nil
	})
}

func (l *lavalinkImpl) Player(guildID snowflake.Snowflake) Player {
	return l.PlayerOnNode("", guildID)
}

func (l *lavalinkImpl) PlayerOnNode(name string, guildID snowflake.Snowflake) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	if player, ok := l.players[guildID]; ok {
		return player
	}
	node := l.Node(name)
	if node == nil {
		node = l.BestNode()
	}
	player := NewPlayer(node, l, guildID)
	for _, pl := range l.config.Plugins {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNewPlayer(player)
		}
	}
	l.players[guildID] = player
	return player
}

func (l *lavalinkImpl) RestorePlayer(restoreState PlayerRestoreState) (Player, error) {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	node := l.Node(restoreState.Node)
	if node == nil {
		node = l.BestNode()
	}
	player, err := newResumingPlayer(node, l, restoreState)
	if err != nil {
		return nil, err
	}
	for _, pl := range l.config.Plugins {
		if plugin, ok := pl.(PluginEventHandler); ok {
			plugin.OnNewPlayer(player)
		}
	}
	l.players[restoreState.GuildID] = player
	return player, nil
}

func (l *lavalinkImpl) ExistingPlayer(guildID snowflake.Snowflake) Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	return l.players[guildID]
}

func (l *lavalinkImpl) RemovePlayer(guildID snowflake.Snowflake) {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	delete(l.players, guildID)
}

func (l *lavalinkImpl) Players() map[snowflake.Snowflake]Player {
	l.playersMu.Lock()
	defer l.playersMu.Unlock()
	players := make(map[snowflake.Snowflake]Player, len(l.players))
	for guildID, player := range l.players {
		players[guildID] = player
	}
	return players
}

func (l *lavalinkImpl) UserID() snowflake.Snowflake {
	return l.config.UserID
}

func (l *lavalinkImpl) SetUserID(userID snowflake.Snowflake) {
	l.config.UserID = userID
}

func (l *lavalinkImpl) Close() {
	l.nodesMu.Lock()
	defer l.nodesMu.Unlock()
	for _, node := range l.nodes {
		node.Close()
	}
}

func (l *lavalinkImpl) OnVoiceServerUpdate(voiceServerUpdate VoiceServerUpdate) {
	player := l.ExistingPlayer(voiceServerUpdate.GuildID)
	if player == nil {
		return
	}
	player.OnVoiceServerUpdate(voiceServerUpdate)
}

func (l *lavalinkImpl) OnVoiceStateUpdate(voiceStateUpdate VoiceStateUpdate) {
	player := l.ExistingPlayer(voiceStateUpdate.GuildID)
	if player == nil {
		return
	}
	player.OnVoiceStateUpdate(voiceStateUpdate)
}
