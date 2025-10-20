package snake

// Direction represents snake movement direction
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

// IsOpposite checks if direction is opposite to another
func (d Direction) IsOpposite(other Direction) bool {
	return (d == Up && other == Down) ||
		(d == Down && other == Up) ||
		(d == Left && other == Right) ||
		(d == Right && other == Left)
}

// ToVector converts direction to vector
func (d Direction) ToVector() Point {
	switch d {
	case Up:
		return Point{0, -1}
	case Down:
		return Point{0, 1}
	case Left:
		return Point{-1, 0}
	case Right:
		return Point{1, 0}
	default:
		return Point{0, 0}
	}
}
