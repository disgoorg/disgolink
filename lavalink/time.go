package lavalink

import (
	"strconv"
	"time"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatInt(t.UnixMilli(), 10))), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	timestampStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}
	*t = Time{
		Time: time.UnixMilli(timestamp),
	}
	return nil
}
