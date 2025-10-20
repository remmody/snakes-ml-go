package game

import (
	"fmt"
	"time"

	"snakes-ml/config"
	"snakes-ml/internal/ai"
	"snakes-ml/internal/snake"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Game represents main game structure
type Game struct {
	screenWidth     int
	screenHeight    int
	state           State
	snake           *snake.Snake
	agent           *ai.Agent
	renderer        *Renderer
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

// NewGame creates new game instance using config
func NewGame(screenWidth, screenHeight int) *Game {
	g := &Game{
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		state:           StateMenu,
		maxEpisodes:     config.MaxEpisodes,
		windowSize:      config.WindowSize,
		recentScores:    make([]int, 0, config.WindowSize),
		trainingMode:    true,
		autoRestart:     true,
		speedMultiplier: config.Speed1x,
		lastUpdateTime:  time.Now(),
	}

	g.renderer = NewRenderer(screenWidth, screenHeight)

	aiConfig := ai.DefaultConfig()
	g.agent = ai.NewAgent(config.StateSize, config.ActionSize, aiConfig)

	// Load existing model if available
	if err := g.agent.LoadModel(config.ModelBestName); err == nil {
		fmt.Println("✅ Loaded existing model")
	} else {
		fmt.Println("🆕 Created new model")
	}

	return g
}

// Update updates game state
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
	// Speed control from config
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.speedMultiplier = config.Speed1x
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.speedMultiplier = config.Speed5x
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.speedMultiplier = config.Speed10x
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.speedMultiplier = config.Speed50x
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
		return nil
	}

	// Training loop
	for i := 0; i < int(g.speedMultiplier); i++ {
		if g.snake == nil {
			g.startNewEpisode()
		}

		state := g.snake.GetState()
		action := g.agent.SelectAction(state)
		reward, done := g.snake.Step(action)
		nextState := g.snake.GetState()

		g.agent.Remember(state, action, reward, nextState, done)

		if g.agent.ReplayBufferSize() >= config.MinBufferSize {
			g.agent.Train()
		}

		g.currentScore = g.snake.Score()

		if done {
			g.handleEpisodeEnd()

			if g.agent.EpisodeCount() >= g.maxEpisodes {
				g.agent.SaveModel(config.ModelFinalName)
				fmt.Println("\n✅ Training completed!")
				g.state = StateMenu
				return nil
			}
		}
	}

	// Update stats display
	if time.Since(g.lastUpdateTime) > time.Millisecond*time.Duration(config.StatsUpdateMs) {
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

	if g.frameCount%config.PlayingSpeed == 0 {
		if g.snake == nil {
			g.startNewEpisode()
		}

		// AI plays without exploration
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
	g.startNewEpisode()
}

func (g *Game) startPlaying() {
	g.state = StatePlaying
	g.trainingMode = false
	g.startNewEpisode()
}

func (g *Game) startNewEpisode() {
	g.snake = snake.NewSnake(
		config.InitialFieldWidth,
		config.InitialFieldHeight,
		config.WrapAroundEnabled,
		config.DynamicSizeEnabled,
	)
	g.currentScore = 0
	g.lastMapSize = fmt.Sprintf("%dx%d", g.snake.Width(), g.snake.Height())
}

func (g *Game) handleEpisodeEnd() {
	score := g.snake.Score()
	g.recentScores = append(g.recentScores, score)
	if len(g.recentScores) > g.windowSize {
		g.recentScores = g.recentScores[1:]
	}

	// Mark episode end for agent
	g.agent.EndEpisode()

	// Save best model
	if score > g.bestScore {
		g.bestScore = score
		g.agent.SaveModel(config.ModelBestName)
		fmt.Printf("🏆 New record: %d (episode %d, generation %d)\n",
			score, g.agent.EpisodeCount(), g.agent.Generation())
	}

	// Save generation checkpoints
	if g.agent.EpisodeCount()%config.SaveCheckpointFreq == 0 {
		filename := fmt.Sprintf("%s%d.json", config.ModelGenPrefix, g.agent.Generation())
		g.agent.SaveModel(filename)
		fmt.Printf("💾 Generation %d completed. Checkpoint saved: %s\n",
			g.agent.Generation(), filename)
	}

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
		"Gen: %d (%d/%d) | Ep: %d/%d | Score: %d | Avg: %.1f | Best: %d\n"+
			"ε: %.3f | Buf: %d | Map: %s | Occ: %.0f%% | Obs: %d | x%.0f",
		g.agent.Generation(),
		g.agent.GenerationProgress(),
		config.EpisodesPerGen,
		g.agent.EpisodeCount(),
		g.maxEpisodes,
		g.currentScore,
		avgScore,
		g.bestScore,
		g.agent.Epsilon(),
		g.agent.ReplayBufferSize(),
		g.lastMapSize,
		occupancy,
		len(g.snake.Obstacles()),
		g.speedMultiplier,
	)
}

// Draw renders game
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(config.ColorBackground)

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
	centerX := g.screenWidth / 2
	startY := config.MenuStartY

	// Title
	title := config.MenuTitle
	titleWidth := len(title) * 6
	ebitenutil.DebugPrintAt(screen, title, centerX-titleWidth/2, startY)

	// Subtitle
	subtitle := config.MenuSubtitle
	subtitleWidth := len(subtitle) * 6
	ebitenutil.DebugPrintAt(screen, subtitle, centerX-subtitleWidth/2, startY+40)

	// Separator
	separator := config.MenuSeparator
	sepWidth := len(separator) * 6
	ebitenutil.DebugPrintAt(screen, separator, centerX-sepWidth/2, startY+90)

	// Buttons
	buttonX := centerX - 150
	ebitenutil.DebugPrintAt(screen, config.MenuBtnTraining, buttonX, startY+130)
	ebitenutil.DebugPrintAt(screen, config.MenuBtnPlay, buttonX, startY+160)
	ebitenutil.DebugPrintAt(screen, config.MenuBtnQuit, buttonX, startY+190)

	ebitenutil.DebugPrintAt(screen, separator, centerX-sepWidth/2, startY+230)

	// Statistics
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best Score: %d", g.bestScore), buttonX, startY+270)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Episodes Trained: %d", g.agent.EpisodeCount()), buttonX, startY+300)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Generations: %d", g.agent.Generation()), buttonX, startY+330)

	// Features
	ebitenutil.DebugPrintAt(screen, config.MenuFeatures, buttonX, startY+380)
	ebitenutil.DebugPrintAt(screen, config.MenuFeature1, buttonX, startY+410)
	ebitenutil.DebugPrintAt(screen, config.MenuFeature2, buttonX, startY+440)
	ebitenutil.DebugPrintAt(screen, config.MenuFeature3, buttonX, startY+470)
	ebitenutil.DebugPrintAt(screen, config.MenuFeature4, buttonX, startY+500)
	ebitenutil.DebugPrintAt(screen, config.MenuFeature5, buttonX, startY+530)

	// Controls
	info := config.MenuControls
	infoWidth := len(info) * 6
	ebitenutil.DebugPrintAt(screen, info, centerX-infoWidth/2, g.screenHeight-30)
}

func (g *Game) drawTraining(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	// Stats box
	if g.statsText != "" {
		vector.FillRect(screen, 10, 10, float32(config.StatsBoxWidth), float32(config.StatsBoxHeight), config.ColorTextBg, false)
		ebitenutil.DebugPrintAt(screen, g.statsText, 15, 15)
	}

	// Progress bar
	totalGenerations := g.maxEpisodes / config.EpisodesPerGen
	currentGen := g.agent.Generation()
	progress := float64(currentGen) / float64(totalGenerations)

	g.renderer.DrawProgressBar(screen, progress, currentGen, totalGenerations)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	scoreText := fmt.Sprintf("Score: %d | Best: %d", g.currentScore, g.bestScore)
	textWidth := float32(len(scoreText) * 6)
	vector.FillRect(screen, 10, 10, textWidth+20, 35, config.ColorTextBg, false)
	ebitenutil.DebugPrintAt(screen, scoreText, 15, 15)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	if g.snake != nil {
		g.renderer.DrawSnake(screen, g.snake)
	}

	// Dimming overlay
	vector.FillRect(screen, 0, 0, float32(g.screenWidth), float32(g.screenHeight), config.ColorTextBg, false)

	// Game over box
	boxW := float32(config.GameOverBoxWidth)
	boxH := float32(config.GameOverBoxHeight)
	boxX := float32(g.screenWidth)/2 - boxW/2
	boxY := float32(g.screenHeight)/2 - boxH/2

	vector.FillRect(screen, boxX, boxY, boxW, boxH, config.ColorBackground, false)
	vector.StrokeRect(screen, boxX, boxY, boxW, boxH, 3, config.ColorSnakeHead, false)

	centerX := g.screenWidth / 2
	centerY := g.screenHeight / 2

	ebitenutil.DebugPrintAt(screen, "GAME OVER", centerX-50, centerY-80)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.currentScore), centerX-40, centerY-40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best: %d", g.bestScore), centerX-35, centerY-10)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Play Again", centerX-75, centerY+30)
	ebitenutil.DebugPrintAt(screen, "[ESC] Main Menu", centerX-65, centerY+60)
}

// Layout returns screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
