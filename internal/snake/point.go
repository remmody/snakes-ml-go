package snake

// Point represents coordinates on grid
type Point struct {
	X, Y int
}

// Equal checks if two points are equal
func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// ManhattanDistance calculates Manhattan distance between points
func (p Point) ManhattanDistance(other Point) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
