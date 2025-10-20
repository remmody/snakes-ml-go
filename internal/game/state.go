package game

// State represents game state
type State int

const (
	StateMenu State = iota
	StateTraining
	StatePlaying
	StateGameOver
)
