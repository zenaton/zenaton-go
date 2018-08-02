package zenaton

import "reflect"

//todo: must call NewWorkflow to create this so that we can do validation. I'm not sure how to do that given that we must export fields for the worker library
type Workflow struct {
	data       interface{}
	name       string
	handleFunc interface{}
	// must be Exported for worker library
	OnEvent func(string, interface{}) // todo: have reflect checks on this input type and the data type on Send()
	//todo: in Client.js it says that ID could either be a function or a field
	//todo: do we need this to be a function? because it's kinda funny looking that way
	id        func() string
	canonical string
	//todo: what to do with these?
	OnStart   func(*Task)
	OnSuccess func(*Task, interface{})
	OnFailure func(*Task, error)
	OnTimeout func(*Task)
}

type WorkflowParams struct {
	Data interface{} //required
	//todo: get the name from the function name like in cadence?
	Name       string //required
	HandleFunc interface{}
	//todo: these should take the data object
	OnEvent   func(name string, data interface{})
	OnStart   func(*Task)
	OnSuccess func(*Task, interface{})
	OnFailure func(*Task, error)
	OnTimeout func(*Task)
	ID        func() string
}

func NewWorkflow(params WorkflowParams) *Workflow {

	validateWorkflowParams(params)

	workflow := &Workflow{
		data:       params.Data,
		name:       params.Name,
		handleFunc: params.HandleFunc,
		OnEvent:    params.OnEvent,
		id:         params.ID,
		OnStart:    params.OnStart,
		OnSuccess:  params.OnSuccess,
		OnFailure:  params.OnFailure,
		OnTimeout:  params.OnTimeout,
	}

	workflowManager := NewWorkflowManager()
	workflowManager.setClass(params.Name, workflow)

	return workflow
}

//todo: handle onstart, onsuccess, onfailure, ontimeout
//todo: should panic?
func validateWorkflowParams(params WorkflowParams) error {
	if params.Name == "" {
		panic("must set a Name for the workflow")
	}
	if params.HandleFunc == nil {
		panic("must set a HandleFunc for the workflow")
	}

	t := reflect.TypeOf(params.HandleFunc)
	if t.Kind() != reflect.Func {
		panic("handler argument to NewWorkflow must be a function. instead it is of kind: " + t.Kind().String())
	}
	if t.NumIn() > 1 {
		panic("must take a maximum of 1 argument which will receive the data object associated with the workflow: " + t.Kind().String())
	}

	// if Data is defined, the type of Data must be the same as the type of the reciever for HandleFunc
	if params.Data != nil {

		t := reflect.TypeOf(params.HandleFunc)
		if t.NumIn() != 1 {
			panic("if you specify a data field for a workflow, your handler function must have a receiver to accept that data" + t.Kind().String())
		}
		if t.In(0) != reflect.TypeOf(params.Data) {
			panic("type of data must be the same as the parameter type of the handler function. Handler function type: " +
				t.String() + " Data type: " + reflect.TypeOf(params.Data).String())
		}
	}

	return nil
}

func (wf *Workflow) Handle() (interface{}, error) {

	handlFuncValue := reflect.ValueOf(wf.handleFunc)
	handlFuncType := reflect.TypeOf(wf.handleFunc)

	var in []reflect.Value
	if handlFuncType.NumIn() > 0 {
		in = []reflect.Value{reflect.ValueOf(wf.data)}
	}

	values := handlFuncValue.Call(in)

	var err error

	if len(values) == 0 {
		return nil, nil
	}

	if !values[len(values)-1].IsNil() {
		err = values[len(values)-1].Interface().(error)
	}

	if len(values) == 1 {
		return nil, err
	}

	return values[0].Interface(), err
}

//func (wf *Workflow) AsyncHandle(channel chan interface{}) {
//	c := NewClient(false)
//
//	channel <- c.StartWorkflow(wf.name, wf.canonical, wf.GetCustomID(), wf.data)
//}

//todo: should this ever return anything? it seems to be only called form user code and all it does is send an http request. In php for example, there is no return here, just a an exception if the http request doesn't work
func (wf *Workflow) Dispatch() error {
	e := NewEngine()
	return e.Dispatch([]Job{wf})
}

func (wf *Workflow) Execute() ([]interface{}, error) {
	e := NewEngine()
	return e.Execute([]Job{wf})
}

func (wf *Workflow) WhereID(id string) *Builder {
	return NewBuilder(wf).WhereID(id)
}

// todo: in js this is with an underscore in front, figure out why and make sure I'm copying functionality
func (wf *Workflow) GetCanonical() string {
	return wf.canonical
}

func (wf *Workflow) SetCanonical(canonical string) {
	wf.canonical = canonical
}

func (wf *Workflow) GetName() string {
	return wf.name
}

func (wf *Workflow) GetData() interface{} {
	return wf.data
}

func (wf *Workflow) SetData(data interface{}) {
	wf.data = data
}

func (wf *Workflow) GetCustomID() string {
	var id string
	if wf.id != nil {
		id = wf.id()
	}
	return id
}
