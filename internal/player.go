package internal

import "github.com/DisgoOrg/disgolink/api"

type PlayerOptions struct {
	Guild    string
	Channel  string
	SelfMute bool
	SelfDeaf bool
	Volume   int
}

type Player struct {
	api.PlayerState
	PlayerOptions
	state   *api.VoiceState
	Manager *Manager
	Node    *Node
	Filters *FilterManager
}

func (p *Player) Connect(channel string, selfMute, selfDeaf bool) {
	type Data struct {
		SelfDeaf  bool   `json:"self_deaf"`
		GuildId   string `json:"guild_id"`
		ChannelId string `json:"channel_id"`
		SelfMute  bool   `json:"self_mute"`
	}
	type ConnectMessage struct {
		Op int  `json:"op"`
		D  Data `json:"d"`
	}

	p.Manager.Options.Send(p.Guild, ConnectMessage{
		Op: 4,
		D: Data{
			SelfDeaf:  selfDeaf,
			GuildId:   p.Guild,
			ChannelId: channel,
			SelfMute:  selfMute,
		},
	})
}

func (p *Player) Disconnect() {
	type Data struct {
		SelfDeaf  bool   `json:"self_deaf"`
		GuildId   string `json:"guild_id"`
		ChannelId string `json:"channel_id"`
		SelfMute  bool   `json:"self_mute"`
	}
	type ConnectMessage struct {
		Op int  `json:"op"`
		D  Data `json:"d"`
	}

	p.Manager.Options.Send(p.Guild, ConnectMessage{
		Op: 4,
		D: Data{
			SelfDeaf:  p.SelfDeaf,
			GuildId:   p.Guild,
			ChannelId: "",
			SelfMute:  p.SelfMute,
		},
	})
}

func (p *Player) Destroy() error {
	p.Disconnect()

	type DestroyMessage struct {
		OP      api.OpType `json:"op"`
		GuildId string     `json:"guildId"`
	}

	defer delete(p.Manager.Players, p.Guild)
	return p.Node.Send(DestroyMessage{
		OP:      api.DestroyOp,
		GuildId: p.Guild,
	})
}

func (p *Player) stop() {
	type StopMessage struct {
		OP      api.OpType `json:"op"`
		GuildId string     `json:"guildId"`
	}
	_ = p.Node.Send(StopMessage{
		OP:      api.StopOp,
		GuildId: p.Guild,
	})
}

func (p *Player) Pause(state bool) error {
	type PauseMessage struct {
		OP      api.OpType `json:"op"`
		GuildId string     `json:"guildId"`
		State   bool       `json:"pause"`
	}

	return p.Node.Send(PauseMessage{
		OP:      api.PauseOP,
		GuildId: p.Guild,
		State:   state,
	})
}

func (p *Player) Seek(position int) error {
	type SeekMessage struct {
		Op       api.OpType `json:"op"`
		GuildId  string     `json:"guildId"`
		Position int        `json:"position"`
	}

	return p.Node.Send(SeekMessage{
		Op:       api.SeekOP,
		GuildId:  p.Guild,
		Position: position,
	})
}
