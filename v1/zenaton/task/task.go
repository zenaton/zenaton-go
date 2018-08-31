package task

import (
	"reflect"

	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
)

type Task struct {
	name string
	//todo: would be nice if the handle func could take many arguments, instead of just one. would have to think how that would be done (maybe pass in argments into execute?)
	handler interfaces.Handler
	//todo: MaxProcessingTime func() int64
}

func (t Task) GetName() string { return t.name }

//todo: change
func (t Task) GetData() interface{} { return t.handler }

func (t Task) Handle() (interface{}, error) {
	t.handler.Handle()
	return "", nil
}

func (t Task) Async() error {
	t.handler.Handle()
	return nil
}

type MaxProcessingTimer interface {
	MaxProcessingTime() int64
}

func (t Task) MaxProcessingTime() int64 {
	maxer, ok := t.handler.(MaxProcessingTimer)
	if ok {
		return maxer.MaxProcessingTime()
	}
	return -1
}

func New(h interfaces.Handler) *Task {

	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewTask")
	}

	task := Task{
		name:    reflect.TypeOf(h).Elem().Name(),
		handler: h,
	}

	NewTaskManager().setClass(task.name, &task)
	return &task
}

//todo: would be great if we could take a pointer to execute and modify that like json.unmarshal does, but it's hard to figure out how they do it
func (t *Task) Execute() (interface{}, error) {
	e := engine.NewEngine()
	output, err := e.Execute([]interfaces.Job{t})
	//todo: make sure this is impossible to get index out of bounds
	if output == nil {
		return nil, err
	}
	return output[0], err
}

func (t *Task) Dispatch() error {
	e := engine.NewEngine()
	err := e.Dispatch([]interfaces.Job{t})
	return err
}

type Parallel []*Task

func (ts Parallel) Dispatch() error {
	e := engine.NewEngine()
	var jobs []interfaces.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)
}

func (ts Parallel) Execute() ([]interface{}, error) {
	e := engine.NewEngine()
	var jobs []interfaces.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}
