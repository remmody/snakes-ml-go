package ai

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

func DefaultConfig() Config {
	return Config{
		LearningRate: 0.001,
		BufferSize:   1000000,
		EpsilonStart: 1.0,
		EpsilonMin:   0.01,
		EpsilonDecay: 0.995,
		Gamma:        0.95,
		BatchSize:    64,
		UpdateFreq:   100,
	}
}
