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
