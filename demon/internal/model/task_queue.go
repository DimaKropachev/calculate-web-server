package model

type TaskQueue struct {
	Tasks   chan *Task
	Results chan *Response
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		Tasks: make(chan *Task),
		Results: make(chan *Response),
	}
}