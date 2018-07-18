package recursive

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

func NewDisplayTask(id string) *task.Task {

	return &task.Task{
		Name: "DisplayTask",
		HandleFunc: func() interface{} {
			fmt.Print(id)
			time.Sleep(1 * time.Second)
			return nil
		},
		ID: func() string {
			return id
		},
	}
}

func NewRelaunchTask(id int, max int) *task.Task {
	return &task.Task{
		Name: "RelaunchTask",
		HandleFunc: func() interface{} {
			if id >= max {
				return nil
			}

			newID := 1 + id
			fmt.Printf("\nIteration: %v", newID)

			//todo: do I need to catch this result?
			NewRecursiveWorkflow(newID, max).Dispatch()
			return nil
		},

		ID: func() string {
			return strconv.Itoa(id)
		},
	}
}
