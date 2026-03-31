package game

type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
	StateControls
    StateRetry
    StateVictory 
)