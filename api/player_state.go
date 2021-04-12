package api

type PlayerState struct {
	Duration  int  `json:"time"`
	Position  int  `json:"position"`
	Connected bool `json:"connected"`
}
