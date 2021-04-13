package api

import (
	"github.com/DisgoOrg/log"
)

type Lavalink interface {
	Logger() log.Logger
	AddNode(options *NodeOptions)
	RemoveNode(name string)
	Player(guildID string) Player
	ExistingPlayer(guildID string) Player
	Players() map[string]Player
	RestClient() RestClient
	UserID() string
	ClientName() string
	Shutdown()
	VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate)
}
