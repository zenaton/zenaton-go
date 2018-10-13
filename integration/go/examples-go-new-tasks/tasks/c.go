package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var C = task.NewCustom("TaskC", &CData{})

type CData struct {
	Workflow string
}

func (a *CData) Handle() (interface{}, error) {
	fmt.Println("\n" + a.Workflow + ": Task C starts")
	time.Sleep(7 * time.Second)
	fmt.Println("\n" + a.Workflow + ": Task C ends")

	return 2, nil
}

func (a *CData) Init(workflow string) {
	a.Workflow = workflow
}
