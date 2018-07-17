package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskA = &task.Task{
	Name: "TaskA",
	HandleFunc: func() interface{} {
		fmt.Println("Task A")
		time.Sleep(3 * time.Millisecond)
		//todo: figure out what's happening with the done function
		fmt.Println("done with task A")
		return "TaskA"
	},
}

var TaskB = &task.Task{
	Name: "TaskB",
	HandleFunc: func() interface{} {
		fmt.Println("Task B")
		time.Sleep(5 * time.Millisecond)
		//todo: figure out what's happening with the done function
		fmt.Println("done with task B")
		return "TaskB"
	},
}

var TaskC = &task.Task{
	Name: "TaskC",
	HandleFunc: func() interface{} {
		fmt.Println("Task C")
		time.Sleep(7 * time.Millisecond)
		//todo: figure out what's happening with the done function
		fmt.Println("done with task C")
		return "TaskC"
	},
}

//var TaskA = task.New("TaskA",
//	// todo: this is ugly, would be nice to not have to define your functions with wait groups might need to though.
//	func() interface{} {
//		fmt.Println("Task A")
//		time.Sleep(3 * time.Millisecond)
//		//todo: figure out what's happening with the done function
//		fmt.Println("done with task A")
//		return nil
//	})
//
//var TaskB = task.New("TaskB",
//	// todo: this is ugly, would be nice to not have to define your functions with wait groups might need to though.
//	func() interface{} {
//		fmt.Println("Task B")
//		time.Sleep(3 * time.Millisecond)
//		//todo: figure out what's happening with the done function
//		fmt.Println("done with task B")
//		return nil
//	})
//
//var TaskC = task.New("TaskC",
//	// todo: this is ugly, would be nice to not have to define your functions with wait groups might need to though.
//	func() interface{} {
//		fmt.Println("Task C")
//		time.Sleep(3 * time.Millisecond)
//		//todo: figure out what's happening with the done function
//		fmt.Println("done with task c")
//		return nil
//	})
