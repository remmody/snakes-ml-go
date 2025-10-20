package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
)

type NeuralNetwork struct {
	layers       []int
	weights      [][][]float64
	biases       [][]float64
	learningRate float64
	mu           sync.RWMutex
}

func NewNeuralNetwork(layers []int, learningRate float64) *NeuralNetwork {
	nn := &NeuralNetwork{layers: layers, learningRate: learningRate}
	nn.weights = make([][][]float64, len(layers)-1)
	nn.biases = make([][]float64, len(layers)-1)

	for i := 0; i < len(layers)-1; i++ {
		nn.weights[i] = make([][]float64, layers[i])
		nn.biases[i] = make([]float64, layers[i+1])
		limit := math.Sqrt(6.0 / float64(layers[i]+layers[i+1]))
		for j := 0; j < layers[i]; j++ {
			nn.weights[i][j] = make([]float64, layers[i+1])
			for k := 0; k < layers[i+1]; k++ {
				nn.weights[i][j][k] = (rand.Float64()*2 - 1) * limit
			}
		}
		for k := 0; k < layers[i+1]; k++ {
			nn.biases[i][k] = (rand.Float64()*2 - 1) * limit
		}
	}
	return nn
}

func relu(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}

func reluDerivative(x float64) float64 {
	if x > 0 {
		return 1
	}
	return 0
}

func (nn *NeuralNetwork) Forward(input []float64) []float64 {
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

func (nn *NeuralNetwork) BackwardAndUpdate(input, target []float64) float64 {
	nn.mu.Lock()
	defer nn.mu.Unlock()
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

	deltas := make([][]float64, len(nn.layers)-1)
	lastIdx := len(activations) - 1
	deltas[lastIdx-1] = make([]float64, nn.layers[lastIdx])
	loss := 0.0
	for i := 0; i < nn.layers[lastIdx]; i++ {
		error := target[i] - activations[lastIdx][i]
		deltas[lastIdx-1][i] = error
		loss += error * error
	}

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

func (nn *NeuralNetwork) Clone() *NeuralNetwork {
	nn.mu.RLock()
	defer nn.mu.RUnlock()
	clone := &NeuralNetwork{
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

func (nn *NeuralNetwork) SaveToFile(filename string) error {
	nn.mu.RLock()
	defer nn.mu.RUnlock()
	data, err := json.Marshal(struct {
		Layers  []int         `json:"layers"`
		Weights [][][]float64 `json:"weights"`
		Biases  [][]float64   `json:"biases"`
	}{Layers: nn.layers, Weights: nn.weights, Biases: nn.biases})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func (nn *NeuralNetwork) LoadFromFile(filename string) error {
	nn.mu.Lock()
	defer nn.mu.Unlock()
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var loaded struct {
		Layers  []int         `json:"layers"`
		Weights [][][]float64 `json:"weights"`
		Biases  [][]float64   `json:"biases"`
	}
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	nn.layers = loaded.Layers
	nn.weights = loaded.Weights
	nn.biases = loaded.Biases
	return nil
}
