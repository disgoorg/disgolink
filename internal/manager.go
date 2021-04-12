package internal

import "github.com/DisgoOrg/disgolink/api"

type ManagerOptions struct {
	Nodes      []NodeOptions
	ClientID   *string
	ClientName *string
	Shards     int
	Send       func(guildID string, payload interface{})
}

type Manager struct {
	Options ManagerOptions
	Nodes   []*Node
	Players map[string]*Player
}

func CreateManager() {

}

func (m *Manager) OnVoiceServerUpdate(voiceServerUpdateEvent *api.VoiceServer) error {
	type VoiceServerMessage struct {
		Op        api.OpType      `json:"op"`
		GuildID   string          `json:"guildId"`
		SessionID string          `json:"sessionId"`
		Event     api.VoiceServer `json:"event"`
	}

	player, ok := m.Players[voiceServerUpdateEvent.GuildId]
	if ok && player.state != nil {
		return player.Node.Send(VoiceServerMessage{
			Op:        api.PlayerUpdateOp,
			GuildID:   voiceServerUpdateEvent.GuildId,
			SessionID: player.state.SessionId,
			Event:     *voiceServerUpdateEvent,
		})
	}
	return nil
}

func (m *Manager) OnVoiceStateUpdate(voiceStateUpdateEvent *api.VoiceState) {
	player, ok := m.Players[voiceStateUpdateEvent.GuildId]
	if ok {
		player.state = voiceStateUpdateEvent
	}
}
