package lavalink

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/disgoorg/json"
)

var (
	_ driver.Valuer = (*Track)(nil)
	_ sql.Scanner   = (*Track)(nil)
)

type Track struct {
	Encoded    string    `json:"encoded"`
	Info       TrackInfo `json:"info"`
	PluginInfo RawData   `json:"pluginInfo"`
	UserData   RawData   `json:"userData"`
}

// WithUserData returns a copy of the Track with the given userData.
func (t Track) WithUserData(userData any) (Track, error) {
	userDataRaw, err := json.Marshal(userData)
	if err != nil {
		return t, fmt.Errorf("failed to marshal userData: %w", err)
	}
	t.UserData = userDataRaw
	return t, nil
}

func (Track) loadResultData() {}

func (t Track) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Track) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &t)
}

type TrackInfo struct {
	Identifier string   `json:"identifier"`
	Author     string   `json:"author"`
	Length     Duration `json:"length"`
	IsStream   bool     `json:"isStream"`
	Title      string   `json:"title"`
	URI        *string  `json:"uri"`
	SourceName string   `json:"sourceName"`
	Position   Duration `json:"position"`
	ArtworkURL *string  `json:"artworkUrl"`
	ISRC       *string  `json:"isrc"`
}
