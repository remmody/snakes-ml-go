package main

import (
	"math"
	"math/rand"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Point struct {
	X, Y int
}

type Snake struct {
	width       int
	height      int
	snake       []Point
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
	s.snake = []Point{{centerX, centerY}}
	s.direction = Right
	s.score = 0
	s.steps = 0
	s.obstacles = nil
	s.SpawnFood()
}

func (s *Snake) SpawnFood() {
	for {
		s.food = Point{X: rand.Intn(s.width), Y: rand.Intn(s.height)}
		valid := true
		for _, segment := range s.snake {
			if s.food == segment {
				valid = false
				break
			}
		}
		for _, obs := range s.obstacles {
			if s.food == obs {
				valid = false
				break
			}
		}
		if valid {
			break
		}
	}
}

func (s *Snake) AddObstacles(count int) {
	for i := 0; i < count; i++ {
		for {
			obs := Point{X: rand.Intn(s.width), Y: rand.Intn(s.height)}
			valid := true
			for _, segment := range s.snake {
				if obs == segment {
					valid = false
					break
				}
			}
			if obs == s.food {
				valid = false
			}
			for _, o := range s.obstacles {
				if obs == o {
					valid = false
					break
				}
			}
			if valid {
				s.obstacles = append(s.obstacles, obs)
				break
			}
		}
	}
}

func (s *Snake) Step(action int) (float64, bool) {
	s.steps++
	newDir := Direction(action)
	if !((s.direction == Up && newDir == Down) || (s.direction == Down && newDir == Up) ||
		(s.direction == Left && newDir == Right) || (s.direction == Right && newDir == Left)) {
		s.direction = newDir
	}

	head := s.snake[0]
	var newHead Point
	switch s.direction {
	case Up:
		newHead = Point{head.X, head.Y - 1}
	case Down:
		newHead = Point{head.X, head.Y + 1}
	case Left:
		newHead = Point{head.X - 1, head.Y}
	case Right:
		newHead = Point{head.X + 1, head.Y}
	}

	if s.wrapAround {
		if newHead.X < 0 {
			newHead.X = s.width - 1
		} else if newHead.X >= s.width {
			newHead.X = 0
		}
		if newHead.Y < 0 {
			newHead.Y = s.height - 1
		} else if newHead.Y >= s.height {
			newHead.Y = 0
		}
	}

	reward := -0.01
	if !s.wrapAround && (newHead.X < 0 || newHead.X >= s.width || newHead.Y < 0 || newHead.Y >= s.height) {
		return -10.0, true
	}

	for _, segment := range s.snake {
		if newHead == segment {
			return -10.0, true
		}
	}

	for _, obs := range s.obstacles {
		if newHead == obs {
			return -5.0, true
		}
	}

	s.snake = append([]Point{newHead}, s.snake...)

	if newHead == s.food {
		s.score++
		reward = 10.0
		s.SpawnFood()
		if s.dynamicSize && s.score%5 == 0 && s.width < s.initialSize*3 {
			s.width++
			s.height++
			s.maxSteps = s.width * s.height * 2
		}
		if s.score%5 == 0 {
			s.AddObstacles(1)
		}
	} else {
		s.snake = s.snake[:len(s.snake)-1]
		oldDist := math.Abs(float64(head.X-s.food.X)) + math.Abs(float64(head.Y-s.food.Y))
		newDist := math.Abs(float64(newHead.X-s.food.X)) + math.Abs(float64(newHead.Y-s.food.Y))
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
	head := s.snake[0]
	var dangerStraight, dangerRight, dangerLeft float64
	directions := []Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	currentDir := int(s.direction)
	
	straightPos := Point{head.X + directions[currentDir].X, head.Y + directions[currentDir].Y}
	if s.IsDanger(straightPos) {
		dangerStraight = 1
	}
	
	rightDir := (currentDir + 1) % 4
	rightPos := Point{head.X + directions[rightDir].X, head.Y + directions[rightDir].Y}
	if s.IsDanger(rightPos) {
		dangerRight = 1
	}
	
	leftDir := (currentDir + 3) % 4
	leftPos := Point{head.X + directions[leftDir].X, head.Y + directions[leftDir].Y}
	if s.IsDanger(leftPos) {
		dangerLeft = 1
	}
	
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
	
	distX := float64(s.food.X-head.X) / float64(s.width)
	distY := float64(s.food.Y-head.Y) / float64(s.height)
	snakeLength := float64(len(s.snake)) / float64(s.width*s.height)
	
	return []float64{
		dangerStraight, dangerRight, dangerLeft,
		foodUp, foodRight, foodDown, foodLeft,
		dirUp, dirRight, dirDown, dirLeft,
		distX, distY, snakeLength,
	}
}

func (s *Snake) IsDanger(pos Point) bool {
	if !s.wrapAround {
		if pos.X < 0 || pos.X >= s.width || pos.Y < 0 || pos.Y >= s.height {
			return true
		}
	} else {
		if pos.X < 0 {
			pos.X = s.width - 1
		} else if pos.X >= s.width {
			pos.X = 0
		}
		if pos.Y < 0 {
			pos.Y = s.height - 1
		} else if pos.Y >= s.height {
			pos.Y = 0
		}
	}
	
	for _, segment := range s.snake {
		if pos == segment {
			return true
		}
	}
	for _, obs := range s.obstacles {
		if pos == obs {
			return true
		}
	}
	return false
}
