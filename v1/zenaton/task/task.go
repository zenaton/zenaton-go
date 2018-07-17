package task

import "github.com/zenaton/zenaton-go/v1/zenaton/engine"

type Task struct {
	Name       string
	HandleFunc func() interface{}
	Data       interface{}
}

func (t *Task) Handle() interface{} {
	return t.HandleFunc()
}

func (t *Task) AsyncHandle(channel chan interface{}) {
	channel <- t.HandleFunc()
}

func (t *Task) Execute() interface{} {
	e := engine.New()
	return e.Execute([]engine.Job{t})
}

func (t *Task) Dispatch() chan interface{} {
	e := engine.New()
	return e.Dispatch([]engine.Job{t})[0]
}

func (ts Tasks) Dispatch() chan interface{} {
	e := engine.New()
	var jobs []engine.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)[0]
}

type Tasks []*Task

func (ts Tasks) Execute() []interface{} {
	e := engine.New()
	var jobs []engine.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}
