package engine

import (
	"sync"

	"github.com/zenaton/zenaton-go/v1/zenaton/internal/client"
)

var instance = &Engine{
	client: client.NewClient(false),
}

type Engine struct {
	client    *client.Client
	processor Processor
}

func NewEngine() *Engine {
	return instance
}

type Processor interface {
	Process([]Job, bool) ([]interface{}, []string, []error)
}

type LaunchInfo struct {
	Type      string
	Name      string
	Canonical string
	ID        string
	Data      interface{}
}

type Handler interface {
	Handle() (interface{}, error)
}

type Job interface {
	Handle() (interface{}, error)
	GetName() string
	GetData() Handler
}

func (e *Engine) Execute(jobs []Job) ([]interface{}, []string, []error) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		var outputs []interface{}
		var errs []error
		mu := &sync.Mutex{}

		wg := sync.WaitGroup{}
		for _, job := range jobs {

			job := job //gotcha!

			wg.Add(1)
			go func() {

				defer wg.Done()

				out, err := job.GetData().Handle()
				mu.Lock()
				errs = append(errs, err)
				outputs = append(outputs, out)
				mu.Unlock()
			}()
		}
		wg.Wait()

		return outputs, nil, errs
	}

	outputValues, serializedOutputs, errs := e.processor.Process(jobs, true)
	return outputValues, serializedOutputs, errs
}

func (e *Engine) Dispatch(jobs []Job, launchInfos []LaunchInfo) {

	if e.processor == nil || len(jobs) == 0 {

		for i, job := range jobs {
			job := job
			li := launchInfos[i]
			if li.Type == "workflow" {
				client.NewClient(false).StartWorkflow(li.Name, li.Canonical, li.ID, li.Data)
			} else {
				go job.Handle()
			}
		}
		return
	}

	e.processor.Process(jobs, false)
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
