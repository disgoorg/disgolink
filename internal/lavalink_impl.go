package internal

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"net/http"
	"time"

	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/log"
)

var _ api.Lavalink = (*LavalinkImpl)(nil)

func NewLavalinkImpl(logger log.Logger, userID dapi.Snowflake) api.Lavalink {
	if logger == nil {
		logger = log.Default()
	}
	lavalink := &LavalinkImpl{
		logger:     logger,
		userID:     userID,
		httpClient: &http.Client{},
		nodes:      map[string]api.Node{},
		players:    map[dapi.Snowflake]api.Player{},
	}
	return lavalink
}

type LavalinkImpl struct {
	logger     log.Logger
	userID     dapi.Snowflake
	httpClient *http.Client
	nodes      map[string]api.Node
	players    map[dapi.Snowflake]api.Player
}

func (l *LavalinkImpl) Logger() log.Logger {
	return l.logger
}

func (l *LavalinkImpl) AddNode(options *api.NodeOptions) {
	node := &NodeImpl{
		options:  options,
		lavalink: l,
	}
	node.restClient = newRestClientImpl(node, l.httpClient)

	l.nodes[options.Name] = node
	go func() {
		delay := 500
		for {
			err := node.Open()
			if err == nil {
				break
			}
			delay += int(float64(delay) * 1.2)
			l.Logger().Errorf("error while connecting to node: %s, waiting %ds, error: %s", node.Name(), delay/1000, err)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}()
}

func (l *LavalinkImpl) Node(name string) api.Node {
	return l.nodes[name]
}

func (l *LavalinkImpl) BestNode() api.Node {
	var bestNode api.Node
	for _, node := range l.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}
	return bestNode
}

func (l *LavalinkImpl) RemoveNode(name string) {
	delete(l.nodes, name)
}

func (l LavalinkImpl) RestClient() api.RestClient {
	if len(l.nodes) == 0 {
		return nil
	}
	return l.BestNode().RestClient()
}

func (l *LavalinkImpl) Player(guildID dapi.Snowflake) api.Player {
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.BestNode(), guildID)
	l.players[guildID] = player
	return player
}

func (l *LavalinkImpl) ExistingPlayer(guildID dapi.Snowflake) api.Player {
	return l.players[guildID]
}

func (l *LavalinkImpl) Players() map[dapi.Snowflake]api.Player {
	return l.players
}

func (l *LavalinkImpl) UserID() dapi.Snowflake {
	return l.userID
}

func (l *LavalinkImpl) SetUserID(userID dapi.Snowflake) {
	l.userID = userID
}

func (l *LavalinkImpl) ClientName() string {
	return "disgolink"
}

func (l *LavalinkImpl) Close() {
	for _, node := range l.nodes {
		node.Close()
	}
}

func (l *LavalinkImpl) VoiceServerUpdate(voiceServerUpdate *api.VoiceServerUpdate) {
	player := l.players[voiceServerUpdate.GuildID]
	if player == nil {
		return
	}
	player.Node().Send(api.EventCommand{
		GenericOp: api.GenericOp{
			Op: api.OpVoiceUpdate,
		},
		GuildID:   voiceServerUpdate.GuildID,
		SessionID: *player.LastSessionID(),
		Event:     voiceServerUpdate,
	})
}

func (l *LavalinkImpl) VoiceStateUpdate(voiceStateUpdate *api.VoiceStateUpdate) {
	player := l.players[voiceStateUpdate.GuildID]
	if player == nil {
		return
	}
	player.SetChannelID(voiceStateUpdate.ChannelID)
	player.SetLastSessionID(voiceStateUpdate.SessionID)
}
