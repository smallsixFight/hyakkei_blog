package task

import "time"

type Task struct {
	Job      func()
	duration time.Duration
	stop     chan bool
}

func (t *Task) Run() {
	ti := time.NewTimer(t.duration)
	for {
		select {
		case <-ti.C:
			go t.Job()
			t.stop <- true
		case <-t.stop:
			return
		}
	}
}
