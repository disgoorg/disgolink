package lavalink

import (
	"context"
	"github.com/DisgoOrg/log"
	"net/http"
	"time"
)

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
		config:  *config,
		nodes:   map[string]Node{},
		players: map[string]Player{},
	}
	return lavalink
}

type Lavalink interface {
	Logger() log.Logger

	AddNode(config NodeConfig)
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

var _ Lavalink = (*lavalinkImpl)(nil)

type lavalinkImpl struct {
	config  Config
	nodes   map[string]Node
	players map[string]Player
}

func (l *lavalinkImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *lavalinkImpl) AddNode(config NodeConfig) {
	node := &nodeImpl{
		config:   config,
		lavalink: l,
	}
	node.restClient = newRestClientImpl(node, l.config.HTTPClient)

	l.nodes[config.Name] = node
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
}

func (l *lavalinkImpl) Node(name string) Node {
	return l.nodes[name]
}

func (l *lavalinkImpl) BestNode() Node {
	var bestNode Node
	for _, node := range l.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}
	return bestNode
}

func (l lavalinkImpl) BestRestClient() RestClient {
	if len(l.nodes) == 0 {
		return nil
	}
	return l.BestNode().RestClient()
}

func (l *lavalinkImpl) RemoveNode(name string) {
	delete(l.nodes, name)
}

func (l *lavalinkImpl) AddPlugins(plugins ...Plugin) {
	for _, plugin := range plugins {
		l.config.Plugins = append(l.config.Plugins, plugin)
	}
}

func (l *lavalinkImpl) Plugins() []Plugin {
	return l.config.Plugins
}

func (l *lavalinkImpl) RemovePlugins(plugins ...Plugin) {
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
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.BestNode(), guildID)
	l.players[guildID] = player
	return player
}

func (l *lavalinkImpl) ExistingPlayer(guildID string) Player {
	return l.players[guildID]
}

func (l *lavalinkImpl) Players() map[string]Player {
	return l.players
}

func (l *lavalinkImpl) UserID() string {
	return l.config.UserID
}

func (l *lavalinkImpl) SetUserID(userID string) {
	l.config.UserID = userID
}

func (l *lavalinkImpl) Close(ctx context.Context) {
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
