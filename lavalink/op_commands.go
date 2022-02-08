package lavalink

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/snowflake"
)

type PlayCommand struct {
	GuildID   snowflake.Snowflake `json:"guildId"`
	Track     string              `json:"track"`
	StartTime int64               `json:"startTime,omitempty"`
	EndTime   int64               `json:"endTime,omitempty"`
	NoReplace bool                `json:"noReplace,omitempty"`
	Pause     bool                `json:"pause,omitempty"`
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
func (PlayCommand) OpCommand() {}

type StopCommand struct {
	GuildID snowflake.Snowflake `json:"guildId"`
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
func (StopCommand) OpCommand() {}

type DestroyCommand struct {
	GuildID snowflake.Snowflake `json:"guildId"`
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
func (DestroyCommand) OpCommand() {}

type PauseCommand struct {
	GuildID snowflake.Snowflake `json:"guildId"`
	Pause   bool                `json:"pause"`
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
func (PauseCommand) OpCommand() {}

type SeekCommand struct {
	GuildID  snowflake.Snowflake `json:"guildId"`
	Position int64               `json:"position"`
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
func (SeekCommand) OpCommand() {}

type VolumeCommand struct {
	GuildID snowflake.Snowflake `json:"guildId"`
	Volume  int                 `json:"volume"`
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
func (VolumeCommand) OpCommand() {}

type VoiceUpdateCommand struct {
	GuildID   snowflake.Snowflake `json:"guildId"`
	SessionID string              `json:"sessionId"`
	Event     VoiceServerUpdate   `json:"event"`
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
func (VoiceUpdateCommand) OpCommand() {}

type ConfigureResumingCommand struct {
	Key     string        `json:"key"`
	Timeout time.Duration `json:"timeout"`
}

func (c ConfigureResumingCommand) MarshalJSON() ([]byte, error) {
	type cmd ConfigureResumingCommand
	return json.Marshal(struct {
		Op      OpType `json:"op"`
		Timeout int64  `json:"timeout"`
		cmd
	}{
		Op:      c.Op(),
		Timeout: c.Timeout.Milliseconds(),
		cmd:     cmd(c),
	})
}
func (ConfigureResumingCommand) Op() OpType { return OpTypeConfigureResuming }
func (ConfigureResumingCommand) OpCommand() {}

type FiltersCommand struct {
	GuildID snowflake.Snowflake `json:"guildId"`
	Filters
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
func (FiltersCommand) OpCommand() {}
