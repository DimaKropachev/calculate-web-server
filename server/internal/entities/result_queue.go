package entities

import (
	"sync"

	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
)

type ResultsQueue struct {
	results map[string]*models.Result
	mu      sync.RWMutex
}

func NewResultsQueue() *ResultsQueue {
	return &ResultsQueue{
		results: make(map[string]*models.Result),
		mu:      sync.RWMutex{},
	}
}

func (rq *ResultsQueue) Get(id string) (*models.Result, bool) {
	rq.mu.RLock()
	defer rq.mu.RUnlock()

	res, ok := rq.results[id]

	return res, ok
}

func (rq *ResultsQueue) Add(result *models.Result) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	rq.results[result.Id] = result
}
