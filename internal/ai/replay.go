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
