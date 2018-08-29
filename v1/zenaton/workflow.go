package zenaton

import (
	"fmt"
	"reflect"

	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

//todo: must call NewWorkflow to create this so that we can do validation. I'm not sure how to do that given that we must export fields for the worker library
type Workflow struct {
	data       interface{}
	name       string
	handleFunc interface{}
	// must be Exported for worker library
	OnEvent func(string, interface{}) // todo: have reflect checks on this input type and the data type on Send()
	//todo: in Client.js it says that ID could either be a function or a field
	//todo: ask gilles why id is a function?
	id        string
	canonical string
	//todo: what to do with these?
	OnStart               func(*Task)
	OnSuccess             func(*Task, interface{})
	OnFailure             func(*Task, error)
	OnTimeout             func(*Task)
	shouldExecuteNextTask bool
}

type WorkflowParams struct {
	Data       interface{} //required
	Name       string      //required
	HandleFunc interface{}
	//todo: these should take the data object
	OnEvent   func(name string, data interface{})
	OnStart   func(*Task)
	OnSuccess func(*Task, interface{})
	OnFailure func(*Task, error)
	OnTimeout func(*Task)
	// ID recieves the data object and should return a string
	ID interface{}
}

func NewWorkflow(params WorkflowParams) *Workflow {

	validateWorkflowParams(params)

	workflow := &Workflow{
		data:                  params.Data,
		name:                  params.Name,
		handleFunc:            params.HandleFunc,
		OnEvent:               params.OnEvent,
		OnStart:               params.OnStart,
		OnSuccess:             params.OnSuccess,
		OnFailure:             params.OnFailure,
		OnTimeout:             params.OnTimeout,
		shouldExecuteNextTask: true,
	}

	workflowManager := NewWorkflowManager()
	workflowManager.setClass(params.Name, workflow)

	//todo: add this validation to validateWorkflowParams
	if params.ID != nil {
		idType := reflect.TypeOf(params.ID)
		idValue := reflect.ValueOf(params.ID)
		var in []reflect.Value
		if idType.NumIn() > 0 {
			if params.Data == nil {
				panic(fmt.Sprint("workflow: ", params.Name, " has an ID function that specifies an input, but no Data was set on workflow params"))
			}
			in = []reflect.Value{reflect.ValueOf(params.Data)}
		}

		values := idValue.Call(in)

		workflow.id = values[0].String()
	}

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

	validateHandlerFunc(params.HandleFunc, params.Data)
	if params.ID != nil {
		validateIDFunc(params.ID, params.Data)
	}

	return nil
}

// todo: put this somewhere else, as it's used by both workflow and task
func validateHandlerFunc(fn interface{}, data interface{}) {
	fnName := "HandlerFunc"
	fnType := reflect.TypeOf(fn)
	// check that the HandleFunc is in fact a function
	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("%s must be a function, not a : %s", fnName, fnType.Kind().String()))
	}
	// check that the number inputs to the HandleFunc is not greater than 1
	if fnType.NumIn() > 1 {
		panic(fmt.Sprintf("%s must take a maximum of 1 argument", fnName))
	}

	// check that the number of outputs is either 0, 1 or 2 (for either (result, error), or just error, or no return)
	if fnType.NumOut() > 2 {
		panic(fmt.Sprintf("%s must return (result, error) or just error, but found %d return values", fnName, fnType.NumOut()))
	}

	// check that the return type is valid (channels, functions, and unsafe pointers cannot be serialized)
	if fnType.NumOut() > 1 && !isValidResultType(fnType.Out(0)) {
		panic(fmt.Sprintf("%s's first return value cannot be a channel, function, or unsafe pointer; found: %v", fnName, fnType.Out(0).Kind()))
	}

	// check that the last return value is an error
	if fnType.NumOut() > 0 && !isError(fnType.Out(fnType.NumOut()-1)) {
		panic(fmt.Sprintf("expected second return value of %s to return error but found %v", fnName, fnType.Out(fnType.NumOut()-1).Kind()))
	}

	// if Data is defined, the type of Data must be the same as the type of the receiver for HandleFunc
	if data != nil {
		if fnType.NumIn() != 1 {
			panic("if you specify a data field for a task, your function must have a receiver to accept that data")
		}
		if fnType.In(0) != reflect.TypeOf(data) {
			panic(fmt.Sprint("type of data must be the same as the parameter type of the function. ", fnName, "parameter type: ", fnType.String(), " Data type: ", reflect.TypeOf(data).String()))
		}
	}
}

func validateIDFunc(fn interface{}, data interface{}) {
	fnName := "ID function"
	fnType := reflect.TypeOf(fn)
	// check that the HandleFunc is in fact a function
	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("%s must be a function, not a : %s", fnName, fnType.Kind().String()))
	}
	// check that the number inputs to the HandleFunc is not greater than 1
	if fnType.NumIn() > 1 {
		panic(fmt.Sprintf("%s must take a maximum of 1 argument", fnName))
	}

	validateData(fn, data)

	//todo: this is ugly, make it better (at least use a constant)
	//todo: test this
	// check that the ID function returns a string
	if fnType.Out(0).String() != "string" {
		panic(fmt.Sprint(fnName, " should return a string. instead it returns a: ", reflect.TypeOf(fnType.Out(0)).String()))
	}
	if fnType.NumOut() != 1 {
		panic(fmt.Sprint(fnName, " should return a string"))
	}
}

//todo: refactor validation, it's a bit ugly
func validateData(fn interface{}, data interface{}) {
	fnType := reflect.TypeOf(fn)
	// if Data is defined, the type of Data must be the same as the type of the receiver for HandleFunc
	if data != nil {
		if fnType.NumIn() != 1 {
			panic("if you specify a data field for a task, your function must have a receiver to accept that data")
		}
		if fnType.In(0) != reflect.TypeOf(data) {
			fmt.Println("problematic data: ", data)
			panic(fmt.Sprint("type of data must be the same as the parameter type of the function. Function parameter type: ", fnType.In(0).String(), ". Data type: ", reflect.TypeOf(data).String()))
		}
	}
}

func (wf *Workflow) Handle() (interface{}, error) {

	handlFuncValue := reflect.ValueOf(wf.handleFunc)
	handlFuncType := reflect.TypeOf(wf.handleFunc)

	var in []reflect.Value
	if handlFuncType.NumIn() > 0 {
		in = []reflect.Value{reflect.ValueOf(wf.data)}
	}

	//todo: test the error case here
	values := handlFuncValue.Call(in)

	if len(values) == 0 {
		return nil, nil
	}

	if !values[len(values)-1].IsNil() {
		err := values[len(values)-1].Interface().(error)
		return nil, err
	}

	return values[0].Interface(), nil
}

//func (wf *Workflow) AsyncHandle(channel chan interface{}) {
//	c := NewClient(false)
//
//	channel <- c.startWorkflow(wf.name, wf.canonical, wf.GetCustomID(), wf.data)
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

func (wf *Workflow) SetData(data interface{}) *Workflow {
	validateData(wf.handleFunc, data)
	if wf.id != "" {
		validateData(wf.id, data)
	}
	wf.data = data
	return wf
}

func (wf *Workflow) SetDataByEncodedString(encodedData string) error {
	typeHandlerFunc := reflect.TypeOf(wf.handleFunc)
	if typeHandlerFunc.NumIn() > 0 {
		data := reflect.New(typeHandlerFunc.In(0)).Interface()
		err := serializer.Decode(encodedData, data)
		if err != nil {
			return err
		}

		wf.data = reflect.ValueOf(data).Elem().Interface()
	} else {
		wf.data = nil
	}
	return nil
}

func (wf *Workflow) GetCustomID() string {
	return wf.id
}

func (wf *Workflow) WhereID(id string) *Builder {
	return NewBuilder(wf).WhereID(id)
}

func (wf *Workflow) Kill() *Builder {
	return NewBuilder(wf).Kill()
}

func (wf *Workflow) Pause() *Builder {
	return NewBuilder(wf).Pause()
}

func (wf *Workflow) Resume() *Builder {
	return NewBuilder(wf).Resume()
}

func (wf *Workflow) GetHandleFunc() interface{} {
	return wf.handleFunc
}
