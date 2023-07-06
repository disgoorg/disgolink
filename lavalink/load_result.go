package lavalink

import (
	"fmt"

	"github.com/disgoorg/json"
)

type LoadType string

const (
	LoadTypeTrack    LoadType = "track"
	LoadTypePlaylist LoadType = "playlist"
	LoadTypeSearch   LoadType = "search"
	LoadTypeEmpty    LoadType = "empty"
	LoadTypeError    LoadType = "error"
)

type LoadResultData interface {
	loadResultData()
}

type LoadResult struct {
	LoadType LoadType       `json:"loadType"`
	Data     LoadResultData `json:"data"`
}

func (r *LoadResult) UnmarshalJSON(data []byte) error {
	var raw struct {
		LoadType LoadType        `json:"loadType"`
		Data     json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.LoadType = raw.LoadType
	switch raw.LoadType {
	case LoadTypeTrack:
		var track Track
		if err := json.Unmarshal(raw.Data, &track); err != nil {
			return err
		}
		r.Data = track
	case LoadTypePlaylist:
		var playlist Playlist
		if err := json.Unmarshal(raw.Data, &playlist); err != nil {
			return err
		}
		r.Data = playlist
	case LoadTypeSearch:
		var search Search
		if err := json.Unmarshal(raw.Data, &search); err != nil {
			return err
		}
		r.Data = search
	case LoadTypeEmpty:
		r.Data = Empty{}
	case LoadTypeError:
		var exception Exception
		if err := json.Unmarshal(raw.Data, &exception); err != nil {
			return err
		}
		r.Data = exception
	default:
		return fmt.Errorf("unknown load type %q", raw.LoadType)
	}
	return nil
}

var _ error = (*Exception)(nil)

type Search []Track

func (Search) loadResultData() {}

type Empty struct{}

func (Empty) loadResultData() {}

type Exception struct {
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Cause    *string  `json:"cause,omitempty"`
}

func (Exception) loadResultData() {}

func (e Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.Severity, e.Message)
}

type Severity string

const (
	SeverityCommon     Severity = "common"
	SeveritySuspicious Severity = "suspicious"
	SeverityFault      Severity = "fault"
)
