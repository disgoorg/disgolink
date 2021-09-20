package disgolink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

type Lavalink interface {
	Logger() log.Logger
	AddNode(options *NodeOptions)
	Node(name string) Node
	BestNode() Node
	RemoveNode(name string)
	Player(guildID discord.Snowflake) Player
	ExistingPlayer(guildID discord.Snowflake) Player
	Players() map[discord.Snowflake]Player
	RestClient() RestClient
	UserID() discord.Snowflake
	ClientName() string
	Close()
	VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate)
}
