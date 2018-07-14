package engine

import "github.com/zenaton/zenaton-go/v1/zenaton/client"

var instance *Engine

type Engine struct {
	client *client.Client
}

//func (e *Engine) ExecuteTask(tasks []task.Task) {
//	//todo: figure out the precessor stuff
//	// local execution
//	//var outputs []interface{} //todo: fix this interface{}
//	//for _, task := range tasks {
//	//	outputs = append(outputs, Handle())
//	//}
//
//	//Engine().execute
//}
//
//func (e *Engine) DispatchTask(tasks []task.Task) {
//}

// todo: does it make sense to use this?
type Job struct {
	Name      string
	Canonical string
}

func (e *Engine) Dispatch(jobs []Job) []interface{} {
	if e.client == nil {
		e.client = &client.Client{}
	}
	var outputs []interface{}
	for _, job := range jobs {
		outputs = append(outputs, e.client.StartWorkflow(job.Name, job.Canonical))
	}
	return outputs
}
