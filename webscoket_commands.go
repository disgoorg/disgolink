package disgolink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/filters"
)

type GenericOp struct {
	Op Op `json:"op"`
}

type EventCommand struct {
	GenericOp
	SessionID string         `json:"sessionId"`
	Event     interface{}    `json:"event"`
	GuildID   discord.Snowflake `json:"guildId"`
}

type PlayerCommand struct {
	GenericOp
	GuildID discord.Snowflake `json:"guildId"`
}

func NewPlayerCommand(op Op, p Player) PlayerCommand {
	return PlayerCommand{
		GenericOp: GenericOp{
			Op: op,
		},
		GuildID: p.GuildID(),
	}
}

type PlayPlayerCommand struct {
	PlayerCommand
	Track     string `json:"track"`
	StartTime *int   `json:"startTime,omitempty"`
	EndTime   *int   `json:"endTime,omitempty"`
	Volume    *int   `json:"volume,omitempty"`
	NoReplace *bool  `json:"noReplace,omitempty"`
	Pause     *bool  `json:"pause,omitempty"`
}

type DestroyPlayerCommand struct {
	PlayerCommand
}

type StopPlayerCommand struct {
	PlayerCommand
}

type PausePlayerCommand struct {
	PlayerCommand
	Pause bool `json:"pause"`
}

type SeekPlayerCommand struct {
	PlayerCommand
	Position int `json:"position"`
}

type VolumePlayerCommand struct {
	PlayerCommand
	Volume int `json:"volume"`
}

type FilterPlayerCommand struct {
	PlayerCommand
	*filters.Filters
}
