package internal

import (
	"github.com/DisgoOrg/disgolink/api"
	"github.com/DisgoOrg/disgolink/api/player"
)

type LinkImpl struct {
	player    player.Player
	node      api.Node
	lavalink  api.Lavalink
	guildID   api.Snowflake
	channelID *api.Snowflake
}

func (l *LinkImpl) Player() player.Player {
	return l.player
}
func (l *LinkImpl) Lavalink() api.Lavalink {
	return l.lavalink
}
func (l *LinkImpl) GuildID() api.Snowflake {
	return l.guildID
}
func (l *LinkImpl) ChannelID() *api.Snowflake {
	return l.channelID
}
func (l *LinkImpl) Node() api.Node {
	return l.node
}
func (l *LinkImpl) ChangeNode(node api.Node) {
	l.node = node
}
