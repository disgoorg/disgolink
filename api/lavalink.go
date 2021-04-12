package api

import (
	"github.com/DisgoOrg/log"
)

type Lavalink interface {
	Logger() log.Logger
	AddNode(options NodeOptions)
	RemoveNode(name string)
	Player(guildID Snowflake) Player
	ExistingPlayer(guildID Snowflake) Player
	Players() map[Snowflake]Player
	UserID() Snowflake
	ClientName() string
	Shutdown()
	VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate)
}
