package events

import "github.com/DisgoOrg/disgolink/api"

type Stats struct {
	Op    api.OpType `json:"op"`
	Stats *api.Stats
}
