package disgolink

import "github.com/DisgoOrg/disgolink/api"

func NewDisgolink() Disgolink {
	return nil
}

type Disgolink interface {
	api.Lavalink
}
