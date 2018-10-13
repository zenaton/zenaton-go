package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var B = task.NewCustom("TaskB", &BData{})

type BData struct {
	Workflow string
}

func (b *BData) Handle() (interface{}, error) {
	fmt.Println("\n" + b.Workflow + ": Task B starts")
	time.Sleep(5 * time.Second)
	fmt.Println("\n" + b.Workflow + ": Task B ends")

	return 1, nil
}

func (b *BData) Init(workflow string) {
	b.Workflow = workflow
}
