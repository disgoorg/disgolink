package lavalink

import (
	"fmt"

	"github.com/disgoorg/json/v2"
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

	var (
		resultData LoadResultData
		err        error
	)

	switch raw.LoadType {
	case LoadTypeTrack:
		var v Track
		err = json.Unmarshal(raw.Data, &v)
		resultData = v
	case LoadTypePlaylist:
		var v Playlist
		err = json.Unmarshal(raw.Data, &v)
		resultData = v
	case LoadTypeSearch:
		var v Search
		err = json.Unmarshal(raw.Data, &v)
		resultData = v
	case LoadTypeEmpty:
		resultData = Empty{}
	case LoadTypeError:
		var v Exception
		err = json.Unmarshal(raw.Data, &v)
		resultData = v
	default:
		return fmt.Errorf("unknown load type %q", raw.LoadType)
	}
	if err != nil {
		return fmt.Errorf("error while unmarshalling load result data: %w", err)
	}

	r.LoadType = raw.LoadType
	r.Data = resultData

	return nil
}

var _ LoadResultData = (*Search)(nil)

type Search []Track

func (Search) loadResultData() {}

var _ LoadResultData = (*Empty)(nil)

type Empty struct{}

func (Empty) loadResultData() {}

var (
	_ error          = (*Exception)(nil)
	_ LoadResultData = (*Exception)(nil)
)

type Exception struct {
	Message         string   `json:"message"`
	Severity        Severity `json:"severity"`
	Cause           string   `json:"cause"`
	CauseStackTrace string   `json:"causeStackTrace"`
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
