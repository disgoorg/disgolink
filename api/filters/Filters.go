package filters

type Filters struct {
	Volume     float32     `json:"volume"`
	Equalizer  *Equalizer  `json:"equalizer"`
	Timescale  *Timescale  `json:"timescale"`
	Tremolo    *Tremolo    `json:"tremolo"`
	Vibrato    *Vibrato    `json:"vibrato"`
	Rotation   *Rotation   `json:"rotation"`
	Karaoke    *Karaoke    `json:"karaoke"`
	Distortion *Distortion `json:"distortion"`
}

func (f Filters) CreateFilterManager() *Filters {
	return &Filters{
		Volume: 1.0,
		Equalizer: &Equalizer{
			EqBand{
				Band: 0,
				Gain: 0.2,
			},
		},
		Timescale: &Timescale{
			Speed: 1.0,
			Pitch: 1.0,
			Rate:  1.0,
		},
		Tremolo: &Tremolo{
			Frequency: 2.0,
			Depth:     0.5,
		},
		Vibrato: &Vibrato{
			Frequency: 2.0,
			Depth:     0.5,
		},
		Rotation: &Rotation{RotationHz: 0},
		Karaoke: &Karaoke{
			Level:       1.0,
			MonoLevel:   1.0,
			FilterBand:  220.0,
			FilterWidth: 100.0,
		},
		Distortion: &Distortion{
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
