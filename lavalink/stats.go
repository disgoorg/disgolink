package lavalink

type Stats struct {
	Players        int         `json:"players"`
	PlayingPlayers int         `json:"playingPlayers"`
	Uptime         Duration    `json:"uptime"`
	Memory         Memory      `json:"memory"`
	CPU            CPU         `json:"cpu"`
	FrameStats     *FrameStats `json:"frameStats"`
}

func (s Stats) Better(stats Stats) bool {
	if s.CPU.Cores <= 0 {
		return false
	}
	if stats.CPU.Cores <= 0 {
		return true
	}
	return s.CPU.SystemLoad/float64(s.CPU.Cores) < stats.CPU.SystemLoad/float64(stats.CPU.Cores)
}

type Memory struct {
	Free       int `json:"free"`
	Used       int `json:"used"`
	Allocated  int `json:"allocated"`
	Reservable int `json:"reservable"`
}

type CPU struct {
	Cores        int     `json:"cores"`
	SystemLoad   float64 `json:"systemLoad"`
	LavalinkLoad float64 `json:"lavalinkLoad"`
}

type FrameStats struct {
	Sent    int `json:"sent"`
	Nulled  int `json:"nulled"`
	Deficit int `json:"deficit"`
}
