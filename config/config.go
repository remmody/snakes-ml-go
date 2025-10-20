package config

import "image/color"

// ================================
// WINDOW SETTINGS
// ================================
const (
	WindowWidth  = 1280
	WindowHeight = 720
	WindowTitle  = "AI Snake Game - Deep Q-Learning"
)

// ================================
// GAME SETTINGS
// ================================
const (
	InitialFieldWidth  = 20
	InitialFieldHeight = 15

	WrapAroundEnabled  = true
	DynamicSizeEnabled = true
	MaxFieldExpansion  = 3

	InitialObstaclesMin = 2
	InitialObstaclesMax = 3
	ObstacleAddInterval = 10
	ObstacleSafeRadius  = 3

	ExpansionThreshold = 0.9
	ExpansionIncrement = 2

	PlayingSpeed = 5
)

// ================================
// AI/TRAINING SETTINGS
// ================================
const (
	MaxEpisodes    = 10000
	EpisodesPerGen = 100
	WindowSize     = 100

	ModelBestName      = "snake_ai_model_best.json"
	ModelFinalName     = "snake_ai_model_final.json"
	ModelGenPrefix     = "snake_ai_model_gen"
	SaveCheckpointFreq = 100

	StateSize    = 28  // ✅ Увеличено с 24 до 28 (добавлено больше информации о пространстве)
	ActionSize   = 4
	HiddenLayer1 = 256
	HiddenLayer2 = 256
	HiddenLayer3 = 128 // ✅ НОВОЕ: добавлен третий слой

	LearningRate  = 0.0003 // ✅ Еще меньше для стабильности
	BufferSize    = 1000000
	EpsilonStart  = 1.0
	EpsilonMin    = 0.01
	EpsilonDecay  = 0.9997 // ✅ Еще медленнее
	Gamma         = 0.99
	BatchSize     = 128
	UpdateFreq    = 200    // ✅ Реже обновляем target network
	MinBufferSize = 256    // ✅ Больше минимум

	Speed1x  = 1.0
	Speed5x  = 5.0
	Speed10x = 10.0
	Speed50x = 50.0
)

// ================================
// REWARD SYSTEM (улучшена)
// ================================
const (
	RewardStep          = -0.01  // Штраф за шаг
	RewardFood          = 10.0   // Награда за еду
	RewardDeath         = -20.0  // Штраф за смерть
	RewardMoveToFood    = 0.5    // ✅ Увеличена награда за движение к еде
	RewardMoveFromFood  = -0.5   // Штраф за отдаление от еды
	RewardDanger        = -2.0   // ✅ Увеличен штраф за опасные позиции
	RewardSafeMove      = 0.2    // Награда за безопасное движение
	RewardNearBody      = -1.0   // ✅ Увеличен штраф за движение рядом с телом
	RewardCycle         = -0.8   // ✅ НОВОЕ: штраф за циклические движения
	RewardFreeSpace     = 0.3    // ✅ НОВОЕ: награда за движение в открытое пространство
	RewardTrap          = -3.0   // ✅ НОВОЕ: большой штраф за попадание в ловушку
)

// ================================
// UI/RENDERING SETTINGS
// ================================
const (
	GridStartX   = 30
	GridStartY   = 150
	CellSizeMin  = 8
	CellSizeMax  = 30
	CellSizeInit = 25
	GridPadding  = 80

	MenuStartY     = 120
	StatsBoxWidth  = 720
	StatsBoxHeight = 65
	StatsUpdateMs  = 1000

	ProgressBarHeight = 30
	ProgressBarMargin = 50

	GameOverBoxWidth  = 360
	GameOverBoxHeight = 240

	ColorSnakeBodyMin uint8 = 100
)

// ================================
// COLOR SCHEME
// ================================
var (
	ColorBackground = color.RGBA{20, 20, 30, 255}
	ColorTextBg     = color.RGBA{0, 0, 0, 180}

	ColorGrid = color.RGBA{50, 50, 60, 255}

	ColorObstacle       = color.RGBA{255, 200, 0, 255}
	ColorObstacleBorder = color.RGBA{180, 140, 0, 255}

	ColorFood       = color.RGBA{255, 80, 80, 255}
	ColorFoodBorder = color.RGBA{200, 50, 50, 255}

	ColorSnakeHead       = color.RGBA{100, 255, 100, 255}
	ColorSnakeHeadBorder = color.RGBA{50, 200, 50, 255}

	ColorProgressBg     = color.RGBA{40, 40, 50, 255}
	ColorProgressBorder = color.RGBA{100, 100, 120, 255}
)

// ================================
// MENU TEXT
// ================================
const (
	MenuTitle     = "AI SNAKE GAME - Deep Q-Learning"
	MenuSubtitle  = "Self-learning snake with A* pathfinding + DQN"
	MenuSeparator = "================================================"

	MenuBtnTraining = "[SPACE] - Start Training"
	MenuBtnPlay     = "[P]     - Play with Trained AI"
	MenuBtnQuit     = "[Q]     - Quit"

	MenuFeatures = "Features:"
	MenuFeature1 = "  • A* pathfinding + Deep Q-Learning"
	MenuFeature2 = "  • 28-parameter state space"
	MenuFeature3 = "  • Wrap-around collision detection"
	MenuFeature4 = "  • Cycle prevention algorithm"
	MenuFeature5 = "  • 100 episodes = 1 generation"

	MenuControls = "Controls: [1] 1x [2] 5x [3] 10x [4] 50x speed | [ESC] Menu"
)

// ================================
// HELPER FUNCTIONS
// ================================

func GetNeuralLayers() []int {
	return []int{StateSize, HiddenLayer1, HiddenLayer2, HiddenLayer3, ActionSize}
}

func GetInitialObstacles() int {
	return InitialObstaclesMin
}
