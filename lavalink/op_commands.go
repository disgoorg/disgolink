package lavalink

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgolink/filters"
)

type PlayCommand struct {
	GuildID   discord.Snowflake `json:"guildId"`
	Track     string            `json:"track"`
	StartTime int               `json:"startTime,omitempty"`
	EndTime   int               `json:"endTime,omitempty"`
	NoReplace bool              `json:"noReplace,omitempty"`
	Pause     bool              `json:"pause,omitempty"`
}

func (c PlayCommand) MarshalJSON() ([]byte, error) {
	type playCommand PlayCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		playCommand
	}{
		Op:          c.Op(),
		playCommand: playCommand(c),
	})
}
func (PlayCommand) Op() OpType { return OpTypePlay }
func (PlayCommand) opCommand() {}

type StopCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
}

func (c StopCommand) MarshalJSON() ([]byte, error) {
	type cmd StopCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (StopCommand) Op() OpType { return OpTypeStop }
func (StopCommand) opCommand() {}

type DestroyCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
}

func (c DestroyCommand) MarshalJSON() ([]byte, error) {
	type cmd DestroyCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (DestroyCommand) Op() OpType { return OpTypeDestroy }
func (DestroyCommand) opCommand() {}

type PauseCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
	Pause   bool              `json:"pause"`
}

func (c PauseCommand) MarshalJSON() ([]byte, error) {
	type cmd PauseCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (PauseCommand) Op() OpType { return OpTypePause }
func (PauseCommand) opCommand() {}

type SeekCommand struct {
	GuildID  discord.Snowflake `json:"guildId"`
	Position int               `json:"position"`
}

func (c SeekCommand) MarshalJSON() ([]byte, error) {
	type cmd SeekCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (SeekCommand) Op() OpType { return OpTypeSeek }
func (SeekCommand) opCommand() {}

type VolumeCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
	Volume  int               `json:"volume"`
}

func (c VolumeCommand) MarshalJSON() ([]byte, error) {
	type cmd VolumeCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (VolumeCommand) Op() OpType { return OpTypeVolume }
func (VolumeCommand) opCommand() {}

type VoiceUpdateCommand struct {
	GuildID   discord.Snowflake `json:"guildId"`
	SessionID string            `json:"sessionId"`
	Event     VoiceServerUpdate `json:"event"`
}

func (c VoiceUpdateCommand) MarshalJSON() ([]byte, error) {
	type cmd VoiceUpdateCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (VoiceUpdateCommand) Op() OpType { return OpTypeVoiceUpdate }
func (VoiceUpdateCommand) opCommand() {}

type ConfigureResumingCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
}

func (c ConfigureResumingCommand) MarshalJSON() ([]byte, error) {
	type cmd ConfigureResumingCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (ConfigureResumingCommand) Op() OpType { return OpTypeConfigureResuming }
func (ConfigureResumingCommand) opCommand() {}

type FiltersCommand struct {
	GuildID discord.Snowflake `json:"guildId"`
	filters.Filters
}

func (c FiltersCommand) MarshalJSON() ([]byte, error) {
	type cmd FiltersCommand
	return json.Marshal(struct {
		Op OpType `json:"op"`
		cmd
	}{
		Op:  c.Op(),
		cmd: cmd(c),
	})
}
func (FiltersCommand) Op() OpType { return OpTypeFilters }
func (FiltersCommand) opCommand() {}
