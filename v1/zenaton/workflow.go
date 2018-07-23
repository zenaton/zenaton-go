package zenaton

import "reflect"

type Workflow struct {
	// todo: maybe unexport these so you have to use the constructor functions?
	data       interface{}
	name       string
	handleFunc interface{}
	onEvent    func(string, interface{}) // todo: should this be an empty interface?
	//todo: in Client.js it says that ID could either be a function or a field
	//todo: do we need this to be a function? because it's kinda funny looking that way
	id        func() string
	canonical string
}

func NewWorkflow(name string, handlfunc interface{}) *Workflow {
	t := reflect.TypeOf(handlfunc)
	if t.Kind() != reflect.Func {
		panic("handler argument to NewWorkflow must by of kind function. instead it is of kind: " + t.Kind().String())
	}
	if t.NumIn() > 1 {
		panic("must take a maximum of 1 argument which will receive the data object associated with the workflow: " + t.Kind().String())
	}

	return &Workflow{
		name:       name,
		handleFunc: handlfunc,
	}
	//todo: workflowManager.setClass(name, WorkflowClass)
}

func (wf *Workflow) IDFunc(idFunc func() string) *Workflow {
	wf.id = idFunc
	return wf
}

//todo: allow this to take multiple arguments, so they don't have to build a struct to make this work?
func (wf *Workflow) Data(data interface{}) *Workflow {

	t := reflect.TypeOf(wf.handleFunc)
	if t.NumIn() != 1 {
		panic("if you specify a data field for a workflow, your handler function must have a receiver to accept that data" + t.Kind().String())
	}
	if t.In(0) != reflect.TypeOf(data) {
		panic("type of data must be the same as the parameter type of the handler function. Handler function type: " +
			t.String() + " Data type: " + reflect.TypeOf(data).String())
	}

	wf.data = data
	return wf
}

func (wf *Workflow) WithOnEvent(onEvent func(string, interface{})) *Workflow {
	wf.onEvent = onEvent
	return wf
}

func (wf *Workflow) Handle() interface{} {

	handlFuncValue := reflect.ValueOf(wf.handleFunc)
	handlFuncType := reflect.TypeOf(wf.handleFunc)
	if handlFuncType.NumIn() > 0 {
		in := []reflect.Value{reflect.ValueOf(wf.data)}
		return handlFuncValue.Call(in)
	} else {
		handlFuncValue.Call(nil)
	}
	//todo: fix this so that it actually returns something. I like the way "go.uber.org/cadence/internal/internal_workflow.go:1160" does it
	return nil
}

func (wf *Workflow) AsyncHandle(channel chan interface{}) {
	c := NewClient(false)

	id := ""
	if wf.id != nil {
		id = wf.id()
	}
	channel <- c.StartWorkflow(wf.name, wf.canonical, id)
}

func (wf *Workflow) Dispatch() {
	e := NewEngine()
	e.Dispatch([]Job{wf})
}

func (wf *Workflow) Execute() []interface{} {
	e := NewEngine()
	return e.Execute([]Job{wf})
}

// todo: should the Builder returned here be a pointer?
func (wf *Workflow) WhereID(id string) *Builder {
	return NewBuilder(wf.name).WhereID(id)
}

// todo: in js this is with an underscore in front, figure out why and make sure I'm copying functionality
func (wf *Workflow) GetCanonical() string {
	return wf.canonical
}

func (wf *Workflow) SetCanonical(canonical string) {
	wf.canonical = canonical
}
