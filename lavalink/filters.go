package lavalink

import "github.com/disgoorg/json"

var DefaultVolume Volume = 1.0

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

func (f *Filters) MarshalJSON() ([]byte, error) {
	filters := make(map[string]any)

	if f.Volume != nil {
		filters["volume"] = f.Volume
	}
	if f.Equalizer != nil {
		filters["equalizer"] = f.Equalizer
	}
	if f.Timescale != nil {
		filters["timescale"] = f.Timescale
	}
	if f.Tremolo != nil {
		filters["tremolo"] = f.Tremolo
	}
	if f.Vibrato != nil {
		filters["vibrato"] = f.Vibrato
	}
	if f.Rotation != nil {
		filters["rotation"] = f.Rotation
	}
	if f.Karaoke != nil {
		filters["karaoke"] = f.Karaoke
	}
	if f.Distortion != nil {
		filters["distortion"] = f.Distortion
	}
	if f.ChannelMix != nil {
		filters["channelMix"] = f.ChannelMix
	}
	if f.LowPass != nil {
		filters["lowPass"] = f.LowPass
	}
	for k, v := range f.PluginFilters {
		filters[k] = v
	}

	return json.Marshal(filters)
}

func UnmarshalFilters(data json.RawMessage) (*Filters, error) {
	var filters map[string]json.RawMessage
	err := json.Unmarshal(data, &filters)
	if err != nil {
		return nil, err
	}

	f := new(Filters)

	for k, v := range filters {
		switch k {
		case "volume":
			f.Volume = new(Volume)
			err = json.Unmarshal(v, f.Volume)
		case "equalizer":
			f.Equalizer = new(Equalizer)
			err = json.Unmarshal(v, f.Equalizer)
		case "timescale":
			f.Timescale = new(Timescale)
			err = json.Unmarshal(v, f.Timescale)
		case "tremolo":
			f.Tremolo = new(Tremolo)
			err = json.Unmarshal(v, f.Tremolo)
		case "vibrato":
			f.Vibrato = new(Vibrato)
			err = json.Unmarshal(v, f.Vibrato)
		case "rotation":
			f.Rotation = new(Rotation)
			err = json.Unmarshal(v, f.Rotation)
		case "karaoke":
			f.Karaoke = new(Karaoke)
			err = json.Unmarshal(v, f.Karaoke)
		case "distortion":
			f.Distortion = new(Distortion)
			err = json.Unmarshal(v, f.Distortion)
		case "channelMix":
			f.ChannelMix = new(ChannelMix)
			err = json.Unmarshal(v, f.ChannelMix)
		case "lowPass":
			f.LowPass = new(LowPass)
			err = json.Unmarshal(v, f.LowPass)
		default:
			if f.PluginFilters == nil {
				f.PluginFilters = make(map[string]any)
			}
			f.PluginFilters[k] = v
		}
	}
	return f, nil
}

func (f *Filters) GetVolume() *Volume {
	if f.Volume == nil {
		f.Volume = &DefaultVolume
	}
	return f.Volume
}

func (f *Filters) GetEqualizer() *Equalizer {
	if f.Equalizer == nil {
		f.Equalizer = new(Equalizer)
	}
	return f.Equalizer
}

func (f *Filters) GetTimescale() *Timescale {
	if f.Timescale == nil {
		f.Timescale = new(Timescale)
	}
	return f.Timescale
}

func (f *Filters) GetTremolo() *Tremolo {
	if f.Tremolo == nil {
		f.Tremolo = new(Tremolo)
	}
	return f.Tremolo
}

func (f *Filters) GetVibrato() *Vibrato {
	if f.Vibrato == nil {
		f.Vibrato = new(Vibrato)
	}
	return f.Vibrato
}

func (f *Filters) GetRotation() *Rotation {
	if f.Rotation == nil {
		f.Rotation = new(Rotation)
	}
	return f.Rotation
}

func (f *Filters) GetKaraoke() *Karaoke {
	if f.Karaoke == nil {
		f.Karaoke = new(Karaoke)
	}
	return f.Karaoke
}

func (f *Filters) GetDistortion() *Distortion {
	if f.Distortion == nil {
		f.Distortion = new(Distortion)
	}
	return f.Distortion
}

func (f *Filters) Clear() *Filters {
	f.Volume = nil
	f.Equalizer = nil
	f.Timescale = nil
	f.Tremolo = nil
	f.Vibrato = nil
	f.Rotation = nil
	f.Karaoke = nil
	f.Distortion = nil
	return f
}

type LowPass struct {
	Smoothing float64 `json:"smoothing"`
}

type ChannelMix struct {
	LeftToLeft   float64 `json:"leftToLeft,omitempty"`
	LeftToRight  float64 `json:"leftToRight,omitempty"`
	RightToLeft  float64 `json:"rightToLeft,omitempty"`
	RightToRight float64 `json:"rightToRight,omitempty"`
}

type Distortion struct {
	SinOffset int `json:"sinOffset"`
	SinScale  int `json:"sinScale"`
	CosOffset int `json:"cosOffset"`
	CosScale  int `json:"cosScale"`
	TanOffset int `json:"tanOffset"`
	TanScale  int `json:"tanScale"`
	Offset    int `json:"offset"`
	Scale     int `json:"scale"`
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
	Speed float32 `json:"speed"`
	Pitch float32 `json:"pitch"`
	Rate  float32 `json:"rate"`
}

type Tremolo struct {
	Frequency float32 `json:"frequency"`
	Depth     float32 `json:"depth"`
}

type Vibrato struct {
	Frequency float32 `json:"frequency"`
	Depth     float32 `json:"depth"`
}

type Volume float32

type Equalizer [15]float32

type EqBand struct {
	Band int     `json:"band"`
	Gain float32 `json:"gain"`
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
