package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/query"
)

type Workflow struct {
	// todo: maybe unexport these so you have to use the constructor functions?
	Name       string
	HandleFunc func() interface{}
	OnEvent    func(string, interface{}) // todo: should this be an empty interface?
	//todo: in Client.js it says that ID could either be a function or a field
	//todo: do we need this to be a function? because it's kinda funny looking that way
	ID        func() string
	canonical string
}

func New(name string, handlfunc func() interface{}) *Workflow {
	return &Workflow{
		Name:       name,
		HandleFunc: handlfunc,
	}
	//todo: workflowManager.setClass(name, WorkflowClass)
}

func (wf *Workflow) IDFunc(idFunc func() string) *Workflow {
	wf.ID = idFunc
	return wf
}

func (wf *Workflow) WithOnEvent(onEvent func(string, interface{})) *Workflow {
	wf.OnEvent = onEvent
	return wf
}

func (wf *Workflow) Handle() interface{} {
	return wf.HandleFunc()
}

func (wf *Workflow) AsyncHandle(channel chan interface{}) {
	c := client.New(false)

	id := ""
	if wf.ID != nil {
		id = wf.ID()
	}
	channel <- c.StartWorkflow(wf.Name, wf.canonical, id)
}

func (wf *Workflow) Dispatch() {
	e := engine.New()
	e.Dispatch([]engine.Job{wf})
}

func (wf *Workflow) Execute() []interface{} {
	e := engine.New()
	return e.Execute([]engine.Job{wf})
}

// todo: should the Builder returned here be a pointer?
func (wf *Workflow) WhereID(id string) *query.Builder {
	return query.New(wf.Name).WhereID(id)
}

// todo: in js this is with an underscore in front, figure out why and make sure I'm copying functionality
func (wf *Workflow) GetCanonical() string {
	return wf.canonical
}

func (wf *Workflow) SetCanonical(canonical string) {
	wf.canonical = canonical
}
