package filters

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
