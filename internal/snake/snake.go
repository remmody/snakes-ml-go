package snake

import "math/rand/v2"

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

func (s *Snake) Reset() {
	centerX, centerY := s.width/2, s.height/2
	s.body = []Point{{X: centerX, Y: centerY}}
	s.direction = Right
	s.score = 0
	s.steps = 0
	s.obstacles = nil
	s.spawnFood()
	initialObstacles := 3 + rand.IntN(3)
	s.addObstacles(initialObstacles)
}

// Геттеры
func (s *Snake) Width() int       { return s.width }
func (s *Snake) Height() int      { return s.height }
func (s *Snake) Body() []Point    { return s.body }
func (s *Snake) Food() Point      { return s.food }
func (s *Snake) Obstacles() []Point { return s.obstacles }
func (s *Snake) Score() int       { return s.score }
func (s *Snake) Steps() int       { return s.steps }
func (s *Snake) CurrentDirection() Direction { return s.direction }
func (s *Snake) Length() int      { return len(s.body) }

// ✅ ИСПРАВЛЕНИЕ: Добавлен метод GetOccupancy
func (s *Snake) GetOccupancy() float64 {
	totalCells := s.width * s.height
	if totalCells == 0 {
		return 0
	}
	return float64(len(s.body)) / float64(totalCells)
}

func (s *Snake) spawnFood() {
	for attempt := 0; attempt < 1000; attempt++ {
		s.food = Point{X: rand.IntN(s.width), Y: rand.IntN(s.height)}
		if s.isCellFree(s.food) {
			return
		}
	}
}

func (s *Snake) addObstacles(count int) {
	centerX, centerY := s.width/2, s.height/2
	for i := 0; i < count; i++ {
		for attempt := 0; attempt < 100; attempt++ {
			obs := Point{X: rand.IntN(s.width), Y: rand.IntN(s.height)}
			dx, dy := obs.X-centerX, obs.Y-centerY
			if dx < 0 { dx = -dx }
			if dy < 0 { dy = -dy }
			if len(s.body) == 1 && dx <= 2 && dy <= 2 {
				continue
			}
			if s.isCellFree(obs) && !s.food.Equal(obs) {
				s.obstacles = append(s.obstacles, obs)
				break
			}
		}
	}
}

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

func (s *Snake) Step(action int) (float64, bool) {
	s.steps++
	newDir := Direction(action)
	if !s.direction.IsOpposite(newDir) {
		s.direction = newDir
	}

	head := s.body[0]
	delta := s.direction.ToVector()
	newHead := Point{X: head.X + delta.X, Y: head.Y + delta.Y}

	if s.wrapAround {
		if newHead.X < 0 { newHead.X = s.width - 1 }
		if newHead.X >= s.width { newHead.X = 0 }
		if newHead.Y < 0 { newHead.Y = s.height - 1 }
		if newHead.Y >= s.height { newHead.Y = 0 }
	}

	reward := -0.01

	if !s.wrapAround && (newHead.X < 0 || newHead.X >= s.width || newHead.Y < 0 || newHead.Y >= s.height) {
		return -10.0, true
	}

	for _, segment := range s.body {
		if newHead.Equal(segment) {
			return -10.0, true
		}
	}

	for _, obs := range s.obstacles {
		if newHead.Equal(obs) {
			return -10.0, true
		}
	}

	s.body = append([]Point{newHead}, s.body...)

	if newHead.Equal(s.food) {
		s.score++
		reward = 10.0
		s.spawnFood()
		if s.dynamicSize && s.GetOccupancy() >= 0.9 && s.width < s.initialSize*3 {
			s.width += 2
			s.height += 2
			s.maxSteps = s.width * s.height * 2
		}
		s.addObstacles(1)
	} else {
		s.body = s.body[:len(s.body)-1]
		oldDist := head.ManhattanDistance(s.food)
		newDist := newHead.ManhattanDistance(s.food)
		if newDist < oldDist {
			reward += 0.1
		}
	}

	if s.steps > s.maxSteps {
		return -10.0, true
	}

	return reward, false
}

func (s *Snake) GetState() []float64 {
	head := s.body[0]
	var dangerStraight, dangerRight, dangerLeft float64

	directions := []Direction{Up, Right, Down, Left}
	currentDir := int(s.direction)

	straightDelta := directions[currentDir].ToVector()
	straightPos := Point{X: head.X + straightDelta.X, Y: head.Y + straightDelta.Y}
	if s.isDanger(straightPos) { dangerStraight = 1 }

	rightDir := (currentDir + 1) % 4
	rightDelta := directions[rightDir].ToVector()
	rightPos := Point{X: head.X + rightDelta.X, Y: head.Y + rightDelta.Y}
	if s.isDanger(rightPos) { dangerRight = 1 }

	leftDir := (currentDir + 3) % 4
	leftDelta := directions[leftDir].ToVector()
	leftPos := Point{X: head.X + leftDelta.X, Y: head.Y + leftDelta.Y}
	if s.isDanger(leftPos) { dangerLeft = 1 }

	var foodUp, foodRight, foodDown, foodLeft float64
	if s.food.Y < head.Y { foodUp = 1 }
	if s.food.Y > head.Y { foodDown = 1 }
	if s.food.X > head.X { foodRight = 1 }
	if s.food.X < head.X { foodLeft = 1 }

	var dirUp, dirRight, dirDown, dirLeft float64
	switch s.direction {
	case Up: dirUp = 1
	case Right: dirRight = 1
	case Down: dirDown = 1
	case Left: dirLeft = 1
	}

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

func (s *Snake) isDanger(pos Point) bool {
	if !s.wrapAround {
		if pos.X < 0 || pos.X >= s.width || pos.Y < 0 || pos.Y >= s.height {
			return true
		}
	} else {
		if pos.X < 0 { pos.X = s.width - 1 }
		if pos.X >= s.width { pos.X = 0 }
		if pos.Y < 0 { pos.Y = s.height - 1 }
		if pos.Y >= s.height { pos.Y = 0 }
	}

	for _, segment := range s.body {
		if pos.Equal(segment) { return true }
	}
	for _, obs := range s.obstacles {
		if pos.Equal(obs) { return true }
	}
	return false
}
