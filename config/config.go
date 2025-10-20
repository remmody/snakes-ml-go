package config

import "image/color"

// ================================
// WINDOW SETTINGS
// ================================
const (
	// Window resolution and title
	WindowWidth  = 1280
	WindowHeight = 720
	WindowTitle  = "AI Snake Game - Deep Q-Learning"
)

// ================================
// GAME SETTINGS
// ================================
const (
	// Initial game field size
	InitialFieldWidth  = 20 // Starting width in cells
	InitialFieldHeight = 15 // Starting height in cells

	// Game mechanics
	WrapAroundEnabled  = true // Snake can pass through walls
	DynamicSizeEnabled = true // Field expands when snake grows
	MaxFieldExpansion  = 3    // Maximum expansion multiplier (e.g., 3x means 60x45)

	// Obstacle settings
	InitialObstaclesMin = 2  // Minimum starting obstacles
	InitialObstaclesMax = 3  // Maximum starting obstacles (will be min + random(max-min))
	ObstacleAddInterval = 10 // Add 1 obstacle every N points
	ObstacleSafeRadius  = 3  // Minimum distance from snake when spawning

	// Field expansion
	ExpansionThreshold = 0.9 // Expand when 90% occupied
	ExpansionIncrement = 2   // Add 2 cells per dimension when expanding

	// Game speed settings
	PlayingSpeed = 5 // Update every N frames in playing mode
)

// ================================
// AI/TRAINING SETTINGS
// ================================
const (
	// Training parameters
	MaxEpisodes    = 10000 // Total episodes to train
	EpisodesPerGen = 100   // Episodes per generation
	WindowSize     = 100   // Window for averaging scores

	// Model saving
	ModelBestName      = "snake_ai_model_best.json"  // Best model filename
	ModelFinalName     = "snake_ai_model_final.json" // Final model filename
	ModelGenPrefix     = "snake_ai_model_gen"        // Generation checkpoint prefix
	SaveCheckpointFreq = 100                         // Save checkpoint every N episodes

	// Neural network architecture
	StateSize    = 14  // Input size (danger, food direction, etc.)
	ActionSize   = 4   // Output size (up, right, down, left)
	HiddenLayer1 = 128 // First hidden layer neurons
	HiddenLayer2 = 128 // Second hidden layer neurons

	// DQN hyperparameters
	LearningRate  = 0.001    // Learning rate for neural network
	BufferSize    = 1000000  // Experience replay buffer size
	EpsilonStart  = 1.0      // Starting exploration rate
	EpsilonMin    = 0.01     // Minimum exploration rate
	EpsilonDecay  = 0.995    // Exploration decay rate
	Gamma         = 0.95     // Discount factor for future rewards
	BatchSize     = 64       // Batch size for training
	UpdateFreq    = 100      // Update target network every N steps
	MinBufferSize = 64       // Minimum buffer size before training

	// Speed multipliers for training
	Speed1x  = 1.0
	Speed5x  = 5.0
	Speed10x = 10.0
	Speed50x = 50.0
)

// ================================
// REWARD SYSTEM
// ================================
const (
	RewardStep       = -0.01 // Penalty for each step (encourages efficiency)
	RewardFood       = 10.0  // Reward for eating food
	RewardDeath      = -10.0 // Penalty for dying
	RewardMoveToFood = 0.1   // Small reward for moving toward food
)

// ================================
// UI/RENDERING SETTINGS
// ================================
const (
	// Grid rendering
	GridStartX   = 30 // Starting X position for game grid
	GridStartY   = 150 // Starting Y position for game grid
	CellSizeMin  = 8  // Minimum cell size in pixels
	CellSizeMax  = 30 // Maximum cell size in pixels
	CellSizeInit = 25 // Initial cell size
	GridPadding  = 80 // Bottom padding for grid

	// UI text positioning
	MenuStartY     = 120  // Menu start Y position
	StatsBoxWidth  = 680  // Stats box width
	StatsBoxHeight = 45   // Stats box height
	StatsUpdateMs  = 1000 // Update stats every N milliseconds

	// Progress bar
	ProgressBarHeight = 30
	ProgressBarMargin = 50

	// Game over box
	GameOverBoxWidth  = 360
	GameOverBoxHeight = 240

	// Snake body color intensity (✅ ИСПРАВЛЕНО: uint8 вместо int)
	ColorSnakeBodyMin uint8 = 100 // Minimum body color intensity
)

// ================================
// COLOR SCHEME
// ================================
var (
	// Background colors
	ColorBackground = color.RGBA{20, 20, 30, 255}
	ColorTextBg     = color.RGBA{0, 0, 0, 180}

	// Grid colors
	ColorGrid = color.RGBA{50, 50, 60, 255}

	// Obstacle colors (yellow)
	ColorObstacle       = color.RGBA{255, 200, 0, 255}
	ColorObstacleBorder = color.RGBA{180, 140, 0, 255}

	// Food colors (red)
	ColorFood       = color.RGBA{255, 80, 80, 255}
	ColorFoodBorder = color.RGBA{200, 50, 50, 255}

	// Snake colors (green)
	ColorSnakeHead       = color.RGBA{100, 255, 100, 255}
	ColorSnakeHeadBorder = color.RGBA{50, 200, 50, 255}

	// Progress bar colors
	ColorProgressBg     = color.RGBA{40, 40, 50, 255}
	ColorProgressBorder = color.RGBA{100, 100, 120, 255}
)

// ================================
// MENU TEXT
// ================================
const (
	MenuTitle    = "AI SNAKE GAME - Deep Q-Learning"
	MenuSubtitle = "Self-learning snake powered by neural networks"
	MenuSeparator = "================================================"

	MenuBtnTraining = "[SPACE] - Start Training"
	MenuBtnPlay     = "[P]     - Play with Trained AI"
	MenuBtnQuit     = "[Q]     - Quit"

	MenuFeatures = "Features:"
	MenuFeature1 = "  • Random yellow obstacles"
	MenuFeature2 = "  • Auto map expansion at 90% occupancy"
	MenuFeature3 = "  • Wrap-around boundaries"
	MenuFeature4 = "  • Deep Q-Learning with Experience Replay"
	MenuFeature5 = "  • 100 episodes = 1 generation"

	MenuControls = "Controls: [1] 1x [2] 5x [3] 10x [4] 50x speed | [ESC] Menu"
)

// ================================
// HELPER FUNCTIONS
// ================================

// GetNeuralLayers returns the neural network architecture
func GetNeuralLayers() []int {
	return []int{StateSize, HiddenLayer1, HiddenLayer2, ActionSize}
}

// GetInitialObstacles returns random number of initial obstacles
func GetInitialObstacles() int {
	return InitialObstaclesMin
}
