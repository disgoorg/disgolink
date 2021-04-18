package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/log"
)

func NewLavalinkImpl(logger log.Logger, userID string) api.Lavalink {
	lavalink := &LavalinkImpl{
		logger:     logger,
		userID:     userID,
		httpClient: &http.Client{},
		players:    map[string]api.Player{},
	}
	return lavalink
}

type LavalinkImpl struct {
	logger     log.Logger
	userID     string
	nodes      []api.Node
	players    map[string]api.Player
	httpClient *http.Client
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

	l.nodes = append(l.nodes, node)
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

func (l *LavalinkImpl) RemoveNode(name string) {
	for i, node := range l.nodes {
		if node.Name() == name {
			l.nodes = append(l.nodes[:i], l.nodes[i+1:]...)
			return
		}
	}
}
func (l LavalinkImpl) RestClient() api.RestClient {
	if len(l.nodes) == 0 {
		return nil
	}
	// TODO: return best one
	return l.nodes[0].RestClient()
}
func (l *LavalinkImpl) Player(guildID string) api.Player {
	if player, ok := l.players[guildID]; ok {
		return player
	}
	player := NewPlayer(l.nodes[0], guildID)
	l.players[guildID] = player
	return player
}
func (l *LavalinkImpl) ExistingPlayer(guildID string) api.Player {
	return l.players[guildID]
}
func (l *LavalinkImpl) Players() map[string]api.Player {
	return nil //l.players
}
func (l *LavalinkImpl) UserID() string {
	return l.userID
}
func (l *LavalinkImpl) SetUserID(userID string) {
	l.userID = userID
}
func (l *LavalinkImpl) ClientName() string {
	return "disgolink"
}
func (l *LavalinkImpl) Shutdown() {

}

func (l *LavalinkImpl) VoiceServerUpdate(voiceServerUpdate *api.VoiceServerUpdate) {
	fmt.Printf("voiceServerUpdate: %+v", voiceServerUpdate)
	player := l.players[voiceServerUpdate.GuildID]
	if player == nil {
		return
	}
	player.Node().Send(api.EventCommand{
		GenericOpCommand: &api.GenericOpCommand{
			Op:      api.OpVoiceUpdate,
			GuildID: voiceServerUpdate.GuildID,
		},
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
