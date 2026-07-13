package game 

type GameState int

const (
    StateTitle GameState = iota
    StateGameplay
    StateTraining
)
