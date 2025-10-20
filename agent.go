package main

import (
	"math/rand"
	"sync"
)

type Experience struct {
	State     []float64
	Action    int
	Reward    float64
	NextState []float64
	Done      bool
}

type ReplayBuffer struct {
	buffer   []Experience
	capacity int
	mu       sync.Mutex
}

func NewReplayBuffer(capacity int) *ReplayBuffer {
	return &ReplayBuffer{buffer: make([]Experience, 0, capacity), capacity: capacity}
}

func (rb *ReplayBuffer) Add(exp Experience) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if len(rb.buffer) < rb.capacity {
		rb.buffer = append(rb.buffer, exp)
	} else {
		rb.buffer = append(rb.buffer[1:], exp)
	}
}

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

func (rb *ReplayBuffer) Size() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer)
}

type DQNAgent struct {
	qNetwork      *NeuralNetwork
	targetNetwork *NeuralNetwork
	replayBuffer  *ReplayBuffer
	epsilon       float64
	epsilonMin    float64
	epsilonDecay  float64
	gamma         float64
	batchSize     int
	updateFreq    int
	stepCount     int
}

type AgentConfig struct {
	LearningRate float64
	BufferSize   int
	EpsilonStart float64
	EpsilonMin   float64
	EpsilonDecay float64
	Gamma        float64
	BatchSize    int
	UpdateFreq   int
}

func NewDQNAgent(stateSize, actionSize int, config AgentConfig) *DQNAgent {
	layers := []int{stateSize, 128, 128, actionSize}
	return &DQNAgent{
		qNetwork:      NewNeuralNetwork(layers, config.LearningRate),
		targetNetwork: NewNeuralNetwork(layers, config.LearningRate),
		replayBuffer:  NewReplayBuffer(config.BufferSize),
		epsilon:       config.EpsilonStart,
		epsilonMin:    config.EpsilonMin,
		epsilonDecay:  config.EpsilonDecay,
		gamma:         config.Gamma,
		batchSize:     config.BatchSize,
		updateFreq:    config.UpdateFreq,
	}
}

func (agent *DQNAgent) SelectAction(state []float64) int {
	if rand.Float64() < agent.epsilon {
		return rand.Intn(4)
	}
	return argmax(agent.qNetwork.Forward(state))
}

func argmax(values []float64) int {
	maxIdx, maxVal := 0, values[0]
	for i, v := range values {
		if v > maxVal {
			maxVal, maxIdx = v, i
		}
	}
	return maxIdx
}

func (agent *DQNAgent) Remember(state []float64, action int, reward float64, nextState []float64, done bool) {
	agent.replayBuffer.Add(Experience{state, action, reward, nextState, done})
}

func (agent *DQNAgent) Train() float64 {
	if agent.replayBuffer.Size() < agent.batchSize {
		return 0
	}
	batch := agent.replayBuffer.Sample(agent.batchSize)
	totalLoss := 0.0
	for _, exp := range batch {
		target := agent.qNetwork.Forward(exp.State)
		if exp.Done {
			target[exp.Action] = exp.Reward
		} else {
			nextQValues := agent.targetNetwork.Forward(exp.NextState)
			maxQ := nextQValues[argmax(agent.qNetwork.Forward(exp.NextState))]
			target[exp.Action] = exp.Reward + agent.gamma*maxQ
		}
		totalLoss += agent.qNetwork.BackwardAndUpdate(exp.State, target)
	}
	agent.stepCount++
	if agent.stepCount%agent.updateFreq == 0 {
		agent.UpdateTargetNetwork()
	}
	if agent.epsilon > agent.epsilonMin {
		agent.epsilon *= agent.epsilonDecay
	}
	return totalLoss / float64(len(batch))
}

func (agent *DQNAgent) UpdateTargetNetwork() {
	agent.targetNetwork = agent.qNetwork.Clone()
}

func (agent *DQNAgent) SaveModel(filename string) error {
	return agent.qNetwork.SaveToFile(filename)
}

func (agent *DQNAgent) LoadModel(filename string) error {
	err := agent.qNetwork.LoadFromFile(filename)
	if err == nil {
		agent.targetNetwork = agent.qNetwork.Clone()
	}
	return err
}
