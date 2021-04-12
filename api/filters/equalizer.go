package filters

type Volume struct{
	volume int
}

func (v *Volume) setVolume(amount int)  {
	v.volume = amount
}

