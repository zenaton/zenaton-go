package workflow

import (
	"reflect"

	"encoding/json"

	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

type Workflow struct {
	name string
	interfaces.Handler
	canonical string
	id        string
}

type WorkflowType struct {
	name            string
	defaultWorkflow *Workflow
	canonical       string
}

func (wft *WorkflowType) SetCanonical(canonical string) {
	wft.canonical = canonical
}

func (wft *WorkflowType) NewInstance(handlers ...interfaces.Handler) *Workflow {

	if len(handlers) > 1 {
		panic("must only pass one handler to WorkflowType.NewInstance()")
	}

	if len(handlers) == 1 {
		h := handlers[0]
		validateHandler(h)
		return newInstance(wft.name, h)
	} else {
		return wft.defaultWorkflow
	}
}

type defaultHandler struct{
	handlerFunc func() (interface{}, error)
}
func (dh *defaultHandler) Handle () (interface{}, error){
	return dh.handlerFunc()
}

func NewDefault(name string, handlerFunc func() (interface{}, error)) *WorkflowType {
	return New(name, &defaultHandler{
		handlerFunc: handlerFunc,
	})
}

func New(name string, h interfaces.Handler) *WorkflowType {

	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewWorkflow")
	}

	validateHandler(h)

	wf := WorkflowType{
		name:            name,
		defaultWorkflow: newInstance(name, h),
	}

	NewWorkflowManager().setClass(wf.name, &wf)
	return &wf
}

func validateHandler(value interface{}) {

	name := reflect.Indirect(reflect.ValueOf(value)).Type().Name()

	jsonV, err := json.Marshal(value)
	if err != nil {
		panic("handler type '" + name + "' must be able to be marshaled to json. " + err.Error())
	}

	newV := reflect.New(reflect.TypeOf(value)).Interface()

	err = json.Unmarshal(jsonV, newV)
	if err != nil {
		panic("handler type '" + name + "' must be able to be unmarshaled from json. " + err.Error())
	}
}

func newInstance(name string, h interfaces.Handler) *Workflow {
	return &Workflow{
		name:    name,
		Handler: h,
	}
}

func (wf Workflow) GetName() string      { return wf.name }
func (wf Workflow) GetData() interface{} { return wf.Handler }

func (wf Workflow) Async() error {
	return client.NewClient(false).StartWorkflow(wf.name, wf.canonical, wf.GetCustomID(), wf.Handler)
}

func (wf Workflow) OnStart(task *task.Task) {
	starter, ok := wf.Handler.(interfaces.Starter)
	if ok {
		starter.Start()
	}
}

func (wf Workflow) OnEvent(eventName string, input interface{}) {

	eventer, ok := wf.Handler.(interfaces.OnEventer)
	if ok {
		eventer.OnEvent(eventName, input)
	}
}

func (wf Workflow) OnSuccess(task *task.Task, output interface{}) {
	successer, ok := wf.Handler.(interfaces.Successer)
	if ok {
		successer.Success()
	}
}

//todo: should this ever return anything? it seems to be only called form user code and all it does is send an http request. In php for example, there is no return here, just a an exception if the http request doesn't work
func (wf *Workflow) Dispatch() error {
	e := engine.NewEngine()
	return e.Dispatch([]interfaces.Job{wf})
}

// todo: in js this is with an underscore in front, figure out why and make sure I'm copying functionality
func (wf *Workflow) GetCanonical() string {
	return wf.canonical
}

func (wf *Workflow) SetCanonical(canonical string) {
	wf.canonical = canonical
}

func (wf *Workflow) SetDataByEncodedString(encodedData string) error {
	err := serializer.Decode(encodedData, wf.Handler)
	if err != nil {
		panic(err)
	}
	return nil
}

func (wf *Workflow) GetCustomID() string {
	ider, ok := wf.Handler.(interfaces.IDer)
	if ok {
		return ider.ID()
	}
	return ""
}

func (wft *WorkflowType) WhereID(id string) *QueryBuilder {
	return NewBuilder(wft).WhereID(id)
}
