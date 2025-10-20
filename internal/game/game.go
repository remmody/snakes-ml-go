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
		maxEpisodes:     1000000,
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
		"Gen: %d | Ep: %d/%d | Score: %d | Avg: %.1f | Best: %d\n"+
			"Epsilon: %.3f | Buffer: %d | Map: %s | Occ: %.0f%% | Obs: %d | x%.0f",
		g.generation, g.episode, g.maxEpisodes,
		g.currentScore, avgScore, g.bestScore,
		g.agent.Epsilon(), g.agent.ReplayBufferSize(),
		g.lastMapSize, occupancy, len(g.snake.Obstacles()), g.speedMultiplier,
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
	// ✅ ИСПРАВЛЕНО: центрирование меню
	centerX := g.screenWidth / 2
	startY := 120

	title := "AI SNAKE GAME - Deep Q-Learning"
	titleWidth := len(title) * 6
	ebitenutil.DebugPrintAt(screen, title, centerX-titleWidth/2, startY)

	subtitle := "Self-learning snake powered by neural networks"
	subtitleWidth := len(subtitle) * 6
	ebitenutil.DebugPrintAt(screen, subtitle, centerX-subtitleWidth/2, startY+40)

	// Разделитель
	separator := "================================================"
	sepWidth := len(separator) * 6
	ebitenutil.DebugPrintAt(screen, separator, centerX-sepWidth/2, startY+90)

	// Кнопки (выровнены по левому краю от центра)
	buttonX := centerX - 150
	ebitenutil.DebugPrintAt(screen, "[SPACE] - Start Training", buttonX, startY+130)
	ebitenutil.DebugPrintAt(screen, "[P]     - Play with Trained AI", buttonX, startY+160)
	ebitenutil.DebugPrintAt(screen, "[Q]     - Quit", buttonX, startY+190)

	ebitenutil.DebugPrintAt(screen, separator, centerX-sepWidth/2, startY+230)

	// Статистика
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best Score: %d", g.bestScore), buttonX, startY+270)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Episodes Trained: %d", g.episode-1), buttonX, startY+300)

	// Особенности
	ebitenutil.DebugPrintAt(screen, "Features:", buttonX, startY+350)
	ebitenutil.DebugPrintAt(screen, "  • Random yellow obstacles", buttonX, startY+380)
	ebitenutil.DebugPrintAt(screen, "  • Auto map expansion at 90% occupancy", buttonX, startY+410)
	ebitenutil.DebugPrintAt(screen, "  • Wrap-around boundaries", buttonX, startY+440)
	ebitenutil.DebugPrintAt(screen, "  • Deep Q-Learning with Experience Replay", buttonX, startY+470)

	// Управление внизу
	info := "Controls: [1] 1x [2] 5x [3] 10x [4] 50x speed | [ESC] Menu"
	infoWidth := len(info) * 6
	ebitenutil.DebugPrintAt(screen, info, centerX-infoWidth/2, g.screenHeight-30)
}

func (g *Game) drawTraining(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	// ✅ ИСПРАВЛЕНО: компактный фон статистики
	if g.statsText != "" {
		textWidth := float32(len(g.statsText) * 4)
		if textWidth < 600 {
			textWidth = 600
		}
		vector.FillRect(screen, 10, 10, textWidth, 45, ui.TextBg, false)
		ebitenutil.DebugPrintAt(screen, g.statsText, 15, 15)
	}

	g.renderer.DrawProgressBar(screen, float64(g.episode)/float64(g.maxEpisodes), g.episode, g.maxEpisodes)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	// ✅ ИСПРАВЛЕНО: компактная статистика
	scoreText := fmt.Sprintf("Score: %d | Best: %d", g.currentScore, g.bestScore)
	textWidth := float32(len(scoreText) * 6)
	vector.FillRect(screen, 10, 10, textWidth+20, 35, ui.TextBg, false)
	ebitenutil.DebugPrintAt(screen, scoreText, 15, 15)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	// Затемнение
	vector.FillRect(screen, 0, 0, float32(g.screenWidth), float32(g.screenHeight), ui.TextBg, false)

	// Окно Game Over
	boxW, boxH := float32(360), float32(240)
	boxX := float32(g.screenWidth)/2 - boxW/2
	boxY := float32(g.screenHeight)/2 - boxH/2

	vector.FillRect(screen, boxX, boxY, boxW, boxH, ui.Background, false)
	vector.StrokeRect(screen, boxX, boxY, boxW, boxH, 3, ui.SnakeHead, false)

	centerX := g.screenWidth / 2
	centerY := g.screenHeight / 2

	ebitenutil.DebugPrintAt(screen, "GAME OVER", centerX-50, centerY-80)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.currentScore), centerX-40, centerY-40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best: %d", g.bestScore), centerX-35, centerY-10)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Play Again", centerX-75, centerY+30)
	ebitenutil.DebugPrintAt(screen, "[ESC] Main Menu", centerX-65, centerY+60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
