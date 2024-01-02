package taskqueue

import (
	"log"
	"sort"
	"sync"
	"time"
)

const (
	TaskStatusNotRunning = 0
	TaskStatusRunning    = 1
	TaskStatusDone       = 2
	TaskStatusFail       = 3
)

type Task struct {
	ID      string
	Status  int32
	Exec    func() int32
	TimeAdd time.Time
}

type TaskQueue struct {
	Tasks chan *Task
	Peek  map[string]*Task
	mu    sync.Mutex
	wg    sync.WaitGroup
}

func NewTaskQueue() *TaskQueue {
	t := TaskQueue{
		Tasks: make(chan *Task, 1000),
		Peek:  make(map[string]*Task),
	}
	t.Start()
	return &t
}

func (q *TaskQueue) Add(id string, exec func() int32) int32 {
	t := Task{
		ID:      id,
		Status:  TaskStatusNotRunning,
		Exec:    exec,
		TimeAdd: time.Now(),
	}
	return q.addTask(&t)
}

func (q *TaskQueue) addTask(t *Task) int32 {
	q.mu.Lock()
	_, ok := q.Peek[t.ID]
	if !ok {
		q.Peek[t.ID] = t
	}
	q.mu.Unlock()

	if ok {
		log.Println("task id", t.ID, "already exist")
		return -1
	}

	q.Tasks <- t
	q.wg.Add(1)
	return 0
}

func (q *TaskQueue) Start() {
	go func() {
		for task := range q.Tasks {
			task.Status = TaskStatusRunning
			log.Println(task.ID, "is running")

			ret := task.Exec()

			if ret == 0 {
				task.Status = TaskStatusDone
				log.Println(task.ID, "is done")
			} else {
				task.Status = TaskStatusFail
				log.Println(task.ID, "is done")
			}

			q.mu.Lock()
			delete(q.Peek, task.ID)
			q.mu.Unlock()

			q.wg.Done()
		}
	}()
}

func (q *TaskQueue) WaitToStop() {
	q.wg.Wait()
	close(q.Tasks)
}

func (q *TaskQueue) PeekTask() []*Task {
	var result []*Task
	for _, v := range q.Peek {
		result = append(result, v)
	}

	// Sort slice based on values
	sort.Slice(result, func(i, j int) bool {
		return result[i].TimeAdd.Before(result[j].TimeAdd)
	})

	return result
}
