package snake

import (
	"math/rand/v2"
	"snakes-ml/config"
)

// Snake represents game field and snake
type Snake struct {
	width        int
	height       int
	body         []Point
	food         Point
	obstacles    []Point
	direction    Direction
	score        int
	steps        int
	maxSteps     int
	wrapAround   bool
	dynamicSize  bool
	initialSize  int
	lastPositions []Point // ✅ НОВОЕ: для отслеживания цикличности
}

// NewSnake creates new snake instance using config
func NewSnake(width, height int, wrapAround, dynamicSize bool) *Snake {
	s := &Snake{
		width:         width,
		height:        height,
		wrapAround:    wrapAround,
		dynamicSize:   dynamicSize,
		initialSize:   width,
		maxSteps:      width * height * 3,
		lastPositions: make([]Point, 0, 10),
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
	s.lastPositions = make([]Point, 0, 10)
	s.spawnFood()

	initialObstacles := config.InitialObstaclesMin + rand.IntN(config.InitialObstaclesMax-config.InitialObstaclesMin+1)
	s.addObstacles(initialObstacles)
}

// Getters
func (s *Snake) Width() int                  { return s.width }
func (s *Snake) Height() int                 { return s.height }
func (s *Snake) Body() []Point               { return s.body }
func (s *Snake) Food() Point                 { return s.food }
func (s *Snake) Obstacles() []Point          { return s.obstacles }
func (s *Snake) Score() int                  { return s.score }
func (s *Snake) Steps() int                  { return s.steps }
func (s *Snake) CurrentDirection() Direction { return s.direction }
func (s *Snake) Length() int                 { return len(s.body) }

// GetOccupancy returns field occupancy percentage
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

// addObstacles adds random obstacles
func (s *Snake) addObstacles(count int) {
	safeRadius := config.ObstacleSafeRadius

	for i := 0; i < count; i++ {
		maxAttempts := 200
		placed := false

		for attempt := 0; attempt < maxAttempts; attempt++ {
			obs := Point{X: rand.IntN(s.width), Y: rand.IntN(s.height)}

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

			if tooCloseToSnake || s.food.Equal(obs) || !s.isCellFree(obs) || s.wouldCreateTrap(obs) {
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

// wouldCreateTrap checks if obstacle creates trap
func (s *Snake) wouldCreateTrap(newObs Point) bool {
	neighbors := newObs.GetNeighbors()
	blockedCount := 0

	for _, neighbor := range neighbors {
		if !s.wrapAround {
			if neighbor.X < 0 || neighbor.X >= s.width || neighbor.Y < 0 || neighbor.Y >= s.height {
				blockedCount++
				continue
			}
		} else {
			neighbor = s.normalizePos(neighbor)
		}

		for _, obs := range s.obstacles {
			if neighbor.Equal(obs) {
				blockedCount++
				break
			}
		}
	}

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

// normalizePos normalizes position for wrap-around
func (s *Snake) normalizePos(pos Point) Point {
	if s.wrapAround {
		for pos.X < 0 {
			pos.X += s.width
		}
		for pos.X >= s.width {
			pos.X -= s.width
		}
		for pos.Y < 0 {
			pos.Y += s.height
		}
		for pos.Y >= s.height {
			pos.Y -= s.height
		}
	}
	return pos
}

// ✅ КРИТИЧЕСКИ ИСПРАВЛЕНО: Step с правильной проверкой wrap-around
func (s *Snake) Step(action int) (float64, bool) {
	s.steps++

	newDir := Direction(action)
	if !s.direction.IsOpposite(newDir) {
		s.direction = newDir
	}

	head := s.body[0]
	delta := s.direction.ToVector()
	newHead := Point{X: head.X + delta.X, Y: head.Y + delta.Y}

	// ✅ ИСПРАВЛЕНО: сначала нормализуем, ПОТОМ проверяем столкновения
	originalNewHead := newHead
	if s.wrapAround {
		newHead = s.normalizePos(newHead)
	}

	reward := config.RewardStep

	// Проверка границ (только если нет wrap-around)
	if !s.wrapAround && (originalNewHead.X < 0 || originalNewHead.X >= s.width || 
		originalNewHead.Y < 0 || originalNewHead.Y >= s.height) {
		return config.RewardDeath, true
	}

	// ✅ ИСПРАВЛЕНО: проверка столкновения с телом ПОСЛЕ нормализации
	willEatFood := newHead.Equal(s.food)
	for i, segment := range s.body {
		// Нормализуем позицию сегмента тела для корректного сравнения
		normalizedSegment := s.normalizePos(segment)
		
		// Пропускаем хвост если не едим еду
		if !willEatFood && i == len(s.body)-1 {
			continue
		}
		
		if newHead.Equal(normalizedSegment) {
			return config.RewardDeath, true
		}
	}

	// Проверка препятствий
	for _, obs := range s.obstacles {
		if newHead.Equal(obs) {
			return config.RewardDeath, true
		}
	}

	// ✅ ИСПРАВЛЕНО: проверка на циклическое движение И ИСПОЛЬЗОВАНИЕ ПЕРЕМЕННОЙ
	isCyclic := false
	for _, oldPos := range s.lastPositions {
		if newHead.Equal(oldPos) {
			isCyclic = true
			break
		}
	}

	// ✅ ИСПОЛЬЗУЕМ isCyclic: дополнительный штраф если движение циклическое
	if isCyclic {
		reward += config.RewardCycle
	}

	// Обновляем историю позиций
	s.lastPositions = append(s.lastPositions, newHead)
	if len(s.lastPositions) > 10 {
		s.lastPositions = s.lastPositions[1:]
	}

	// ✅ НОВОЕ: проверка свободного пространства вокруг
	freeSpaceCount := s.countFreeSpace(newHead)
	if freeSpaceCount >= 3 {
		reward += config.RewardFreeSpace
	} else if freeSpaceCount <= 1 {
		reward += config.RewardTrap
	}

	// Добавляем новую голову
	s.body = append([]Point{newHead}, s.body...)

	// Обработка поедания еды
	if willEatFood {
		s.score++
		reward = config.RewardFood
		s.spawnFood()

		// Очищаем историю позиций при поедании еды (новая игра)
		s.lastPositions = make([]Point, 0, 10)

		if s.dynamicSize && s.GetOccupancy() >= config.ExpansionThreshold && 
			s.width < s.initialSize*config.MaxFieldExpansion {
			s.width += config.ExpansionIncrement
			s.height += config.ExpansionIncrement
			s.maxSteps = s.width * s.height * 3
		}

		if s.score%config.ObstacleAddInterval == 0 {
			s.addObstacles(1)
		}
	} else {
		// Удаляем хвост
		s.body = s.body[:len(s.body)-1]

		// Награда за приближение
		oldDist := head.ManhattanDistance(s.food)
		newDist := newHead.ManhattanDistance(s.food)
		if newDist < oldDist {
			reward += config.RewardMoveToFood
		} else {
			reward += config.RewardMoveFromFood
		}
	}

	// Штраф за близость к телу
	minDistToBody := float64(s.width + s.height)
	for _, segment := range s.body[1:] {
		dist := newHead.ManhattanDistance(segment)
		if dist < minDistToBody {
			minDistToBody = dist
		}
	}
	if minDistToBody <= 1.5 {
		reward += config.RewardNearBody
	}

	// Таймаут
	if s.steps > s.maxSteps {
		return config.RewardDeath, true
	}

	return reward, false
}


// ✅ НОВОЕ: подсчет свободного пространства вокруг позиции
func (s *Snake) countFreeSpace(pos Point) int {
	neighbors := pos.GetNeighbors()
	count := 0

	for _, neighbor := range neighbors {
		neighbor = s.normalizePos(neighbor)
		
		if !s.wrapAround {
			if neighbor.X < 0 || neighbor.X >= s.width || 
				neighbor.Y < 0 || neighbor.Y >= s.height {
				continue
			}
		}

		isFree := true
		
		// Проверка тела
		for _, segment := range s.body {
			if s.normalizePos(neighbor).Equal(s.normalizePos(segment)) {
				isFree = false
				break
			}
		}

		// Проверка препятствий
		if isFree {
			for _, obs := range s.obstacles {
				if neighbor.Equal(obs) {
					isFree = false
					break
				}
			}
		}

		if isFree {
			count++
		}
	}

	return count
}

// ✅ РАСШИРЕНО: состояние теперь 28 параметров
func (s *Snake) GetState() []float64 {
	head := s.body[0]

	// 1-3: Опасность в трех направлениях
	var dangerStraight, dangerRight, dangerLeft float64
	directions := []Direction{Up, Right, Down, Left}
	currentDir := int(s.direction)

	straightDelta := directions[currentDir].ToVector()
	straightPos := s.normalizePos(head.Add(straightDelta))
	if s.isDanger(straightPos) {
		dangerStraight = 1
	}

	rightDir := (currentDir + 1) % 4
	rightDelta := directions[rightDir].ToVector()
	rightPos := s.normalizePos(head.Add(rightDelta))
	if s.isDanger(rightPos) {
		dangerRight = 1
	}

	leftDir := (currentDir + 3) % 4
	leftDelta := directions[leftDir].ToVector()
	leftPos := s.normalizePos(head.Add(leftDelta))
	if s.isDanger(leftPos) {
		dangerLeft = 1
	}

	// 4-7: Направление к еде
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

	// 8-11: Текущее направление
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

	// 12-13: Нормализованная дистанция до еды
	distX := float64(s.food.X-head.X) / float64(s.width)
	distY := float64(s.food.Y-head.Y) / float64(s.height)

	// 14: Заполненность
	occupancy := s.GetOccupancy()

	// 15-18: Расстояние до препятствий
	distObsUp := s.getDistanceToObstacle(head, Up) / float64(s.height)
	distObsRight := s.getDistanceToObstacle(head, Right) / float64(s.width)
	distObsDown := s.getDistanceToObstacle(head, Down) / float64(s.height)
	distObsLeft := s.getDistanceToObstacle(head, Left) / float64(s.width)

	// 19-22: Расстояние до своего тела
	distBodyUp := s.getDistanceToBody(head, Up) / float64(s.height)
	distBodyRight := s.getDistanceToBody(head, Right) / float64(s.width)
	distBodyDown := s.getDistanceToBody(head, Down) / float64(s.height)
	distBodyLeft := s.getDistanceToBody(head, Left) / float64(s.width)

	// 23: Длина змейки
	lengthNorm := float64(len(s.body)) / float64(s.width*s.height)

	// 24: Манхэттенское расстояние до еды
	manhattanDist := head.ManhattanDistance(s.food) / float64(s.width+s.height)

	// ✅ 25-28: НОВОЕ - свободное пространство в 4 направлениях
	freeSpaceUp := float64(s.countFreeSpaceInDirection(head, Up)) / 4.0
	freeSpaceRight := float64(s.countFreeSpaceInDirection(head, Right)) / 4.0
	freeSpaceDown := float64(s.countFreeSpaceInDirection(head, Down)) / 4.0
	freeSpaceLeft := float64(s.countFreeSpaceInDirection(head, Left)) / 4.0

	return []float64{
		dangerStraight, dangerRight, dangerLeft,
		foodUp, foodRight, foodDown, foodLeft,
		dirUp, dirRight, dirDown, dirLeft,
		distX, distY, occupancy,
		distObsUp, distObsRight, distObsDown, distObsLeft,
		distBodyUp, distBodyRight, distBodyDown, distBodyLeft,
		lengthNorm, manhattanDist,
		freeSpaceUp, freeSpaceRight, freeSpaceDown, freeSpaceLeft, // ✅ НОВОЕ
	}
}

// ✅ НОВОЕ: подсчет свободного пространства в направлении
func (s *Snake) countFreeSpaceInDirection(from Point, dir Direction) int {
	delta := dir.ToVector()
	pos := s.normalizePos(from.Add(delta))
	return s.countFreeSpace(pos)
}

// getDistanceToObstacle получает расстояние до препятствия
func (s *Snake) getDistanceToObstacle(from Point, dir Direction) float64 {
	delta := dir.ToVector()
	current := from
	distance := 0.0

	for i := 0; i < max(s.width, s.height); i++ {
		current = s.normalizePos(Point{X: current.X + delta.X, Y: current.Y + delta.Y})
		distance++

		if !s.wrapAround {
			if current.X < 0 || current.X >= s.width || current.Y < 0 || current.Y >= s.height {
				return distance
			}
		}

		for _, obs := range s.obstacles {
			if current.Equal(obs) {
				return distance
			}
		}

		if s.wrapAround && i >= max(s.width, s.height) {
			return float64(max(s.width, s.height))
		}
	}

	return float64(max(s.width, s.height))
}

// getDistanceToBody получает расстояние до своего тела
func (s *Snake) getDistanceToBody(from Point, dir Direction) float64 {
	delta := dir.ToVector()
	current := from
	distance := 0.0

	for i := 0; i < max(s.width, s.height); i++ {
		current = s.normalizePos(Point{X: current.X + delta.X, Y: current.Y + delta.Y})
		distance++

		if !s.wrapAround {
			if current.X < 0 || current.X >= s.width || current.Y < 0 || current.Y >= s.height {
				return distance
			}
		}

		for _, segment := range s.body[1:] {
			// ✅ ИСПРАВЛЕНО: сравниваем нормализованные позиции
			if s.normalizePos(current).Equal(s.normalizePos(segment)) {
				return distance
			}
		}

		if s.wrapAround && i >= max(s.width, s.height) {
			return float64(max(s.width, s.height))
		}
	}

	return float64(max(s.width, s.height))
}

// isDanger проверяет опасность позиции
func (s *Snake) isDanger(pos Point) bool {
	pos = s.normalizePos(pos)

	if !s.wrapAround {
		if pos.X < 0 || pos.X >= s.width || pos.Y < 0 || pos.Y >= s.height {
			return true
		}
	}

	// ✅ ИСПРАВЛЕНО: сравниваем нормализованные позиции
	for _, segment := range s.body {
		if s.normalizePos(pos).Equal(s.normalizePos(segment)) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
