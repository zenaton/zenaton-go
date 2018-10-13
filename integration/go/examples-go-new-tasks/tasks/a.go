package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var A = task.NewCustom("TaskA", &AData{})

type AData struct {
	Workflow string
}

func (a *AData) Handle() (interface{}, error) {
	fmt.Println("\n" + a.Workflow + ": Task A starts")
	time.Sleep(3 * time.Second)
	fmt.Println("\n" + a.Workflow + ": Task A ends")

	return 0, nil
}

func (a *AData) Init(workflow string) {
	a.Workflow = workflow
}
