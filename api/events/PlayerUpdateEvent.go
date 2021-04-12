package events

import "github.com/DisgoOrg/disgolink/api"

type PlayerUpdateEvent struct {
	Op      api.OpType      `json:"op"`
	GuildId string          `json:"guildId"`
	State   api.PlayerState `json:"state"`
}
