package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameState int

const (
	StateMenu GameState = iota
	StateTraining
	StatePlaying
	StateGameOver
)

type Game struct {
	state          GameState
	snake          *Snake
	agent          *DQNAgent
	generation     int
	episode        int
	maxEpisodes    int
	bestScore      int
	currentScore   int
	recentScores   []int
	windowSize     int
	trainingMode   bool
	autoRestart    bool
	frameCount     int
	speedMultiplier float64
	statsText      string
	lastUpdateTime time.Time
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	
	game := &Game{
		state:           StateMenu,
		generation:      1,
		episode:         1,
		maxEpisodes:     5000,
		windowSize:      100,
		recentScores:    make([]int, 0, 100),
		trainingMode:    true,
		autoRestart:     true,
		speedMultiplier: 1.0,
		lastUpdateTime:  time.Now(),
	}
	
	stateSize := 14
	actionSize := 4
	config := AgentConfig{
		LearningRate: 0.001,
		BufferSize:   50000,
		EpsilonStart: 1.0,
		EpsilonMin:   0.01,
		EpsilonDecay: 0.995,
		Gamma:        0.95,
		BatchSize:    64,
		UpdateFreq:   100,
	}
	
	game.agent = NewDQNAgent(stateSize, actionSize, config)
	
	if err := game.agent.LoadModel("snake_ai_model_best.json"); err == nil {
		game.agent.epsilon = 0.1
	}
	
	return game
}

func (g *Game) Update() error {
	g.frameCount++
	
	switch g.state {
	case StateMenu:
		return g.updateMenu()
	case StateTraining:
		return g.updateTraining()
	case StatePlaying:
		return g.updatePlaying()
	case StateGameOver:
		return g.updateGameOver()
	}
	
	return nil
}

func (g *Game) updateMenu() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.startTraining()
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.startPlaying()
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) updateTraining() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.speedMultiplier = 1.0
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.speedMultiplier = 5.0
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.speedMultiplier = 10.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
		return nil
	}
	
	for i := 0; i < int(g.speedMultiplier); i++ {
		if g.snake == nil {
			g.startNewEpisode()
		}
		
		state := g.snake.GetState()
		action := g.agent.SelectAction(state)
		reward, done := g.snake.Step(action)
		nextState := g.snake.GetState()
		
		g.agent.Remember(state, action, reward, nextState, done)
		
		if g.agent.replayBuffer.Size() >= g.agent.batchSize {
			g.agent.Train()
		}
		
		g.currentScore = g.snake.score
		
		if done {
			g.handleEpisodeEnd()
			
			if g.episode > g.maxEpisodes {
				g.agent.SaveModel("snake_ai_model_final.json")
				g.state = StateMenu
				return nil
			}
		}
	}
	
	if time.Since(g.lastUpdateTime) > time.Second {
		g.updateStatsText()
		g.lastUpdateTime = time.Now()
	}
	
	return nil
}

func (g *Game) updatePlaying() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
		return nil
	}
	
	if g.frameCount%5 == 0 {
		if g.snake == nil {
			g.startNewEpisode()
		}
		
		oldEpsilon := g.agent.epsilon
		g.agent.epsilon = 0
		action := g.agent.SelectAction(g.snake.GetState())
		g.agent.epsilon = oldEpsilon
		
		_, done := g.snake.Step(action)
		g.currentScore = g.snake.score
		
		if done {
			g.state = StateGameOver
		}
	}
	
	return nil
}

func (g *Game) updateGameOver() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.startPlaying()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
	return nil
}

func (g *Game) startTraining() {
	g.state = StateTraining
	g.episode = 1
	g.startNewEpisode()
}

func (g *Game) startPlaying() {
	g.state = StatePlaying
	g.trainingMode = false
	g.startNewEpisode()
}

func (g *Game) startNewEpisode() {
	g.snake = NewSnake(20, 15, true, true)
	g.currentScore = 0
}

func (g *Game) handleEpisodeEnd() {
	score := g.snake.score
	g.recentScores = append(g.recentScores, score)
	if len(g.recentScores) > g.windowSize {
		g.recentScores = g.recentScores[1:]
	}
	
	if score > g.bestScore {
		g.bestScore = score
		g.agent.SaveModel("snake_ai_model_best.json")
	}
	
	if g.episode%500 == 0 {
		filename := fmt.Sprintf("snake_ai_model_gen%d_ep%d.json", g.generation, g.episode)
		g.agent.SaveModel(filename)
	}
	
	g.episode++
	
	if g.autoRestart {
		g.startNewEpisode()
	}
}

func (g *Game) updateStatsText() {
	avgScore := 0.0
	if len(g.recentScores) > 0 {
		for _, s := range g.recentScores {
			avgScore += float64(s)
		}
		avgScore /= float64(len(g.recentScores))
	}
	
	g.statsText = fmt.Sprintf(
		"Generation: %d | Episode: %d/%d\n"+
		"Score: %d | Avg: %.2f | Best: %d\n"+
		"Epsilon: %.4f | Buffer: %d\n"+
		"Speed: %.0fx",
		g.generation, g.episode, g.maxEpisodes,
		g.currentScore, avgScore, g.bestScore,
		g.agent.epsilon, g.agent.replayBuffer.Size(),
		g.speedMultiplier,
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})
	
	switch g.state {
	case StateMenu:
		g.drawMenu(screen)
	case StateTraining:
		g.drawTraining(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	title := "AI SNAKE GAME - Deep Q-Learning"
	ebitenutil.DebugPrintAt(screen, title, screenWidth/2-200, 100)
	
	instructions := []string{
		"",
		"[SPACE] - Start Training",
		"[P]     - Play with AI",
		"[Q]     - Quit",
		"",
		fmt.Sprintf("Best Score: %d", g.bestScore),
		fmt.Sprintf("Episodes Trained: %d", g.episode-1),
	}
	
	y := 200
	for _, line := range instructions {
		ebitenutil.DebugPrintAt(screen, line, screenWidth/2-100, y)
		y += 30
	}
	
	info := "Controls during training: [1] 1x [2] 5x [3] 10x speed | [ESC] Menu"
	ebitenutil.DebugPrintAt(screen, info, 20, screenHeight-30)
}

func (g *Game) drawTraining(screen *ebiten.Image) {
	if g.snake != nil {
		g.drawSnake(screen)
	}
	
	if g.statsText != "" {
		ebitenutil.DebugPrintAt(screen, g.statsText, 10, 10)
	}
	
	g.drawProgressBar(screen, float64(g.episode)/float64(g.maxEpisodes))
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	if g.snake != nil {
		g.drawSnake(screen)
	}
	
	scoreText := fmt.Sprintf("Score: %d | Best: %d", g.currentScore, g.bestScore)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 10)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	if g.snake != nil {
		g.drawSnake(screen)
	}
	
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, 
		color.RGBA{0, 0, 0, 150}, false)
	
	ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-50, screenHeight/2-50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.currentScore), screenWidth/2-40, screenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best: %d", g.bestScore), screenWidth/2-40, screenHeight/2+10)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Play Again | [ESC] Menu", screenWidth/2-120, screenHeight/2+50)
}

func (g *Game) drawSnake(screen *ebiten.Image) {
	gridX := 20
	gridY := 100
	cellSize := 20
	
	for y := 0; y < g.snake.height; y++ {
		for x := 0; x < g.snake.width; x++ {
			posX := float32(gridX + x*cellSize)
			posY := float32(gridY + y*cellSize)
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize),
				1, color.RGBA{50, 50, 60, 255}, false)
		}
	}
	
	for _, obs := range g.snake.obstacles {
		posX := float32(gridX + obs.X*cellSize)
		posY := float32(gridY + obs.Y*cellSize)
		vector.DrawFilledRect(screen, posX, posY, float32(cellSize), float32(cellSize),
			color.RGBA{100, 50, 50, 255}, false)
	}
	
	posX := float32(gridX + g.snake.food.X*cellSize)
	posY := float32(gridY + g.snake.food.Y*cellSize)
	vector.DrawFilledRect(screen, posX, posY, float32(cellSize), float32(cellSize),
		color.RGBA{255, 100, 100, 255}, false)
	
	for i, segment := range g.snake.snake {
		posX := float32(gridX + segment.X*cellSize)
		posY := float32(gridY + segment.Y*cellSize)
		
		var col color.RGBA
		if i == 0 {
			col = color.RGBA{100, 255, 100, 255}
		} else {
			col = color.RGBA{100, 200, 100, 255}
		}
		
		vector.DrawFilledRect(screen, posX, posY, float32(cellSize), float32(cellSize), col, false)
	}
	
	infoText := fmt.Sprintf("Length: %d | Steps: %d | Map: %dx%d",
		len(g.snake.snake), g.snake.steps, g.snake.width, g.snake.height)
	ebitenutil.DebugPrintAt(screen, infoText, gridX, gridY-30)
}

func (g *Game) drawProgressBar(screen *ebiten.Image, progress float64) {
	barX := float32(10)
	barY := float32(screenHeight - 40)
	barWidth := float32(screenWidth - 20)
	barHeight := float32(20)
	
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2,
		color.RGBA{100, 100, 120, 255}, false)
	
	fillWidth := barWidth * float32(progress)
	vector.DrawFilledRect(screen, barX, barY, fillWidth, barHeight,
		color.RGBA{100, 200, 100, 255}, false)
	
	progressText := fmt.Sprintf("Progress: %.1f%%", progress*100)
	ebitenutil.DebugPrintAt(screen, progressText, int(barX)+5, int(barY)+5)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
