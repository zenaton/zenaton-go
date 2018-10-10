package workflows

import (
	"fmt"

	"github.com/zenaton/integration/go/tasks"
	"github.com/zenaton/integration/go/tasks/log"
	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WithTaskWorkflow = workflow.NewCustom("WithTaskWorkflow", &WithTask{})

type WithTask struct {
	Task *task.Instance
}

func (wt *WithTask) Init(t *task.Instance) {
	wt.Task = t
}

func (wt *WithTask) Handle() (interface{}, error) {
	fmt.Println("wt: ", wt.Task)
	wt.Task.Execute()
	return nil, nil
}

var TestW = workflow.NewCustom("TestWorkflow", &Test2{})

type Test2 struct {
	Relaunch   bool
	Parallel   bool
	Print      string
	Return     string
	Error      string
	Panic      string
	TaskReturn interface{}
	TaskError  string
	TaskPanic  string
	TaskPrint  string
	IDstr      string
}

func (t *Test2) Init(t2 Test2) {
	t.Relaunch = t2.Relaunch
	t.Parallel = t2.Parallel
	t.Print = t2.Print
	t.Return = t2.Return
	t.Error = t2.Error
	t.Panic = t2.Panic
	t.TaskReturn = t2.TaskReturn
	t.TaskError = t2.TaskError
	t.TaskPanic = t2.TaskPanic
	t.TaskPrint = t2.TaskPrint
	t.IDstr = t2.IDstr
}

func (t *Test2) Handle() (interface{}, error) {
	if t.Parallel {

		var t1Return tasks.Test
		var t2Return tasks.Test

		t1 := tasks.TestTask.New(&tasks.Test{
			Return: t.TaskReturn,
			Error:  t.TaskError,
			Panic:  t.TaskPanic,
			Print:  t.TaskPrint,
		})

		t2 := t1

		task.Parallel{t1, t2}.Execute().Output(&t1Return, &t2Return)
		log.Println("out1: ", t1Return)
		log.Println("out2: ", t2Return)

	} else {
		var out interface{}
		err := tasks.TestTask.New(&tasks.Test{
			Return:   t.TaskReturn,
			Error:    t.TaskError,
			Panic:    t.TaskPanic,
			Print:    t.TaskPrint,
			Relaunch: t.Relaunch,
		}).Execute().Output(&out)

		log.Println("out: ", out)
		log.Println("err: ", err)
	}

	if t.Panic != "" {
		panic(t.Panic)
	}

	if t.Error != "" {
		return nil, errors.New("testWorkflowError", "testErrorMessage")
	}

	return t.Return, nil
}

func (t *Test2) ID() string {
	if t.IDstr != "" {
		return t.IDstr
	}
	return "MyID"
}

var TestRelaunchTaskWorkflow = workflow.NewCustom("TestRelaunchTaskWorkflow", &TestRelaunchTask{})

type TestRelaunchTask struct{}

func (t TestRelaunchTask) Handle() (interface{}, error) {

	var out interface{}
	err := tasks.TaskRunnerTask.New().Execute().Output(&out)

	log.Println("out: ", out)
	log.Println("err: ", err)

	return nil, nil
}

var TestEventValueWorkflow = workflow.NewCustom("TestEventValueWorkflow", &TestEventValue{})

type TestEventValue struct{}

func (tev TestEventValue) Handle() (interface{}, error) {

	event := task.Wait().ForEvent("MyEvent").Seconds(5).Execute()
	fmt.Println("wait for event: ", event)

	return nil, nil
}

func (tev *TestEventValue) OnEvent(eventName string, eventData interface{}) {
	fmt.Println("onEvent: ", eventName, eventData)
}

func (tev TestEventValue) ID() string {
	return "TestEventValueID"
}
