package filters

type Karaoke struct {
	Level       float32 `json:"level"`
	MonoLevel   float32 `json:"monoLevel"`
	FilterBand  float32 `json:"filterBand"`
	FilterWidth float32 `json:"filterWidth"`
}
