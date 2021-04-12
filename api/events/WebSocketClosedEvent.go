package events

import "github.com/DisgoOrg/disgolink/api"

type WebSocketClosedEvent struct {
	Op       api.OpType `json:"op"`
	Type     string `json:"type"`
	GuildId  string `json:"guildId"`
	Code     int    `json:"code"`
	Reason   string `json:"reason"`
	ByRemote bool   `json:"byRemote"`
}