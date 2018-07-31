package zenaton

import "reflect"

var instance *Engine

type Engine struct {
	client    *Client
	processor Processor
}

func NewEngine() *Engine {
	if instance == nil {
		instance = &Engine{
			client: NewClient(false),
		}
	}
	return instance
}

// todo: maybe I don't want this to be exported, so only tasks and workflows can implement this interface
type Job interface {
	Handle() interface{}
	AsyncHandle(chan interface{})
	GetName() string
	GetData() interface{}
}

type Processor interface {
	Process([]Job, bool) []interface{}
}

type chanResult struct {
	result interface{}
	index  int
}

func wrapper(index int, outcome chan chanResult, handle func() interface{}) {
	outcome <- chanResult{
		index:  index,
		result: handle(),
	}
}

//todo: error handling
func (e *Engine) Execute(jobs []Job) []interface{} {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		outputChan := make(chan chanResult)
		for i, job := range jobs {
			go wrapper(i, outputChan, job.Handle)
		}

		outputs := make([]interface{}, len(jobs))
		for range jobs {
			o, ok := <-outputChan
			if !ok {
				break
			}
			outputs[o.index] = o.result
		}

		return outputs
	}

	return e.processor.Process(jobs, true)
}

func (e *Engine) Dispatch(jobs []Job) []chan interface{} {
	// local execution
	var chans []chan interface{}
	for range jobs {
		chans = append(chans, make(chan interface{}))
	}

	if !reflect.ValueOf(e.processor).IsValid() || len(jobs) == 0 {
		for i, job := range jobs {
			go job.AsyncHandle(chans[i])
		}
	}
	return chans
	//todo: (figure out what to do with processor here) return e.processor.Process(jobs, false)
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
