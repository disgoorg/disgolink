package lavalink

import (
	"strconv"
	"time"
)

func Now() Timestamp {
	return Timestamp{
		Time: time.Now(),
	}
}

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.UnixMilli(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp{Time: time.UnixMilli(timestamp)}
	return nil
}
