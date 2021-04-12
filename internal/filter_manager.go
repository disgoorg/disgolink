package internal

import "github.com/DisgoOrg/disgolink/api/filters"

type FilterManager struct {
	Volume     float32             `json:"volume"`
	Equalizer  *filters.Equalizer  `json:"equalizer"`
	Timescale  *filters.Timescale  `json:"timescale"`
	Tremolo    *filters.Tremolo    `json:"tremolo"`
	Vibrato    *filters.Vibrato    `json:"vibrato"`
	Rotation   *filters.Rotation   `json:"rotation"`
	Karaoke    *filters.Karaoke    `json:"karaoke"`
	Distortion *filters.Distortion `json:"distortion"`
}

func (f FilterManager) CreateFilterManager() *FilterManager {
	return &FilterManager{
		Volume: 1.0,
		Equalizer: &filters.Equalizer{
			filters.EqBand{
				Band: 0,
				Gain: 0.2,
			},
		},
		Timescale: &filters.Timescale{
			Speed: 1.0,
			Pitch: 1.0,
			Rate:  1.0,
		},
		Tremolo: &filters.Tremolo{
			Frequency: 2.0,
			Depth:     0.5,
		},
		Vibrato: &filters.Vibrato{
			Frequency: 2.0,
			Depth:     0.5,
		},
		Rotation: &filters.Rotation{RotationHz: 0},
		Karaoke: &filters.Karaoke{
			Level:       1.0,
			MonoLevel:   1.0,
			FilterBand:  220.0,
			FilterWidth: 100.0,
		},
		Distortion: &filters.Distortion{
			SinOffset: 0,
			SinScale:  1,
			CosOffset: 0,
			CosScale:  1,
			TanOffset: 0,
			TanScale:  1,
			Offset:    0,
			Scale:     1,
		},
	}
}
