package game

type State int

const (
	StateMenu State = iota
	StateTraining
	StatePlaying
	StateGameOver
)
