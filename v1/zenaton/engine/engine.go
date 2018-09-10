package engine

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
)

var instance *Engine

type Engine struct {
	client    *client.Client
	processor Processor
}

func NewEngine() *Engine {
	if instance == nil {
		instance = &Engine{
			client: client.NewClient(false),
		}
	}
	return instance
}

type Processor interface {
	Process([]interfaces.Job, bool, ...interface{}) error
}

//todo: error handling
func (e *Engine) Execute(jobs []interfaces.Job, outputs []interface{}) error {

	// local execution
	if e.processor == nil || len(jobs) == 0 {

		//var outputs []interface{}
		//var output interface{}
		var err error

		for _, job := range jobs {
			_, err = job.Handle()
			if err != nil {
				return err
			}
			//outputs = append(outputs, output)
		}

		return nil
	}

	return e.processor.Process(jobs, true, outputs...)
}

func (e *Engine) Dispatch(jobs []interfaces.Job) error {

	if e.processor == nil || len(jobs) == 0 {

		var err error

		for _, job := range jobs {
			err = job.Async()
			if err != nil {
				return err
			}
		}

		return nil
	}

	return e.processor.Process(jobs, false)
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
