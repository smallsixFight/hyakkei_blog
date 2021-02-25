package task

import (
	"fmt"
	"testing"
	"time"
)

func TestTask_Run(t *testing.T) {
	task := &Task{
		Job: func() {
			fmt.Println("run task...")
		},
		duration: time.Second * 2,
		stop:     make(chan bool, 1),
	}
	go task.Run()
	task.stop <- true
	time.Sleep(time.Second * 4)
}
