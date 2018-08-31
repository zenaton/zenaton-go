package engine

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
)

var instance *Engine

type Engine struct {
	client     *client.Client
	processor  Processor
	processor2 Processor2
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
	Process([]interfaces.Job, bool) ([]interface{}, error)
}
type Processor2 interface {
	Process2([]interfaces.Handler, bool) ([]interface{}, error)
}

//todo: error handling
func (e *Engine) Execute(jobs []interfaces.Job) ([]interface{}, error) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {

		var outputs []interface{}
		var output interface{}
		var err error

		for _, job := range jobs {
			output, err = job.Handle()
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, output)
		}

		return outputs, nil
	}

	return e.processor.Process(jobs, true)
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

	_, err := e.processor.Process(jobs, false)
	return err
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
