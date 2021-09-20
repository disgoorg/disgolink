package disgolink

import (
	"github.com/DisgoOrg/disgo/discord"
	"net/http"
	"time"

	"github.com/DisgoOrg/log"
)

var _ Lavalink = (*defaultLavalink)(nil)

func newDefaultLavalink(logger log.Logger, httpClient *http.Client, userID discord.Snowflake) Lavalink {
	if logger == nil {
		logger = log.Default()
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	lavalink := &defaultLavalink{
		logger:     logger,
		userID:     userID,
		httpClient: httpClient,
		nodes:      map[string]Node{},
		players:    map[discord.Snowflake]Player{},
	}
	return lavalink
}

type defaultLavalink struct {
	logger     log.Logger
	userID     discord.Snowflake
	httpClient *http.Client
	nodes      map[string]Node
	players    map[discord.Snowflake]Player
}

func (l *defaultLavalink) Logger() log.Logger {
	return l.logger
}

func (l *defaultLavalink) AddNode(options *NodeOptions) {
	node := &defaultNode{
		options:  options,
		lavalink: l,
	}
	node.restClient = newDefaultRestClient(node, l.httpClient)

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

func (l *defaultLavalink) Node(name string) Node {
	return l.nodes[name]
}

func (l *defaultLavalink) BestNode() Node {
	var bestNode Node
	for _, node := range l.nodes {
		if bestNode == nil || node.Stats().Better(bestNode.Stats()) {
			bestNode = node
		}
	}
	return bestNode
}

func (l *defaultLavalink) RemoveNode(name string) {
	delete(l.nodes, name)
}

func (l defaultLavalink) RestClient() RestClient {
	if len(l.nodes) == 0 {
		return nil
	}
	return l.BestNode().RestClient()
}

func (l *defaultLavalink) Player(guildID discord.Snowflake) Player {
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.BestNode(), guildID)
	l.players[guildID] = player
	return player
}

func (l *defaultLavalink) ExistingPlayer(guildID discord.Snowflake) Player {
	return l.players[guildID]
}

func (l *defaultLavalink) Players() map[discord.Snowflake]Player {
	return l.players
}

func (l *defaultLavalink) UserID() discord.Snowflake {
	return l.userID
}

func (l *defaultLavalink) SetUserID(userID discord.Snowflake) {
	l.userID = userID
}

func (l *defaultLavalink) ClientName() string {
	return "disgolink"
}

func (l *defaultLavalink) Close() {
	for _, node := range l.nodes {
		node.Close()
	}
}

func (l *defaultLavalink) VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate) {
	player := l.players[voiceServerUpdate.GuildID]
	if player == nil {
		return
	}
	player.Node().Send(EventCommand{
		GenericOp: GenericOp{
			Op: OpVoiceUpdate,
		},
		GuildID:   voiceServerUpdate.GuildID,
		SessionID: *player.LastSessionID(),
		Event:     voiceServerUpdate,
	})
}

func (l *defaultLavalink) VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate) {
	player := l.players[voiceStateUpdate.GuildID]
	if player == nil {
		return
	}
	player.SetChannelID(voiceStateUpdate.ChannelID)
	player.SetLastSessionID(voiceStateUpdate.SessionID)
}
