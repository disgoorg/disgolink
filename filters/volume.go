package filters

type Volume float32

func (v *Volume) Set(volume float32) *Volume {
	*v = Volume(volume)
	return v
}

func (v *Volume) Get() *float32 {
	return (*float32)(v)
}
