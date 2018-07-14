package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/query"
)

type Workflow struct {
	Name    string
	Handle  func()
	OnEvent func(string, interface{}) // todo: should this be an empty interface?
	//todo: is there a reason why this should be a function instead of a field?
	// todo: should this return a string?
	ID        func() string
	canonical string
}

func (wf *Workflow) Dispatch() {
	e := engine.Engine{}
	e.Dispatch([]engine.Job{
		{
			Name:      wf.Name,
			Canonical: wf.GetCanonical(),
		},
	})
}

func (wf *Workflow) Execute() {
	wf.Handle()
}

// todo: should the Builder returned here be a pointer?
func (wf *Workflow) WhereID(id string) *query.Builder {
	b := query.Builder{
		WorkflowClass: wf.Name,
		ID:            id,
	}
	return &b
}

// todo: in js this is with an underscore in front, figure out why and make sure I'm copying functionality
func (wf *Workflow) GetCanonical() string {
	return wf.canonical
}
