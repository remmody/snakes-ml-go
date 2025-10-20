package snake

import "math"

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

// EuclideanDistance calculates Euclidean distance
func (p Point) EuclideanDistance(other Point) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Add adds two points
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// GetNeighbors returns 4 neighboring points
func (p Point) GetNeighbors() []Point {
	return []Point{
		{X: p.X, Y: p.Y - 1},     // Up
		{X: p.X + 1, Y: p.Y},     // Right
		{X: p.X, Y: p.Y + 1},     // Down
		{X: p.X - 1, Y: p.Y},     // Left
	}
}
