package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var D = task.NewCustom("TaskD", &DData{})

type DData struct {
	Workflow string
}

func (a *DData) Handle() (interface{}, error) {
	fmt.Println("\n" + a.Workflow + ": Task D starts")
	time.Sleep(9 * time.Second)
	fmt.Println("\n" + a.Workflow + ": Task D ends")

	return 3, nil
}

func (a *DData) Init(workflow string) {
	a.Workflow = workflow
}
