package api

import "github.com/DisgoOrg/disgolink/api/player"

type Link interface {
	Player() player.Player
	Lavalink() Lavalink
	GuildID() Snowflake
	ChannelID() *Snowflake
	Node() Node
	ChangeNode(node Node)
}
