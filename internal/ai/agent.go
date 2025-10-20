package ai

import "math/rand/v2"

// Agent представляет DQN агента с системой поколений
type Agent struct {
	qNetwork         *Network
	targetNetwork    *Network
	replayBuffer     *ReplayBuffer
	epsilon          float64
	epsilonMin       float64
	epsilonDecay     float64
	gamma            float64
	batchSize        int
	updateFreq       int
	stepCount        int
	episodeCount     int          // ✅ НОВОЕ: счетчик эпизодов
	generationSize   int          // ✅ НОВОЕ: размер поколения
	currentGeneration int         // ✅ НОВОЕ: текущее поколение
	totalReward      float64      // ✅ НОВОЕ: накопленная награда
	episodeRewards   []float64    // ✅ НОВОЕ: награды эпизодов
}

// NewAgent создает нового DQN агента
func NewAgent(stateSize, actionSize int, config Config) *Agent {
	layers := []int{stateSize, 128, 128, actionSize}

	return &Agent{
		qNetwork:         NewNetwork(layers, config.LearningRate),
		targetNetwork:    NewNetwork(layers, config.LearningRate),
		replayBuffer:     NewReplayBuffer(config.BufferSize),
		epsilon:          config.EpsilonStart,
		epsilonMin:       config.EpsilonMin,
		epsilonDecay:     config.EpsilonDecay,
		gamma:            config.Gamma,
		batchSize:        config.BatchSize,
		updateFreq:       config.UpdateFreq,
		stepCount:        0,
		episodeCount:     0,
		generationSize:   100,        // ✅ 100 эпизодов = 1 поколение
		currentGeneration: 1,
		episodeRewards:   make([]float64, 0, 100),
	}
}

// SelectAction выбирает действие используя epsilon-greedy стратегию
func (a *Agent) SelectAction(state []float64) int {
	if rand.Float64() < a.epsilon {
		return rand.IntN(4)
	}

	qValues := a.qNetwork.Forward(state)
	return argmax(qValues)
}

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
	
	// ✅ НОВОЕ: накапливаем награду эпизода
	a.totalReward += reward
}

// Train выполняет один шаг обучения
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

	// ✅ ИСПРАВЛЕНО: Epsilon decay после каждого обучения, а не после эпизода
	if a.epsilon > a.epsilonMin {
		a.epsilon *= a.epsilonDecay
	}

	return totalLoss / float64(len(batch))
}

// ✅ НОВОЕ: EndEpisode вызывается в конце каждого эпизода
func (a *Agent) EndEpisode() {
	a.episodeCount++
	a.episodeRewards = append(a.episodeRewards, a.totalReward)
	
	// Ограничиваем хранение наград последними 100 эпизодами
	if len(a.episodeRewards) > 100 {
		a.episodeRewards = a.episodeRewards[1:]
	}
	
	// Сброс накопленной награды
	a.totalReward = 0
	
	// Проверка смены поколения
	if a.episodeCount%a.generationSize == 0 {
		a.currentGeneration++
	}
}

// ✅ НОВОЕ: GetAverageReward возвращает среднюю награду за последние эпизоды
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

// Геттеры
func (a *Agent) Epsilon() float64           { return a.epsilon }
func (a *Agent) SetEpsilon(epsilon float64) { a.epsilon = epsilon }
func (a *Agent) ReplayBufferSize() int      { return a.replayBuffer.Size() }
func (a *Agent) StepCount() int             { return a.stepCount }
func (a *Agent) EpisodeCount() int          { return a.episodeCount }
func (a *Agent) Generation() int            { return a.currentGeneration }
func (a *Agent) GenerationProgress() int    { return a.episodeCount % a.generationSize }
