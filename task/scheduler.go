package task

import (
	"sync"
	"time"
)

type scheduler struct {
	taskList map[int64]*Task
	lock     sync.Mutex
}

var Scheduler *scheduler

func init() {
	Scheduler = &scheduler{
		taskList: make(map[int64]*Task),
		lock:     sync.Mutex{},
	}
}

func (s *scheduler) AddTask(job func(), duration time.Duration) int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	task := &Task{
		Job:      job,
		duration: duration,
		stop:     make(chan bool, 1),
	}
	id := time.Now().UnixNano()
	s.taskList[id] = task
	go task.Run()
	return id
}

func (s *scheduler) RemoveTask(taskId int64) {
	if task, exist := s.taskList[taskId]; exist {
		task.stop <- true
	}
}
