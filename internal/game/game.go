package game

import (
	"fmt"
	"time"

	"snakes-ml/internal/ai"
	"snakes-ml/internal/snake"
	"snakes-ml/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	screenWidth     int
	screenHeight    int
	state           State
	snake           *snake.Snake
	agent           *ai.Agent
	renderer        *Renderer
	generation      int
	episode         int
	maxEpisodes     int
	bestScore       int
	currentScore    int
	recentScores    []int
	windowSize      int
	trainingMode    bool
	autoRestart     bool
	frameCount      int
	speedMultiplier float64
	statsText       string
	lastUpdateTime  time.Time
	lastMapSize     string
}

func NewGame(screenWidth, screenHeight int) *Game {
	g := &Game{
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
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

	g.renderer = NewRenderer(screenWidth, screenHeight)

	config := ai.DefaultConfig()
	g.agent = ai.NewAgent(14, 4, config)

	if err := g.agent.LoadModel("snake_ai_model_best.json"); err == nil {
		fmt.Println("✅ Loaded existing model")
	} else {
		fmt.Println("🆕 Created new model")
	}

	return g
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
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.speedMultiplier = 50.0
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

		if g.agent.ReplayBufferSize() >= 64 {
			g.agent.Train()
		}

		g.currentScore = g.snake.Score()

		if done {
			g.handleEpisodeEnd()

			if g.episode > g.maxEpisodes {
				g.agent.SaveModel("snake_ai_model_final.json")
				fmt.Println("\n✅ Training completed!")
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

		oldEpsilon := g.agent.Epsilon()
		g.agent.SetEpsilon(0)
		action := g.agent.SelectAction(g.snake.GetState())
		g.agent.SetEpsilon(oldEpsilon)

		_, done := g.snake.Step(action)
		g.currentScore = g.snake.Score()

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
	g.snake = snake.NewSnake(20, 15, true, true)
	g.currentScore = 0
	g.lastMapSize = fmt.Sprintf("%dx%d", g.snake.Width(), g.snake.Height())
}

func (g *Game) handleEpisodeEnd() {
	score := g.snake.Score()
	g.recentScores = append(g.recentScores, score)
	if len(g.recentScores) > g.windowSize {
		g.recentScores = g.recentScores[1:]
	}

	if score > g.bestScore {
		g.bestScore = score
		g.agent.SaveModel("snake_ai_model_best.json")
		fmt.Printf("🏆 New record: %d (episode %d)\n", score, g.episode)
	}

	if g.episode%500 == 0 {
		filename := fmt.Sprintf("snake_ai_model_gen%d_ep%d.json", g.generation, g.episode)
		g.agent.SaveModel(filename)
		fmt.Printf("💾 Checkpoint saved: %s\n", filename)
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

	occupancy := 0.0
	if g.snake != nil {
		occupancy = g.snake.GetOccupancy() * 100
	}

	g.statsText = fmt.Sprintf(
		"Generation: %d | Episode: %d/%d\n"+
			"Score: %d | Avg: %.2f | Best: %d\n"+
			"Epsilon: %.4f | Buffer: %d\n"+
			"Map: %s | Occupancy: %.1f%%\n"+
			"Obstacles: %d | Speed: %.0fx",
		g.generation, g.episode, g.maxEpisodes,
		g.currentScore, avgScore, g.bestScore,
		g.agent.Epsilon(), g.agent.ReplayBufferSize(),
		g.lastMapSize, occupancy,
		len(g.snake.Obstacles()), g.speedMultiplier,
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ui.Background)

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
	// ✅ Адаптировано под 1280x720
	title := "AI SNAKE GAME - Deep Q-Learning"
	ebitenutil.DebugPrintAt(screen, title, g.screenWidth/2-200, 100)

	subtitle := "Self-learning snake powered by neural networks"
	ebitenutil.DebugPrintAt(screen, subtitle, g.screenWidth/2-220, 140)

	instructions := []string{
		"",
		"═════════════════════════════════════════════════",
		"",
		"[SPACE] - Start Training",
		"[P]     - Play with Trained AI",
		"[Q]     - Quit",
		"",
		"═════════════════════════════════════════════════",
		"",
		fmt.Sprintf("📊 Best Score: %d", g.bestScore),
		fmt.Sprintf("🎓 Episodes Trained: %d", g.episode-1),
		"",
		"Features:",
		"• Random yellow obstacles",
		"• Auto map expansion at 90% occupancy",
		"• Wrap-around boundaries",
		"• Deep Q-Learning with Experience Replay",
	}

	y := 200
	for _, line := range instructions {
		ebitenutil.DebugPrintAt(screen, line, g.screenWidth/2-250, y)
		y += 30
	}

	info := "Controls: [1] 1x [2] 5x [3] 10x [4] 50x speed | [ESC] Menu"
	ebitenutil.DebugPrintAt(screen, info, 30, g.screenHeight-40)
}

func (g *Game) drawTraining(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	if g.statsText != "" {
		vector.FillRect(screen, 5, 5, 450, 110, ui.TextBg, false)
		ebitenutil.DebugPrintAt(screen, g.statsText, 10, 10)
	}

	g.renderer.DrawProgressBar(screen, float64(g.episode)/float64(g.maxEpisodes), g.episode, g.maxEpisodes)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	vector.FillRect(screen, 5, 5, 300, 30, ui.TextBg, false)
	scoreText := fmt.Sprintf("Score: %d | Best: %d", g.currentScore, g.bestScore)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 10)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	vector.FillRect(screen, 0, 0, float32(g.screenWidth), float32(g.screenHeight), ui.TextBg, false)

	boxX := float32(g.screenWidth/2 - 180)
	boxY := float32(g.screenHeight/2 - 120)
	vector.FillRect(screen, boxX, boxY, 360, 240, ui.Background, false)
	vector.StrokeRect(screen, boxX, boxY, 360, 240, 3, ui.SnakeHead, false)

	ebitenutil.DebugPrintAt(screen, "GAME OVER", g.screenWidth/2-50, g.screenHeight/2-80)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.currentScore), g.screenWidth/2-40, g.screenHeight/2-40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best: %d", g.bestScore), g.screenWidth/2-35, g.screenHeight/2-10)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Play Again", g.screenWidth/2-75, g.screenHeight/2+30)
	ebitenutil.DebugPrintAt(screen, "[ESC] Main Menu", g.screenWidth/2-65, g.screenHeight/2+60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
