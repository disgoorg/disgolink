package lavalink

import (
	"strconv"
	"time"
)

func Now() Time {
	return Time{
		Time: time.Now(),
	}
}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.UnixMilli(), 10)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{Time: time.UnixMilli(timestamp)}
	return nil
}
