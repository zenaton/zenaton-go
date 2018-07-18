package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var (
	TaskA = &task.Task{
		Name: "TaskA",
		HandleFunc: func() interface{} {
			fmt.Println("Task A")
			time.Sleep(3 * time.Millisecond)
			//todo: figure out what's happening with the done function
			return "TaskA"
		},
	}

	TaskB = &task.Task{
		Name: "TaskB",
		HandleFunc: func() interface{} {
			fmt.Println("Task B")
			time.Sleep(5 * time.Millisecond)
			return "TaskB"
		},
	}

	TaskC = &task.Task{
		Name: "TaskC",
		HandleFunc: func() interface{} {
			fmt.Println("Task C")
			time.Sleep(7 * time.Millisecond)
			return "TaskC"
		},
	}

	TaskD = &task.Task{
		Name: "Task D",
		HandleFunc: func() interface{} {
			fmt.Println("Task D")
			time.Sleep(9 * time.Millisecond)
			return "Task D"
		},
	}
)
