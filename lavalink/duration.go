package lavalink

import (
	"strconv"
)

type Duration int64

const (
	Millisecond Duration = 1
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
)

func (d Duration) Milliseconds() int64 {
	return int64(d)
}

func (d Duration) MillisecondsPart() int64 {
	return int64(Duration(d.Milliseconds()) % 1000)
}

func (d Duration) Seconds() int64 {
	return int64(d / Second)
}

func (d Duration) SecondsPart() int64 {
	return int64(Duration(d.Seconds()) % 60)
}

func (d Duration) Minutes() int64 {
	return int64(d / Minute)
}

func (d Duration) MinutesPart() int64 {
	return int64(Duration(d.Minutes()) % 60)
}

func (d Duration) Hours() int64 {
	return int64(d / Hour)
}

func (d Duration) HoursPart() int64 {
	return int64(Duration(d.Hours()) % 24)
}

func (d Duration) Days() int64 {
	return int64(d / Day)
}

func (d Duration) String() string {
	var str string
	if days := d.Days(); days > 0 {
		str += strconv.FormatInt(days, 10) + "d"
	}
	if hours := d.HoursPart(); hours > 0 {
		str += strconv.FormatInt(hours, 10) + "h"
	}
	if minutes := d.MinutesPart(); minutes > 0 {
		str += strconv.FormatInt(minutes, 10) + "m"
	}
	if seconds := d.SecondsPart(); seconds > 0 {
		str += strconv.FormatInt(seconds, 10) + "s"
	}
	if milliseconds := d.MillisecondsPart(); milliseconds > 0 {
		str += strconv.FormatInt(milliseconds, 10) + "ms"
	}
	return str
}
