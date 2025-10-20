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
