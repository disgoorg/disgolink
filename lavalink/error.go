package lavalink

import (
	"fmt"
	"time"
)

type Error struct {
	Timestamp   time.Time `json:"timestamp"`
	Status      int       `json:"status"`
	StatusError string    `json:"error"`
	Trace       string    `json:"trace"`
	Message     string    `json:"message"`
	Path        string    `json:"path"`
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %d - %s", e.Path, e.Status, e.StatusError)
}
