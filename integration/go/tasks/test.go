package tasks

import (
	"fmt"

	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TestTask = task.NewCustom("TestTask", &Test{})

type Test struct {
	Return   interface{}
	Error    string
	Panic    string
	Print    string
	Relaunch bool
}

func (t *Test) Init(t2 *Test) {
	t.Return = t2.Return
	t.Error = t2.Error
	t.Panic = t2.Panic
	t.Print = t2.Print
	t.Relaunch = t2.Relaunch
}

func (t *Test) Handle() (interface{}, error) {
	if t.Print != "" {
		fmt.Println(t.Print)
	}

	if t.Panic != "" {
		panic(t.Panic)
	}

	var err error
	if t.Error != "" {
		err = errors.New("testErrorName", t.Error)
	}

	if t.Relaunch {
		var out interface{}
		err := TestTask.New(&Test{
			Return: t.Return,
			Error:  t.Error,
		}).Execute().Output(&out)
		return out, err
	}

	return t.Return, err
}

var TaskRunnerTask = task.New("TaskRunnerTask",
	func() (interface{}, error) {
		//without return
		TestTask.New().Execute()

		//with return
		var out string
		err := TestTask.New(&Test{
			Return: "test return",
			Error:  "test error",
		}).Execute().Output(&out)

		return out, err
	})
