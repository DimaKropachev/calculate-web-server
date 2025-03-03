package entities

import (
	"sync"

	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
)

type TasksQueue struct {
	tasks []*models.FinalTask
	mu    sync.Mutex
}

func NewTasksQueue() *TasksQueue {
	return &TasksQueue{
		mu: sync.Mutex{},
	}
}

func (tq *TasksQueue) Get() *models.FinalTask {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.tasks) == 0 {
		return nil
	}
	task := tq.tasks[0]
	tq.tasks = tq.tasks[1:]

	return task
}

func (tq *TasksQueue) Add(task *models.FinalTask) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.tasks = append(tq.tasks, task)
}
