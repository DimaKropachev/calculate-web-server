package entities

import (
	"strconv"
	"sync"
)

type SafeId struct {
	counter int
	mu      sync.Mutex
}

func NewSafeId() *SafeId {
	return &SafeId{
		counter: 0,
		mu: sync.Mutex{},
	}
}

func (id *SafeId) Get() string {
	id.mu.Lock()
	id.counter++
	id.mu.Unlock()

	return strconv.Itoa(id.counter)
}