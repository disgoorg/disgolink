package internal

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/log"
)

func NewLavalinkImpl(logger log.Logger, userID api.Snowflake) *LavalinkImpl {
	return &LavalinkImpl{
		logger: logger,
		userID: userID,
	}
}

type LavalinkImpl struct {
	logger log.Logger
	userID api.Snowflake
	nodes  []api.Node
	players  map[api.Snowflake]api.Player
}

func (l *LavalinkImpl) Logger() log.Logger {
	return l.logger
}

func (l *LavalinkImpl) AddNode(options api.NodeOptions) {
	l.nodes = append(l.nodes, &NodeImpl{
		NodeOptions: options,
		lavalink:    l,
	})
}

func (l *LavalinkImpl) RemoveNode(name string) {
	for i, node := range l.nodes {
		if node.Name() == name {
			l.nodes = append(l.nodes[:i], l.nodes[i+1:]...)
			return
		}
	}
}
func (l *LavalinkImpl) Player(guildID api.Snowflake) api.Player {
	if link, ok := l.players[guildID]; ok {
		return link
	}
	// create new link
	return nil
}
func (l *LavalinkImpl) ExistingPlayer(guildID api.Snowflake) api.Player {
	return l.players[guildID]
}
func (l *LavalinkImpl) Players() map[api.Snowflake]api.Player {
	return l.players
}
func (l *LavalinkImpl) UserID() api.Snowflake {
	return l.userID
}
func (l *LavalinkImpl) SetUserID(userID api.Snowflake) {
	l.userID = userID
}
func (l *LavalinkImpl) ClientName() string {
	return "disgolink"
}
func (l *LavalinkImpl) Shutdown() {

}

func (l *LavalinkImpl) VoiceServerUpdate(voiceServerUpdate *api.VoiceServerUpdate) {

}

func (l *LavalinkImpl) VoiceStateUpdate(voiceStateUpdate *api.VoiceStateUpdate) {

}
