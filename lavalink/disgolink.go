package lavalink

import (
	"github.com/DisgoOrg/disgo/core"
)

func NewDisgolink(disgo *core.Bot) Disgolink {
	dgolink := &disgolinkImpl{
		Lavalink: NewLavalink(
			WithLogger(disgo.Logger),
			WithHTTPClient(disgo.RestServices.HTTPClient()),
			WithUserID(disgo.ApplicationID),
		),
	}

	disgo.EventManager.AddEventListeners(dgolink)
	return dgolink
}

type Disgolink interface {
	Lavalink
	core.EventListener
}

var (
	_ Disgolink = (*disgolinkImpl)(nil)
	_ Lavalink = (*disgolinkImpl)(nil)
	_ core.EventListener = (*disgolinkImpl)(nil)
)

type disgolinkImpl struct {
	Lavalink
}

func (l *disgolinkImpl) OnEvent(event core.Event) {
	switch _ := event.(type) {

	}
}