package filters

// TODO add constructor func for all filters with default values

var DefaultVolume Volume = 1

func NewFilters(commitFunc func(filters *Filters)) *Filters {
	return &Filters{commitFunc: commitFunc}
}

type Filters struct {
	Volume     *Volume     `json:"volume,omitempty"`
	Equalizer  Equalizer   `json:"equalizer,omitempty"`
	Timescale  *Timescale  `json:"timescale,omitempty"`
	Tremolo    *Tremolo    `json:"tremolo,omitempty"`
	Vibrato    *Vibrato    `json:"vibrato,omitempty"`
	Rotation   *Rotation   `json:"rotation,omitempty"`
	Karaoke    *Karaoke    `json:"karaoke,omitempty"`
	Distortion *Distortion `json:"distortion,omitempty"`
	commitFunc func(filters *Filters)
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

func (f *Filters) Commit() {
	f.commitFunc(f)
}
