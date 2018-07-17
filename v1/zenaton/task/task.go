package task

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/job"
)

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
	return e.Execute([]job.Job{t})
}

func (t *Task) Dispatch() chan interface{} {
	e := engine.New()
	return e.Dispatch([]job.Job{t})[0]
}

func (ts Tasks) Dispatch() chan interface{} {
	e := engine.New()
	var jobs []job.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)[0]
}

type Tasks []*Task

func (ts Tasks) Execute() []interface{} {
	e := engine.New()
	var jobs []job.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}
