package snake

import (
	"math/rand/v2"
	"snakes-ml/config"
)

// Snake represents game field and snake
type Snake struct {
	width       int
	height      int
	body        []Point
	food        Point
	obstacles   []Point
	direction   Direction
	score       int
	steps       int
	maxSteps    int
	wrapAround  bool
	dynamicSize bool
	initialSize int
}

// NewSnake creates new snake instance using config
func NewSnake(width, height int, wrapAround, dynamicSize bool) *Snake {
	s := &Snake{
		width:       width,
		height:      height,
		wrapAround:  wrapAround,
		dynamicSize: dynamicSize,
		initialSize: width,
		maxSteps:    width * height * 2,
	}
	s.Reset()
	return s
}

// Reset resets game to initial state
func (s *Snake) Reset() {
	centerX, centerY := s.width/2, s.height/2
	s.body = []Point{{X: centerX, Y: centerY}}
	s.direction = Right
	s.score = 0
	s.steps = 0
	s.obstacles = nil
	s.spawnFood()

	// Spawn initial obstacles from config
	initialObstacles := config.InitialObstaclesMin + rand.IntN(config.InitialObstaclesMax-config.InitialObstaclesMin+1)
	s.addObstacles(initialObstacles)
}

// Getters
func (s *Snake) Width() int              { return s.width }
func (s *Snake) Height() int             { return s.height }
func (s *Snake) Body() []Point           { return s.body }
func (s *Snake) Food() Point             { return s.food }
func (s *Snake) Obstacles() []Point      { return s.obstacles }
func (s *Snake) Score() int              { return s.score }
func (s *Snake) Steps() int              { return s.steps }
func (s *Snake) CurrentDirection() Direction { return s.direction }
func (s *Snake) Length() int             { return len(s.body) }

// GetOccupancy returns field occupancy percentage (0.0 - 1.0)
func (s *Snake) GetOccupancy() float64 {
	totalCells := s.width * s.height
	if totalCells == 0 {
		return 0
	}
	return float64(len(s.body)) / float64(totalCells)
}

// spawnFood generates food at random free position
func (s *Snake) spawnFood() {
	for attempt := 0; attempt < 1000; attempt++ {
		s.food = Point{X: rand.IntN(s.width), Y: rand.IntN(s.height)}
		if s.isCellFree(s.food) {
			return
		}
	}
}

// addObstacles adds random obstacles using config parameters
func (s *Snake) addObstacles(count int) {
	safeRadius := config.ObstacleSafeRadius

	for i := 0; i < count; i++ {
		maxAttempts := 200
		placed := false

		for attempt := 0; attempt < maxAttempts; attempt++ {
			obs := Point{X: rand.IntN(s.width), Y: rand.IntN(s.height)}

			// Check safe zone around snake
			tooCloseToSnake := false
			for _, segment := range s.body {
				dx := obs.X - segment.X
				dy := obs.Y - segment.Y
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx <= safeRadius && dy <= safeRadius {
					tooCloseToSnake = true
					break
				}
			}

			if tooCloseToSnake {
				continue
			}

			// Check not on food
			if s.food.Equal(obs) {
				continue
			}

			// Check cell is free
			if !s.isCellFree(obs) {
				continue
			}

			// Check doesn't create trap
			if s.wouldCreateTrap(obs) {
				continue
			}

			s.obstacles = append(s.obstacles, obs)
			placed = true
			break
		}

		if !placed {
			break
		}
	}
}

// wouldCreateTrap checks if obstacle would create inescapable trap
func (s *Snake) wouldCreateTrap(newObs Point) bool {
	// Check 4 neighboring cells
	neighbors := []Point{
		{X: newObs.X, Y: newObs.Y - 1}, // Up
		{X: newObs.X + 1, Y: newObs.Y}, // Right
		{X: newObs.X, Y: newObs.Y + 1}, // Down
		{X: newObs.X - 1, Y: newObs.Y}, // Left
	}

	blockedCount := 0

	for _, neighbor := range neighbors {
		// Check boundaries (if no wrap-around)
		if !s.wrapAround {
			if neighbor.X < 0 || neighbor.X >= s.width || neighbor.Y < 0 || neighbor.Y >= s.height {
				blockedCount++
				continue
			}
		} else {
			// Normalize coordinates with wrap-around
			if neighbor.X < 0 {
				neighbor.X = s.width - 1
			} else if neighbor.X >= s.width {
				neighbor.X = 0
			}
			if neighbor.Y < 0 {
				neighbor.Y = s.height - 1
			} else if neighbor.Y >= s.height {
				neighbor.Y = 0
			}
		}

		// Check obstacles
		for _, obs := range s.obstacles {
			if neighbor.Equal(obs) {
				blockedCount++
				break
			}
		}
	}

	// If 3 or 4 sides blocked - it's a trap
	return blockedCount >= 3
}

// isCellFree checks if cell is free
func (s *Snake) isCellFree(pos Point) bool {
	for _, segment := range s.body {
		if pos.Equal(segment) {
			return false
		}
	}
	for _, obs := range s.obstacles {
		if pos.Equal(obs) {
			return false
		}
	}
	return true
}

// Step performs one game step and returns reward and done flag
func (s *Snake) Step(action int) (float64, bool) {
	s.steps++

	newDir := Direction(action)
	if !s.direction.IsOpposite(newDir) {
		s.direction = newDir
	}

	head := s.body[0]
	delta := s.direction.ToVector()
	newHead := Point{X: head.X + delta.X, Y: head.Y + delta.Y}

	// Wrap-around mode
	if s.wrapAround {
		if newHead.X < 0 {
			newHead.X = s.width - 1
		}
		if newHead.X >= s.width {
			newHead.X = 0
		}
		if newHead.Y < 0 {
			newHead.Y = s.height - 1
		}
		if newHead.Y >= s.height {
			newHead.Y = 0
		}
	}

	reward := config.RewardStep

	// Wall collision
	if !s.wrapAround && (newHead.X < 0 || newHead.X >= s.width || newHead.Y < 0 || newHead.Y >= s.height) {
		return config.RewardDeath, true
	}

	// Self collision
	for _, segment := range s.body {
		if newHead.Equal(segment) {
			return config.RewardDeath, true
		}
	}

	// Obstacle collision
	for _, obs := range s.obstacles {
		if newHead.Equal(obs) {
			return config.RewardDeath, true
		}
	}

	// Add new head
	s.body = append([]Point{newHead}, s.body...)

	// Check food eaten
	if newHead.Equal(s.food) {
		s.score++
		reward = config.RewardFood
		s.spawnFood()

		// Check field expansion using config
		if s.dynamicSize && s.GetOccupancy() >= config.ExpansionThreshold && s.width < s.initialSize*config.MaxFieldExpansion {
			s.width += config.ExpansionIncrement
			s.height += config.ExpansionIncrement
			s.maxSteps = s.width * s.height * 2
		}

		// Add obstacles based on config interval
		if s.score%config.ObstacleAddInterval == 0 {
			s.addObstacles(1)
		}
	} else {
		// Remove tail
		s.body = s.body[:len(s.body)-1]

		// Reward for moving toward food
		oldDist := head.ManhattanDistance(s.food)
		newDist := newHead.ManhattanDistance(s.food)
		if newDist < oldDist {
			reward += config.RewardMoveToFood
		}
	}

	// Timeout
	if s.steps > s.maxSteps {
		return config.RewardDeath, true
	}

	return reward, false
}

// GetState returns current state for AI (14 parameters)
func (s *Snake) GetState() []float64 {
	head := s.body[0]
	var dangerStraight, dangerRight, dangerLeft float64

	directions := []Direction{Up, Right, Down, Left}
	currentDir := int(s.direction)

	// Danger straight
	straightDelta := directions[currentDir].ToVector()
	straightPos := Point{X: head.X + straightDelta.X, Y: head.Y + straightDelta.Y}
	if s.isDanger(straightPos) {
		dangerStraight = 1
	}

	// Danger right
	rightDir := (currentDir + 1) % 4
	rightDelta := directions[rightDir].ToVector()
	rightPos := Point{X: head.X + rightDelta.X, Y: head.Y + rightDelta.Y}
	if s.isDanger(rightPos) {
		dangerRight = 1
	}

	// Danger left
	leftDir := (currentDir + 3) % 4
	leftDelta := directions[leftDir].ToVector()
	leftPos := Point{X: head.X + leftDelta.X, Y: head.Y + leftDelta.Y}
	if s.isDanger(leftPos) {
		dangerLeft = 1
	}

	// Food direction
	var foodUp, foodRight, foodDown, foodLeft float64
	if s.food.Y < head.Y {
		foodUp = 1
	}
	if s.food.Y > head.Y {
		foodDown = 1
	}
	if s.food.X > head.X {
		foodRight = 1
	}
	if s.food.X < head.X {
		foodLeft = 1
	}

	// Current direction
	var dirUp, dirRight, dirDown, dirLeft float64
	switch s.direction {
	case Up:
		dirUp = 1
	case Right:
		dirRight = 1
	case Down:
		dirDown = 1
	case Left:
		dirLeft = 1
	}

	// Normalized distance to food
	distX := float64(s.food.X-head.X) / float64(s.width)
	distY := float64(s.food.Y-head.Y) / float64(s.height)
	occupancy := s.GetOccupancy()

	return []float64{
		dangerStraight, dangerRight, dangerLeft,
		foodUp, foodRight, foodDown, foodLeft,
		dirUp, dirRight, dirDown, dirLeft,
		distX, distY, occupancy,
	}
}

// isDanger checks if position is dangerous
func (s *Snake) isDanger(pos Point) bool {
	if !s.wrapAround {
		if pos.X < 0 || pos.X >= s.width || pos.Y < 0 || pos.Y >= s.height {
			return true
		}
	} else {
		if pos.X < 0 {
			pos.X = s.width - 1
		}
		if pos.X >= s.width {
			pos.X = 0
		}
		if pos.Y < 0 {
			pos.Y = s.height - 1
		}
		if pos.Y >= s.height {
			pos.Y = 0
		}
	}

	for _, segment := range s.body {
		if pos.Equal(segment) {
			return true
		}
	}
	for _, obs := range s.obstacles {
		if pos.Equal(obs) {
			return true
		}
	}
	return false
}
