package filters

type Timescale struct {
	Speed float32 `json:"speed,omitempty"`
	Pitch float32 `json:"pitch,omitempty"`
	Rate  float32 `json:"rate,omitempty"`
}
