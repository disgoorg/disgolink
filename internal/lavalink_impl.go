package internal

import "github.com/DisgoOrg/disgolink/api"

type LavalinkImpl struct {
	userID api.Snowflake
	nodes  []api.Node
	links  map[api.Snowflake]api.Link
}

func (l *LavalinkImpl) AddNode(options api.NodeOptions) {
	l.nodes = append(l.nodes, &NodeImpl{
		NodeOptions: options,
		lavalink:    l,
	})
}

func (l *LavalinkImpl) RemoveNode(name string) {
	for i, node := range l.nodes {
		if node.Name() == name {
			l.nodes = append(l.nodes[:i], l.nodes[i+1:]...)
			return
		}
	}
}
func (l *LavalinkImpl) Link(guildID api.Snowflake) api.Link {
	if link, ok := l.links[guildID]; ok {
		return link
	}
	// create new link
	return nil
}
func (l *LavalinkImpl) ExistingLink(guildID api.Snowflake) api.Link {
	return l.links[guildID]
}
func (l *LavalinkImpl) Links() map[api.Snowflake]api.Link {
	return l.links
}
func (l *LavalinkImpl) UserID() api.Snowflake {
	return l.userID
}
func (l *LavalinkImpl) SetUserID(userID api.Snowflake) {
	l.userID = userID
}
func (l *LavalinkImpl) ClientName() string {
	return "disgolink"
}
func (l *LavalinkImpl) Shutdown() {

}

func (l *LavalinkImpl) OnVoiceServerUpdate(voiceServerUpdateEvent *api.VoiceServer) {

}

func (l *LavalinkImpl) OnVoiceStateUpdate(voiceStateUpdateEvent *api.VoiceState) {

}
