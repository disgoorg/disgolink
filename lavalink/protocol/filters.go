package protocol

import "encoding/json"

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
	PluginFilters map[string]any `json:"-"`
}

func (f *Filters) UnmarshalJSON(data []byte) error {
	type filters Filters
	var v filters
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*f = Filters(v)

	// Unmarshal plugin filters
	var pluginFilters map[string]any
	if err := json.Unmarshal(data, &pluginFilters); err != nil {
		return err
	}
	for k := range pluginFilters {
		if k == "volume" || k == "equalizer" || k == "timescale" || k == "tremolo" || k == "vibrato" || k == "rotation" || k == "karaoke" || k == "distortion" || k == "channelMix" || k == "lowPass" {
			delete(pluginFilters, k)
		}
	}
	f.PluginFilters = pluginFilters
	return nil
}

func (f Filters) MarshalJSON() ([]byte, error) {
	v := make(map[string]any)
	for k, val := range f.PluginFilters {
		v[k] = val
	}

	if f.Volume != nil {
		v["volume"] = f.Volume
	}
	if f.Equalizer != nil {
		v["equalizer"] = f.Equalizer
	}
	if f.Timescale != nil {
		v["timescale"] = f.Timescale
	}
	if f.Tremolo != nil {
		v["tremolo"] = f.Tremolo
	}
	if f.Vibrato != nil {
		v["vibrato"] = f.Vibrato
	}
	if f.Rotation != nil {
		v["rotation"] = f.Rotation
	}
	if f.Karaoke != nil {
		v["karaoke"] = f.Karaoke
	}
	if f.Distortion != nil {
		v["distortion"] = f.Distortion
	}
	if f.ChannelMix != nil {
		v["channelMix"] = f.ChannelMix
	}
	if f.LowPass != nil {
		v["lowPass"] = f.LowPass
	}

	return json.Marshal(v)
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
