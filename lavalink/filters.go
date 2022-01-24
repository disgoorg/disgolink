package lavalink

import "encoding/json"

var DefaultVolume Volume = 1.0

type Filters interface {
	Volume() *Volume
	Equalizer() *Equalizer
	Timescale() *Timescale
	Tremolo() *Tremolo
	Vibrato() *Vibrato
	Rotation() *Rotation
	Karaoke() *Karaoke
	Distortion() *Distortion

	SetVolume(v *Volume) Filters
	SetEqualizer(equalizer *Equalizer) Filters
	SetTimescale(timescale *Timescale) Filters
	SetTremolo(tremolo *Tremolo) Filters
	SetVibrato(vibrato *Vibrato) Filters
	SetRotation(rotation *Rotation) Filters
	SetKaraoke(karaoke *Karaoke) Filters
	SetDistortion(distortion *Distortion) Filters

	Clear() Filters
	Commit() error
}

func NewFilters(commitFunc func(filters Filters) error) Filters {
	return &defaultFilters{commitFunc: commitFunc}
}

var _ Filters = (*defaultFilters)(nil)

type defaultFilters struct {
	FilterVolume     *Volume     `json:"volume,omitempty"`
	FilterEqualizer  *Equalizer  `json:"equalizer,omitempty"`
	FilterTimescale  *Timescale  `json:"timescale,omitempty"`
	FilterTremolo    *Tremolo    `json:"tremolo,omitempty"`
	FilterVibrato    *Vibrato    `json:"vibrato,omitempty"`
	FilterRotation   *Rotation   `json:"rotation,omitempty"`
	FilterKaraoke    *Karaoke    `json:"karaoke,omitempty"`
	FilterDistortion *Distortion `json:"distortion,omitempty"`
	commitFunc       func(filters Filters) error
}

func (f *defaultFilters) Volume() *Volume {
	if f.FilterVolume == nil {
		f.FilterVolume = &DefaultVolume
	}
	return f.FilterVolume
}

func (f *defaultFilters) SetVolume(volume *Volume) Filters {
	f.FilterVolume = volume
	return f
}

func (f *defaultFilters) Equalizer() *Equalizer {
	if f.FilterEqualizer == nil {
		f.FilterEqualizer = new(Equalizer)
	}
	return f.FilterEqualizer
}

func (f *defaultFilters) SetEqualizer(equalizer *Equalizer) Filters {
	f.FilterEqualizer = equalizer
	return f
}

func (f *defaultFilters) Timescale() *Timescale {
	if f.FilterTimescale == nil {
		f.FilterTimescale = new(Timescale)
	}
	return f.FilterTimescale
}

func (f *defaultFilters) SetTimescale(timescale *Timescale) Filters {
	f.FilterTimescale = timescale
	return f
}

func (f *defaultFilters) Tremolo() *Tremolo {
	if f.FilterTremolo == nil {
		f.FilterTremolo = new(Tremolo)
	}
	return f.FilterTremolo
}

func (f *defaultFilters) SetTremolo(tremolo *Tremolo) Filters {
	f.FilterTremolo = tremolo
	return f
}

func (f *defaultFilters) Vibrato() *Vibrato {
	if f.FilterVibrato == nil {
		f.FilterVibrato = new(Vibrato)
	}
	return f.FilterVibrato
}

func (f *defaultFilters) SetVibrato(vibrato *Vibrato) Filters {
	f.FilterVibrato = vibrato
	return f
}

func (f *defaultFilters) Rotation() *Rotation {
	if f.FilterRotation == nil {
		f.FilterRotation = new(Rotation)
	}
	return f.FilterRotation
}

func (f *defaultFilters) SetRotation(rotation *Rotation) Filters {
	f.FilterRotation = rotation
	return f
}

func (f *defaultFilters) Karaoke() *Karaoke {
	if f.FilterKaraoke == nil {
		f.FilterKaraoke = new(Karaoke)
	}
	return f.FilterKaraoke
}

func (f *defaultFilters) SetKaraoke(karaoke *Karaoke) Filters {
	f.FilterKaraoke = karaoke
	return f
}

func (f *defaultFilters) Distortion() *Distortion {
	if f.FilterDistortion == nil {
		f.FilterDistortion = new(Distortion)
	}
	return f.FilterDistortion
}

func (f *defaultFilters) SetDistortion(distortion *Distortion) Filters {
	f.FilterDistortion = distortion
	return f
}

func (f *defaultFilters) Clear() Filters {
	f.FilterVolume = nil
	f.FilterEqualizer = nil
	f.FilterTimescale = nil
	f.FilterTremolo = nil
	f.FilterVibrato = nil
	f.FilterRotation = nil
	f.FilterKaraoke = nil
	f.FilterDistortion = nil
	return f
}

func (f *defaultFilters) Commit() error {
	return f.commitFunc(f)
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

type Equalizer map[int]float32

type EqBand struct {
	Band int     `json:"band"`
	Gain float32 `json:"gain"`
}

// MarshalJSON marshals the map as object array
func (e Equalizer) MarshalJSON() ([]byte, error) {
	var bands []EqBand
	for band, gain := range e {
		bands = append(bands, EqBand{
			Band: band,
			Gain: gain,
		})
	}
	return json.Marshal(bands)
}
