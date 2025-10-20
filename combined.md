# Combined Source Code

Generated from: `.`

---

## cmd\game\main.go

<!-- source: cmd\game\main.go -->

```go
package main

import (
	"log"
	"snakes-ml/config"
	"snakes-ml/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Window configuration from central config
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create and run game
	g := game.NewGame(config.WindowWidth, config.WindowHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
```

---

## config\config.go

<!-- source: config\config.go -->

```go
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
	WrapAroundEnabled  = true  // Snake can pass through walls
	DynamicSizeEnabled = true  // Field expands when snake grows
	MaxFieldExpansion  = 3     // Maximum expansion multiplier (e.g., 3x means 60x45)
	
	// Obstacle settings
	InitialObstaclesMin = 2  // Minimum starting obstacles
	InitialObstaclesMax = 3  // Maximum starting obstacles (will be min + random(max-min))
	ObstacleAddInterval = 10 // Add 1 obstacle every N points
	ObstacleSafeRadius  = 3  // Minimum distance from snake when spawning
	
	// Field expansion
	ExpansionThreshold = 0.9  // Expand when 90% occupied
	ExpansionIncrement = 2    // Add 2 cells per dimension when expanding
	
	// Game speed settings
	PlayingSpeed = 5 // Update every N frames in playing mode
)

// ================================
// AI/TRAINING SETTINGS
// ================================
const (
	// Training parameters
	MaxEpisodes     = 10000 // Total episodes to train
	EpisodesPerGen  = 100   // Episodes per generation
	WindowSize      = 100   // Window for averaging scores
	
	// Model saving
	ModelBestName     = "snake_ai_model_best.json"     // Best model filename
	ModelFinalName    = "snake_ai_model_final.json"    // Final model filename
	ModelGenPrefix    = "snake_ai_model_gen"           // Generation checkpoint prefix
	SaveCheckpointFreq = 100                           // Save checkpoint every N episodes
	
	// Neural network architecture
	StateSize  = 14  // Input size (danger, food direction, etc.)
	ActionSize = 4   // Output size (up, right, down, left)
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
	RewardStep        = -0.01 // Penalty for each step (encourages efficiency)
	RewardFood        = 10.0  // Reward for eating food
	RewardDeath       = -10.0 // Penalty for dying
	RewardMoveToFood  = 0.1   // Small reward for moving toward food
)

// ================================
// UI/RENDERING SETTINGS
// ================================
const (
	// Grid rendering
	GridStartX   = 30  // Starting X position for game grid
	GridStartY   = 150 // Starting Y position for game grid
	CellSizeMin  = 8   // Minimum cell size in pixels
	CellSizeMax  = 30  // Maximum cell size in pixels
	CellSizeInit = 25  // Initial cell size
	GridPadding  = 80  // Bottom padding for grid
	
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
	ColorSnakeBodyMin    = 100 // Minimum body color intensity
	
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
```

---

## go.mod

<!-- source: go.mod -->

```text
module snakes-ml

go 1.25.1

require github.com/hajimehoshi/ebiten/v2 v2.9.2

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
```

---

## go.sum

<!-- source: go.sum -->

```text
github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 h1:+kz5iTT3L7uU+VhlMfTb8hHcxLO3TlaELlX8wa4XjA0=
github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1/go.mod h1:lKJoeixeJwnFmYsBny4vvCJGVFc3aYDalhuDsfZzWHI=
github.com/ebitengine/hideconsole v1.0.0 h1:5J4U0kXF+pv/DhiXt5/lTz0eO5ogJ1iXb8Yj1yReDqE=
github.com/ebitengine/hideconsole v1.0.0/go.mod h1:hTTBTvVYWKBuxPr7peweneWdkUwEuHuB3C1R/ielR1A=
github.com/ebitengine/purego v0.9.0 h1:mh0zpKBIXDceC63hpvPuGLiJ8ZAa3DfrFTudmfi8A4k=
github.com/ebitengine/purego v0.9.0/go.mod h1:iIjxzd6CiRiOG0UyXP+V1+jWqUXVjPKLAI0mRfJZTmQ=
github.com/hajimehoshi/ebiten/v2 v2.9.2 h1:fV9Wh8dL4gSV62s/oygCIckEJ2MPpJRaUiwj6GR4Uos=
github.com/hajimehoshi/ebiten/v2 v2.9.2/go.mod h1:DAt4tnkYYpCvu3x9i1X/nK/vOruNXIlYq/tBXxnhrXM=
github.com/jezek/xgb v1.1.1 h1:bE/r8ZZtSv7l9gk6nU0mYx51aXrvnyb44892TwSaqS4=
github.com/jezek/xgb v1.1.1/go.mod h1:nrhwO0FX/enq75I7Y7G8iN1ubpSGZEiA3v9e9GyRFlk=
golang.org/x/image v0.31.0 h1:mLChjE2MV6g1S7oqbXC0/UcKijjm5fnJLUYKIYrLESA=
golang.org/x/image v0.31.0/go.mod h1:R9ec5Lcp96v9FTF+ajwaH3uGxPH4fKfHHAVbUILxghA=
golang.org/x/sync v0.17.0 h1:l60nONMj9l5drqw6jlhIELNv9I0A4OFgRsG9k2oT9Ug=
golang.org/x/sync v0.17.0/go.mod h1:9KTHXmSnoGruLpwFjVSX0lNNA75CykiMECbovNTZqGI=
golang.org/x/sys v0.36.0 h1:KVRy2GtZBrk1cBYA7MKu5bEZFxQk4NIDV6RLVcC8o0k=
golang.org/x/sys v0.36.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
```

---

## internal\ai\agent.go

<!-- source: internal\ai\agent.go -->

```go
package ai

import (
	"math/rand/v2"
	"snakes-ml/config"
)

// Agent represents DQN agent with generation system
type Agent struct {
	qNetwork          *Network
	targetNetwork     *Network
	replayBuffer      *ReplayBuffer
	epsilon           float64
	epsilonMin        float64
	epsilonDecay      float64
	gamma             float64
	batchSize         int
	updateFreq        int
	stepCount         int
	episodeCount      int
	generationSize    int
	currentGeneration int
	totalReward       float64
	episodeRewards    []float64
}

// NewAgent creates new DQN agent using configuration
func NewAgent(stateSize, actionSize int, cfg Config) *Agent {
	layers := config.GetNeuralLayers()

	return &Agent{
		qNetwork:          NewNetwork(layers, cfg.LearningRate),
		targetNetwork:     NewNetwork(layers, cfg.LearningRate),
		replayBuffer:      NewReplayBuffer(cfg.BufferSize),
		epsilon:           cfg.EpsilonStart,
		epsilonMin:        cfg.EpsilonMin,
		epsilonDecay:      cfg.EpsilonDecay,
		gamma:             cfg.Gamma,
		batchSize:         cfg.BatchSize,
		updateFreq:        cfg.UpdateFreq,
		stepCount:         0,
		episodeCount:      0,
		generationSize:    config.EpisodesPerGen,
		currentGeneration: 1,
		episodeRewards:    make([]float64, 0, 100),
	}
}

// SelectAction chooses action using epsilon-greedy strategy
func (a *Agent) SelectAction(state []float64) int {
	if rand.Float64() < a.epsilon {
		return rand.IntN(config.ActionSize)
	}

	qValues := a.qNetwork.Forward(state)
	return argmax(qValues)
}

// argmax returns index of maximum value
func argmax(values []float64) int {
	if len(values) == 0 {
		return 0
	}

	maxIdx := 0
	maxVal := values[0]

	for i, v := range values {
		if v > maxVal {
			maxVal = v
			maxIdx = i
		}
	}

	return maxIdx
}

// Remember stores experience in replay buffer
func (a *Agent) Remember(state []float64, action int, reward float64, nextState []float64, done bool) {
	a.replayBuffer.Add(Experience{
		State:     state,
		Action:    action,
		Reward:    reward,
		NextState: nextState,
		Done:      done,
	})

	// Accumulate episode reward
	a.totalReward += reward
}

// Train performs one training step using experience replay
func (a *Agent) Train() float64 {
	if a.replayBuffer.Size() < a.batchSize {
		return 0
	}

	batch := a.replayBuffer.Sample(a.batchSize)
	totalLoss := 0.0

	for _, exp := range batch {
		target := a.qNetwork.Forward(exp.State)

		if exp.Done {
			target[exp.Action] = exp.Reward
		} else {
			// Double DQN: use q-network to select action, target-network to evaluate
			nextQValues := a.targetNetwork.Forward(exp.NextState)
			bestAction := argmax(a.qNetwork.Forward(exp.NextState))
			maxQ := nextQValues[bestAction]
			target[exp.Action] = exp.Reward + a.gamma*maxQ
		}

		loss := a.qNetwork.BackwardAndUpdate(exp.State, target)
		totalLoss += loss
	}

	a.stepCount++
	if a.stepCount%a.updateFreq == 0 {
		a.UpdateTargetNetwork()
	}

	// Epsilon decay during training
	if a.epsilon > a.epsilonMin {
		a.epsilon *= a.epsilonDecay
	}

	return totalLoss / float64(len(batch))
}

// EndEpisode marks end of episode and updates generation counter
func (a *Agent) EndEpisode() {
	a.episodeCount++
	a.episodeRewards = append(a.episodeRewards, a.totalReward)

	// Keep only last 100 episode rewards
	if len(a.episodeRewards) > 100 {
		a.episodeRewards = a.episodeRewards[1:]
	}

	// Reset accumulated reward
	a.totalReward = 0

	// Check generation change
	if a.episodeCount%a.generationSize == 0 {
		a.currentGeneration++
	}
}

// GetAverageReward returns average reward over last N episodes
func (a *Agent) GetAverageReward(window int) float64 {
	if len(a.episodeRewards) == 0 {
		return 0
	}

	start := 0
	if len(a.episodeRewards) > window {
		start = len(a.episodeRewards) - window
	}

	sum := 0.0
	count := 0
	for i := start; i < len(a.episodeRewards); i++ {
		sum += a.episodeRewards[i]
		count++
	}

	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// UpdateTargetNetwork copies weights from q-network to target-network
func (a *Agent) UpdateTargetNetwork() {
	a.targetNetwork = a.qNetwork.Clone()
}

// SaveModel saves neural network to file
func (a *Agent) SaveModel(filename string) error {
	return a.qNetwork.SaveToFile(filename)
}

// LoadModel loads neural network from file
func (a *Agent) LoadModel(filename string) error {
	err := a.qNetwork.LoadFromFile(filename)
	if err == nil {
		a.targetNetwork = a.qNetwork.Clone()
	}
	return err
}

// Getters
func (a *Agent) Epsilon() float64           { return a.epsilon }
func (a *Agent) SetEpsilon(epsilon float64) { a.epsilon = epsilon }
func (a *Agent) ReplayBufferSize() int      { return a.replayBuffer.Size() }
func (a *Agent) StepCount() int             { return a.stepCount }
func (a *Agent) EpisodeCount() int          { return a.episodeCount }
func (a *Agent) Generation() int            { return a.currentGeneration }
func (a *Agent) GenerationProgress() int    { return a.episodeCount % a.generationSize }
```

---

## internal\ai\config.go

<!-- source: internal\ai\config.go -->

```go
package ai

import "snakes-ml/config"

// Config holds DQN agent configuration
type Config struct {
	LearningRate float64
	BufferSize   int
	EpsilonStart float64
	EpsilonMin   float64
	EpsilonDecay float64
	Gamma        float64
	BatchSize    int
	UpdateFreq   int
}

// DefaultConfig returns default DQN configuration from central config
func DefaultConfig() Config {
	return Config{
		LearningRate: config.LearningRate,
		BufferSize:   config.BufferSize,
		EpsilonStart: config.EpsilonStart,
		EpsilonMin:   config.EpsilonMin,
		EpsilonDecay: config.EpsilonDecay,
		Gamma:        config.Gamma,
		BatchSize:    config.BatchSize,
		UpdateFreq:   config.UpdateFreq,
	}
}
```

---

## internal\ai\network.go

<!-- source: internal\ai\network.go -->

```go
package ai

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"sync"
)

// Network представляет нейронную сеть с feed-forward архитектурой
type Network struct {
	layers       []int
	weights      [][][]float64
	biases       [][]float64
	learningRate float64
	mu           sync.RWMutex
}

// NewNetwork создает новую нейронную сеть
func NewNetwork(layers []int, learningRate float64) *Network {
	nn := &Network{
		layers:       layers,
		learningRate: learningRate,
	}

	nn.weights = make([][][]float64, len(layers)-1)
	nn.biases = make([][]float64, len(layers)-1)

	// Xavier/Glorot инициализация
	for i := 0; i < len(layers)-1; i++ {
		nn.weights[i] = make([][]float64, layers[i])
		nn.biases[i] = make([]float64, layers[i+1])

		limit := math.Sqrt(6.0 / float64(layers[i]+layers[i+1]))

		for j := 0; j < layers[i]; j++ {
			nn.weights[i][j] = make([]float64, layers[i+1])
			for k := 0; k < layers[i+1]; k++ {
				// ✅ ИСПРАВЛЕНИЕ: используем math/rand/v2
				nn.weights[i][j][k] = (rand.Float64()*2 - 1) * limit
			}
		}

		for k := 0; k < layers[i+1]; k++ {
			nn.biases[i][k] = (rand.Float64()*2 - 1) * limit
		}
	}

	return nn
}

// relu функция активации
func relu(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}

// reluDerivative производная ReLU
func reluDerivative(x float64) float64 {
	if x > 0 {
		return 1
	}
	return 0
}

// Forward прямой проход через сеть
func (nn *Network) Forward(input []float64) []float64 {
	nn.mu.RLock()
	defer nn.mu.RUnlock()

	current := input

	for i := 0; i < len(nn.weights); i++ {
		next := make([]float64, nn.layers[i+1])

		for j := 0; j < nn.layers[i+1]; j++ {
			sum := nn.biases[i][j]
			for k := 0; k < nn.layers[i]; k++ {
				sum += current[k] * nn.weights[i][k][j]
			}

			// ReLU для скрытых слоев, linear для выходного
			if i < len(nn.weights)-1 {
				next[j] = relu(sum)
			} else {
				next[j] = sum
			}
		}

		current = next
	}

	return current
}

// BackwardAndUpdate выполняет обратное распространение и обновляет веса
func (nn *Network) BackwardAndUpdate(input, target []float64) float64 {
	nn.mu.Lock()
	defer nn.mu.Unlock()

	// Forward pass с сохранением активаций
	activations := make([][]float64, len(nn.layers))
	activations[0] = input

	for i := 0; i < len(nn.weights); i++ {
		next := make([]float64, nn.layers[i+1])

		for j := 0; j < nn.layers[i+1]; j++ {
			sum := nn.biases[i][j]
			for k := 0; k < nn.layers[i]; k++ {
				sum += activations[i][k] * nn.weights[i][k][j]
			}

			if i < len(nn.weights)-1 {
				next[j] = relu(sum)
			} else {
				next[j] = sum
			}
		}

		activations[i+1] = next
	}

	// Вычисление ошибки выходного слоя
	deltas := make([][]float64, len(nn.layers)-1)
	lastIdx := len(activations) - 1
	deltas[lastIdx-1] = make([]float64, nn.layers[lastIdx])

	loss := 0.0
	for i := 0; i < nn.layers[lastIdx]; i++ {
		error := target[i] - activations[lastIdx][i]
		deltas[lastIdx-1][i] = error
		loss += error * error
	}

	// Backpropagation для скрытых слоев
	for i := len(nn.weights) - 2; i >= 0; i-- {
		deltas[i] = make([]float64, nn.layers[i+1])

		for j := 0; j < nn.layers[i+1]; j++ {
			sum := 0.0
			for k := 0; k < nn.layers[i+2]; k++ {
				sum += deltas[i+1][k] * nn.weights[i+1][j][k]
			}
			deltas[i][j] = sum * reluDerivative(activations[i+1][j])
		}
	}

	// Обновление весов и смещений
	for i := 0; i < len(nn.weights); i++ {
		for j := 0; j < nn.layers[i]; j++ {
			for k := 0; k < nn.layers[i+1]; k++ {
				nn.weights[i][j][k] += nn.learningRate * deltas[i][k] * activations[i][j]
			}
		}

		for k := 0; k < nn.layers[i+1]; k++ {
			nn.biases[i][k] += nn.learningRate * deltas[i][k]
		}
	}

	return loss / float64(len(target))
}

// Clone создает глубокую копию сети
func (nn *Network) Clone() *Network {
	nn.mu.RLock()
	defer nn.mu.RUnlock()

	clone := &Network{
		layers:       make([]int, len(nn.layers)),
		weights:      make([][][]float64, len(nn.weights)),
		biases:       make([][]float64, len(nn.biases)),
		learningRate: nn.learningRate,
	}

	copy(clone.layers, nn.layers)

	for i := range nn.weights {
		clone.weights[i] = make([][]float64, len(nn.weights[i]))
		clone.biases[i] = make([]float64, len(nn.biases[i]))
		copy(clone.biases[i], nn.biases[i])

		for j := range nn.weights[i] {
			clone.weights[i][j] = make([]float64, len(nn.weights[i][j]))
			copy(clone.weights[i][j], nn.weights[i][j])
		}
	}

	return clone
}

// SaveToFile сохраняет сеть в JSON файл
func (nn *Network) SaveToFile(filename string) error {
	nn.mu.RLock()
	defer nn.mu.RUnlock()

	data, err := json.Marshal(struct {
		Layers  []int         `json:"layers"`
		Weights [][][]float64 `json:"weights"`
		Biases  [][]float64   `json:"biases"`
	}{
		Layers:  nn.layers,
		Weights: nn.weights,
		Biases:  nn.biases,
	})

	if err != nil {
		return fmt.Errorf("marshal network: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadFromFile загружает сеть из JSON файла
func (nn *Network) LoadFromFile(filename string) error {
	nn.mu.Lock()
	defer nn.mu.Unlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var loaded struct {
		Layers  []int         `json:"layers"`
		Weights [][][]float64 `json:"weights"`
		Biases  [][]float64   `json:"biases"`
	}

	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("unmarshal network: %w", err)
	}

	nn.layers = loaded.Layers
	nn.weights = loaded.Weights
	nn.biases = loaded.Biases

	return nil
}

// Layers возвращает архитектуру сети
func (nn *Network) Layers() []int {
	nn.mu.RLock()
	defer nn.mu.RUnlock()
	return nn.layers
}
```

---

## internal\ai\replay.go

<!-- source: internal\ai\replay.go -->

```go
package ai

import (
	"math/rand/v2"
	"sync"
)

// Experience represents single training experience
type Experience struct {
	State     []float64
	Action    int
	Reward    float64
	NextState []float64
	Done      bool
}

// ReplayBuffer implements experience replay buffer
type ReplayBuffer struct {
	buffer   []Experience
	capacity int
	mu       sync.Mutex
}

// NewReplayBuffer creates new replay buffer
func NewReplayBuffer(capacity int) *ReplayBuffer {
	return &ReplayBuffer{
		buffer:   make([]Experience, 0, capacity),
		capacity: capacity,
	}
}

// Add adds experience to buffer (FIFO)
func (rb *ReplayBuffer) Add(exp Experience) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.buffer) < rb.capacity {
		rb.buffer = append(rb.buffer, exp)
	} else {
		// Remove oldest experience
		rb.buffer = append(rb.buffer[1:], exp)
	}
}

// Sample returns random batch of experiences
func (rb *ReplayBuffer) Sample(batchSize int) []Experience {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.buffer) < batchSize {
		batchSize = len(rb.buffer)
	}

	samples := make([]Experience, batchSize)
	indices := rand.Perm(len(rb.buffer))[:batchSize]

	for i, idx := range indices {
		samples[i] = rb.buffer[idx]
	}

	return samples
}

// Size returns current buffer size
func (rb *ReplayBuffer) Size() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer)
}

// Clear empties the buffer
func (rb *ReplayBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.buffer = make([]Experience, 0, rb.capacity)
}

// IsFull checks if buffer is full
func (rb *ReplayBuffer) IsFull() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer) >= rb.capacity
}
```

---

## internal\game\game.go

<!-- source: internal\game\game.go -->

```go
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
```

---

## internal\game\renderer.go

<!-- source: internal\game\renderer.go -->

```go
package game

import (
	"fmt"
	"image/color"

	"snakes-ml/config"
	"snakes-ml/internal/snake"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Renderer handles game rendering
type Renderer struct {
	screenWidth  int
	screenHeight int
}

// NewRenderer creates new renderer
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		screenWidth:  width,
		screenHeight: height,
	}
}

// DrawSnake renders snake and game field
func (r *Renderer) DrawSnake(screen *ebiten.Image, s *snake.Snake) {
	gridX := config.GridStartX
	gridY := config.GridStartY
	cellSize := config.CellSizeInit

	// Calculate adaptive cell size
	maxWidth := r.screenWidth - gridX*2
	maxHeight := r.screenHeight - gridY - config.GridPadding

	cellWidth := maxWidth / s.Width()
	cellHeight := maxHeight / s.Height()

	if cellWidth < cellHeight {
		cellSize = cellWidth
	} else {
		cellSize = cellHeight
	}

	if cellSize < config.CellSizeMin {
		cellSize = config.CellSizeMin
	}
	if cellSize > config.CellSizeMax {
		cellSize = config.CellSizeMax
	}

	totalWidth := s.Width() * cellSize
	totalHeight := s.Height() * cellSize
	gridX = (r.screenWidth - totalWidth) / 2

	// Draw grid
	for y := 0; y < s.Height(); y++ {
		for x := 0; x < s.Width(); x++ {
			posX := float32(gridX + x*cellSize)
			posY := float32(gridY + y*cellSize)
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 1, config.ColorGrid, false)
		}
	}

	// Draw obstacles (yellow)
	for _, obs := range s.Obstacles() {
		posX := float32(gridX + obs.X*cellSize)
		posY := float32(gridY + obs.Y*cellSize)
		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), config.ColorObstacle, false)
		vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorObstacleBorder, false)
	}

	// Draw food (red)
	food := s.Food()
	posX := float32(gridX + food.X*cellSize)
	posY := float32(gridY + food.Y*cellSize)
	vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), config.ColorFood, false)
	vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorFoodBorder, false)

	// Draw snake
	for i, segment := range s.Body() {
		posX := float32(gridX + segment.X*cellSize)
		posY := float32(gridY + segment.Y*cellSize)

		var col color.RGBA
		if i == 0 {
			col = config.ColorSnakeHead
		} else {
			intensity := uint8(200 - i*2)
			if intensity < config.ColorSnakeBodyMin {
				intensity = config.ColorSnakeBodyMin
			}
			col = color.RGBA{50, intensity, 50, 255}
		}

		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), col, false)

		if i == 0 {
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorSnakeHeadBorder, false)
		}
	}

	// Draw info
	occupancy := s.GetOccupancy() * 100
	infoText := fmt.Sprintf("Length: %d | Steps: %d | Map: %dx%d | Occupancy: %.1f%% | Obstacles: %d",
		s.Length(), s.Steps(), s.Width(), s.Height(), occupancy, len(s.Obstacles()))

	vector.FillRect(screen, float32(gridX-5), float32(gridY-40), float32(totalWidth+10), 35, config.ColorTextBg, false)
	ebitenutil.DebugPrintAt(screen, infoText, gridX, gridY-35)
}

// DrawProgressBar renders training progress bar
func (r *Renderer) DrawProgressBar(screen *ebiten.Image, progress float64, current, total int) {
	barX := float32(10)
	barY := float32(r.screenHeight - config.ProgressBarMargin)
	barWidth := float32(r.screenWidth - 20)
	barHeight := float32(config.ProgressBarHeight)

	vector.FillRect(screen, barX, barY, barWidth, barHeight, config.ColorProgressBg, false)
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2, config.ColorProgressBorder, false)

	fillWidth := barWidth * float32(progress)
	if fillWidth > 0 {
		r := uint8(100 - progress*50)
		g := uint8(200 - progress*50)
		b := uint8(100 + progress*100)
		vector.FillRect(screen, barX+2, barY+2, fillWidth-4, barHeight-4, color.RGBA{r, g, b, 255}, false)
	}

	progressText := fmt.Sprintf("Training Progress: %.1f%% (Generation %d/%d)", progress*100, current, total)
	ebitenutil.DebugPrintAt(screen, progressText, int(barX)+10, int(barY)+10)
}
```

---

## internal\game\state.go

<!-- source: internal\game\state.go -->

```go
package game

// State represents game state
type State int

const (
	StateMenu State = iota
	StateTraining
	StatePlaying
	StateGameOver
)
```

---

## internal\snake\direction.go

<!-- source: internal\snake\direction.go -->

```go
package snake

// Direction represents snake movement direction
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

// IsOpposite checks if direction is opposite to another
func (d Direction) IsOpposite(other Direction) bool {
	return (d == Up && other == Down) ||
		(d == Down && other == Up) ||
		(d == Left && other == Right) ||
		(d == Right && other == Left)
}

// ToVector converts direction to vector
func (d Direction) ToVector() Point {
	switch d {
	case Up:
		return Point{0, -1}
	case Down:
		return Point{0, 1}
	case Left:
		return Point{-1, 0}
	case Right:
		return Point{1, 0}
	default:
		return Point{0, 0}
	}
}
```

---

## internal\snake\point.go

<!-- source: internal\snake\point.go -->

```go
package snake

// Point represents coordinates on grid
type Point struct {
	X, Y int
}

// Equal checks if two points are equal
func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// ManhattanDistance calculates Manhattan distance between points
func (p Point) ManhattanDistance(other Point) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
```

---

## internal\snake\snake.go

<!-- source: internal\snake\snake.go -->

```go
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
```

---

## internal\ui\colors.go

<!-- source: internal\ui\colors.go -->

```go
package ui

import "image/color"

var (
	Background    = color.RGBA{20, 20, 30, 255}
	Grid          = color.RGBA{50, 50, 60, 255}
	Obstacle      = color.RGBA{255, 200, 0, 255}
	ObstacleBorder = color.RGBA{180, 140, 0, 255}
	Food          = color.RGBA{255, 80, 80, 255}
	FoodBorder    = color.RGBA{200, 50, 50, 255}
	SnakeHead     = color.RGBA{100, 255, 100, 255}
	SnakeBody     = color.RGBA{100, 200, 100, 255}
	TextBg        = color.RGBA{0, 0, 0, 180}
)
```

---

