package internal

import (
	"time"

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
	logger  log.Logger
	userID  api.Snowflake
	nodes   []api.Node
	players map[api.Snowflake]api.Player
}

func (l *LavalinkImpl) Logger() log.Logger {
	return l.logger
}

func (l *LavalinkImpl) AddNode(options api.NodeOptions) {
	node := &NodeImpl{
		NodeOptions: options,
		lavalink:    l,
	}
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
