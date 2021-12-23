package lavalink

import (
	"context"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
	"net/http"
	"time"
)

func NewLavalink(opts ...ConfigOpt) Lavalink {
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
		players: map[discord.Snowflake]Player{},
	}
	return lavalink
}

type Lavalink interface {
	Logger() log.Logger

	AddPlugins(plugins ...interface{})
	RemovePlugins(plugins ...interface{})
	Plugins() []interface{}

	AddNode(config NodeConfig)
	Node(name string) Node
	BestNode() Node
	BestRestClient() RestClient
	RemoveNode(name string)

	Player(guildID discord.Snowflake) Player
	ExistingPlayer(guildID discord.Snowflake) Player
	Players() map[discord.Snowflake]Player

	UserID() discord.Snowflake
	SetUserID(userID discord.Snowflake)

	Close(ctx context.Context)

	VoiceServerUpdate(voiceServerUpdate VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate VoiceStateUpdate)
}

var _ Lavalink = (*lavalinkImpl)(nil)

type lavalinkImpl struct {
	config  Config
	nodes   map[string]Node
	players map[discord.Snowflake]Player
}

func (l *lavalinkImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *lavalinkImpl) AddPlugins(plugins ...interface{}) {
	l.config.Plugins = append(l.config.Plugins, plugins...)
}

func (l *lavalinkImpl) RemovePlugins(plugins ...interface{}) {
	for _, pl := range plugins {
		for i, p := range l.config.Plugins {
			if p == pl {
				l.config.Plugins = append(l.config.Plugins[:i], l.config.Plugins[i+1:]...)
				break
			}
		}
	}
}

func (l *lavalinkImpl) Plugins() []interface{} {
	return l.config.Plugins
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

func (l *lavalinkImpl) Player(guildID discord.Snowflake) Player {
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.BestNode(), guildID)
	l.players[guildID] = player
	return player
}

func (l *lavalinkImpl) ExistingPlayer(guildID discord.Snowflake) Player {
	return l.players[guildID]
}

func (l *lavalinkImpl) Players() map[discord.Snowflake]Player {
	return l.players
}

func (l *lavalinkImpl) UserID() discord.Snowflake {
	return l.config.UserID
}

func (l *lavalinkImpl) SetUserID(userID discord.Snowflake) {
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
