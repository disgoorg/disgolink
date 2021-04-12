package api

import (
	"github.com/DisgoOrg/log"
)

type Lavalink interface {
	Logger() log.Logger
	AddNode(options NodeOptions)
	RemoveNode(name string)
	Link(guildID Snowflake) Link
	ExistingLink(guildID Snowflake) Link
	Links() map[Snowflake]Link
	UserID() Snowflake
	ClientName() string
	Shutdown()
	VoiceServerUpdate(voiceServerUpdate *VoiceServerUpdate)
	VoiceStateUpdate(voiceStateUpdate *VoiceStateUpdate)
}
