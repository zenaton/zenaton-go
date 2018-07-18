package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/query"
)

type Workflow struct {
	Name       string
	HandleFunc func() interface{}
	OnEvent    func(string, interface{}) // todo: should this be an empty interface?
	//todo: in Client.js it says that ID could either be a function or a field
	ID        func() string
	canonical string
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
