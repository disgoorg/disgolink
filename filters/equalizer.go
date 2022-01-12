package filters

import "encoding/json"

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

func (e Equalizer) SetBand(band int, gain float32) Equalizer {
	if band < 0 || band > 14 {
		return e
	}
	if gain < -0.25 {
		gain = -0.25
	}
	if gain > 1 {
		gain = 1
	}
	e[band] = gain
	return e
}

func (e Equalizer) GetBand(band int) float32 {
	if band < 0 || band > 14 {
		return -1
	}
	return e[band]
}

func (e Equalizer) Reset() Equalizer {
	for k := range e {
		delete(e, k)
	}
	return e
}

func (e Equalizer) ResetBand(band int) Equalizer {
	if band < 0 || band > 14 {
		return e
	}
	e[band] = 0
	return e
}
