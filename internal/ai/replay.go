package ai

import (
	"math/rand/v2"
	"sync"
)

// Experience представляет один опыт для обучения
type Experience struct {
	State     []float64
	Action    int
	Reward    float64
	NextState []float64
	Done      bool
}

// ReplayBuffer реализует Experience Replay буфер
type ReplayBuffer struct {
	buffer   []Experience
	capacity int
	mu       sync.Mutex
}

// NewReplayBuffer создает новый буфер опыта
func NewReplayBuffer(capacity int) *ReplayBuffer {
	return &ReplayBuffer{
		buffer:   make([]Experience, 0, capacity),
		capacity: capacity,
	}
}

// Add добавляет опыт в буфер
func (rb *ReplayBuffer) Add(exp Experience) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.buffer) < rb.capacity {
		rb.buffer = append(rb.buffer, exp)
	} else {
		// Удаляем самый старый опыт
		rb.buffer = append(rb.buffer[1:], exp)
	}
}

// Sample возвращает случайную выборку опытов
func (rb *ReplayBuffer) Sample(batchSize int) []Experience {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.buffer) < batchSize {
		batchSize = len(rb.buffer)
	}

	samples := make([]Experience, batchSize)
	
	// ✅ ИСПРАВЛЕНИЕ: используем math/rand/v2
	indices := rand.Perm(len(rb.buffer))[:batchSize]

	for i, idx := range indices {
		samples[i] = rb.buffer[idx]
	}

	return samples
}

// Size возвращает текущий размер буфера
func (rb *ReplayBuffer) Size() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer)
}

// Clear очищает буфер
func (rb *ReplayBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.buffer = make([]Experience, 0, rb.capacity)
}

// IsFull проверяет заполнен ли буфер
func (rb *ReplayBuffer) IsFull() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer) >= rb.capacity
}
