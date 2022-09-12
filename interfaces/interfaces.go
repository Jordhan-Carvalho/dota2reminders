package interfaces


type mapEvent struct {
	Name      string `json:"name"`
	ClockTime int    `json:"clock_time"`
	GameState string `json:"game_state"` // Would there be a way to make this an enum?
}

type GameEvents struct {
	Map mapEvent `json:"map"`
}
