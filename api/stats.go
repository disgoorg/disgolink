package api

type Stats struct {
	Players        int         `json:"players"`
	PlayingPlayers int         `json:"playingPlayers"`
	Uptime         int         `json:"uptime"`
	Memory         *Memory     `json:"memory"`
	Cpu            *Cpu        `json:"cpu"`
	FrameStats     *FrameStats `json:"frameStats"`
}

type Memory struct {
	Free       int `json:"free"`
	Used       int `json:"used"`
	Allocated  int `json:"allocated"`
	Reservable int `json:"reservable"`
}

type Cpu struct {
	Cores        int `json:"cores"`
	SystemLoad   int `json:"systemLoad"`
	LavalinkLoad int `json:"lavalinkLoad"`
}

type FrameStats struct {
	Sent    int `json:"sent"`
	Nulled  int `json:"nulled"`
	Deficit int `json:"deficit"`
}
