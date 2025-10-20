package ai

import "math/rand/v2"

// Agent представляет DQN агента
type Agent struct {
	qNetwork      *Network
	targetNetwork *Network
	replayBuffer  *ReplayBuffer
	epsilon       float64
	epsilonMin    float64
	epsilonDecay  float64
	gamma         float64
	batchSize     int
	updateFreq    int
	stepCount     int
}

// NewAgent создает нового DQN агента
func NewAgent(stateSize, actionSize int, config Config) *Agent {
	layers := []int{stateSize, 128, 128, actionSize}

	return &Agent{
		qNetwork:      NewNetwork(layers, config.LearningRate),
		targetNetwork: NewNetwork(layers, config.LearningRate),
		replayBuffer:  NewReplayBuffer(config.BufferSize),
		epsilon:       config.EpsilonStart,
		epsilonMin:    config.EpsilonMin,
		epsilonDecay:  config.EpsilonDecay,
		gamma:         config.Gamma,
		batchSize:     config.BatchSize,
		updateFreq:    config.UpdateFreq,
		stepCount:     0,
	}
}

// SelectAction выбирает действие используя epsilon-greedy стратегию
func (a *Agent) SelectAction(state []float64) int {
	// ✅ ИСПРАВЛЕНИЕ: используем math/rand/v2
	if rand.Float64() < a.epsilon {
		return rand.IntN(4) // 4 действия: Up, Right, Down, Left
	}

	qValues := a.qNetwork.Forward(state)
	return argmax(qValues)
}

// argmax возвращает индекс максимального элемента
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

// Remember сохраняет опыт в replay buffer
func (a *Agent) Remember(state []float64, action int, reward float64, nextState []float64, done bool) {
	a.replayBuffer.Add(Experience{
		State:     state,
		Action:    action,
		Reward:    reward,
		NextState: nextState,
		Done:      done,
	})
}

// Train выполняет один шаг обучения
func (a *Agent) Train() float64 {
	if a.replayBuffer.Size() < a.batchSize {
		return 0
	}

	batch := a.replayBuffer.Sample(a.batchSize)
	totalLoss := 0.0

	for _, exp := range batch {
		// Получаем текущие Q-значения
		target := a.qNetwork.Forward(exp.State)

		if exp.Done {
			// Если эпизод завершен, целевое значение = награда
			target[exp.Action] = exp.Reward
		} else {
			// Double DQN: используем q-network для выбора действия,
			// target-network для оценки
			nextQValues := a.targetNetwork.Forward(exp.NextState)
			bestAction := argmax(a.qNetwork.Forward(exp.NextState))
			maxQ := nextQValues[bestAction]

			// Bellman equation
			target[exp.Action] = exp.Reward + a.gamma*maxQ
		}

		// Обратное распространение
		loss := a.qNetwork.BackwardAndUpdate(exp.State, target)
		totalLoss += loss
	}

	// Обновление target network
	a.stepCount++
	if a.stepCount%a.updateFreq == 0 {
		a.UpdateTargetNetwork()
	}

	// Уменьшение epsilon
	if a.epsilon > a.epsilonMin {
		a.epsilon *= a.epsilonDecay
	}

	return totalLoss / float64(len(batch))
}

// UpdateTargetNetwork обновляет target network
func (a *Agent) UpdateTargetNetwork() {
	a.targetNetwork = a.qNetwork.Clone()
}

// SaveModel сохраняет модель в файл
func (a *Agent) SaveModel(filename string) error {
	return a.qNetwork.SaveToFile(filename)
}

// LoadModel загружает модель из файла
func (a *Agent) LoadModel(filename string) error {
	err := a.qNetwork.LoadFromFile(filename)
	if err == nil {
		a.targetNetwork = a.qNetwork.Clone()
	}
	return err
}

// Epsilon возвращает текущее значение epsilon
func (a *Agent) Epsilon() float64 {
	return a.epsilon
}

// SetEpsilon устанавливает значение epsilon
func (a *Agent) SetEpsilon(epsilon float64) {
	a.epsilon = epsilon
}

// ReplayBufferSize возвращает размер replay buffer
func (a *Agent) ReplayBufferSize() int {
	return a.replayBuffer.Size()
}

// StepCount возвращает количество шагов обучения
func (a *Agent) StepCount() int {
	return a.stepCount
}
	