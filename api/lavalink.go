package api

import (
	dapi "github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/log"
)

type Lavalink interface {
	Logger() log.Logger
	AddNode(options *NodeOptions)
	Node(name string) Node
	BestNode() Node
	RemoveNode(name string)
	Player(guildID dapi.Snowflake) Player
	ExistingPlayer(guildID dapi.Snowflake) Player
	Players() map[dapi.Snowflake]Player
	RestClient() RestClient
	UserID() dapi.Snowflake
	ClientName() string
	Close()
	VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate)
}
