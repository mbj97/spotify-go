package utils

import (
	"fmt"
	"sync"
	"time"
)

type TaskScheduler struct {
	Tasks []Task
	mutex sync.Mutex
}

type Task struct {
	ID       string
	Function func()
	Interval int
	stopCh   chan struct{}
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{}
}

func (ts *TaskScheduler) CreateTask(id string, f func(), interval int) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task := Task{
		ID:       id,
		Function: f,
		Interval: interval,
		stopCh:   make(chan struct{}),
	}

	ts.Tasks = append(ts.Tasks, task)
}

func (ts *TaskScheduler) StartTaskByID(id string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	for _, task := range ts.Tasks {
		if task.ID == id {
			go ts.runTask(&task)
			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

func (ts *TaskScheduler) StopTaskByID(id string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	for _, task := range ts.Tasks {
		if task.ID == id {
			close(task.stopCh)
			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

func (ts *TaskScheduler) Start() {
	for _, task := range ts.Tasks {
		go ts.runTask(&task)
	}
}

func (ts *TaskScheduler) runTask(task *Task) {
	ticker := time.NewTicker(time.Duration(task.Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			task.Function()
		case <-task.stopCh:
			ticker.Stop()
			return
		}
	}
}