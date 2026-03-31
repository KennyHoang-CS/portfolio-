package game

type GameState int

const (
    StatePreload GameState = iota
    StateMenu
    StatePlaying
    StateControls
    StateRetry
    StateVictory
)