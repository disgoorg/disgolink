package lavalink

import (
	"encoding/json"
)

var DefaultFilters = []string{"volume", "equalizer", "timescale", "tremolo", "vibrato", "rotation", "karaoke", "distortion", "channelMix", "lowPass"}

type Filters struct {
	Volume        *Volume        `json:"volume,omitempty"`
	Equalizer     *Equalizer     `json:"equalizer,omitempty"`
	Timescale     *Timescale     `json:"timescale,omitempty"`
	Tremolo       *Tremolo       `json:"tremolo,omitempty"`
	Vibrato       *Vibrato       `json:"vibrato,omitempty"`
	Rotation      *Rotation      `json:"rotation,omitempty"`
	Karaoke       *Karaoke       `json:"karaoke,omitempty"`
	Distortion    *Distortion    `json:"distortion,omitempty"`
	ChannelMix    *ChannelMix    `json:"channelMix,omitempty"`
	LowPass       *LowPass       `json:"lowPass,omitempty"`
	PluginFilters map[string]any `json:"pluginFilters,omitempty"`
}

type LowPass struct {
	Smoothing float64 `json:"smoothing"`
}

type ChannelMix struct {
	LeftToLeft   float32 `json:"leftToLeft"`
	LeftToRight  float32 `json:"leftToRight"`
	RightToLeft  float32 `json:"rightToLeft"`
	RightToRight float32 `json:"rightToRight"`
}

type Distortion struct {
	SinOffset float32 `json:"sinOffset"`
	SinScale  float32 `json:"sinScale"`
	CosOffset float32 `json:"cosOffset"`
	CosScale  float32 `json:"cosScale"`
	TanOffset float32 `json:"tanOffset"`
	TanScale  float32 `json:"tanScale"`
	Offset    float32 `json:"offset"`
	Scale     float32 `json:"scale"`
}

type Vibrato struct {
	Frequency float32 `json:"frequency"`
	Depth     float32 `json:"depth"`
}

type Karaoke struct {
	Level       float32 `json:"level"`
	MonoLevel   float32 `json:"monoLevel"`
	FilterBand  float32 `json:"filterBand"`
	FilterWidth float32 `json:"filterWidth"`
}

type Rotation struct {
	RotationHz int `json:"rotationHz"`
}

type Timescale struct {
	Speed float64 `json:"speed"`
	Pitch float64 `json:"pitch"`
	Rate  float64 `json:"rate"`
}

type Tremolo struct {
	Frequency float32 `json:"frequency"`
	Depth     float32 `json:"depth"`
}

type Volume float32

type Equalizer [15]float32

type EqBand struct {
	Band int     `json:"band"`
	Gain float32 `json:"gain"`
}

func (e *Equalizer) UnmarshalJSON(data []byte) error {
	var bands [15]EqBand
	if err := json.Unmarshal(data, &bands); err != nil {
		return err
	}
	for _, band := range bands {
		e[band.Band] = band.Gain
	}
	return nil
}

// MarshalJSON marshals the map as object array
func (e Equalizer) MarshalJSON() ([]byte, error) {
	var bands [15]EqBand
	for band, gain := range e {
		bands[band] = EqBand{
			Band: band,
			Gain: gain,
		}
	}
	return json.Marshal(bands)
}
