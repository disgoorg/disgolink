package filters

type Equalizer []EqBand

type EqBand struct {
	Band int     `json:"band"`
	Gain float32 `json:"gain"`
}
