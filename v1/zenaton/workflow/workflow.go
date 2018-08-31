package workflow

import (
	"reflect"

	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

//func Workfloww (h Handler) func(...interface{}) *Workflow {
//
//	return func (args ...interface{}) *Workflow {
//		h.Init("bob")
//		wf := Workflow{
//			name:    reflect.TypeOf(h).Name(),
//			handler: h,
//		}
//
//		workflowManager := NewWorkflowManager()
//		workflowManager.setClass(wf.name, &wf)
//		return &wf
//	}
//}

func New(h interfaces.Handler) *Workflow {
	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewWorkflow")
	}

	wf := Workflow{
		name:    reflect.TypeOf(h).Elem().Name(),
		handler: h,
	}

	wf.id = wf.GetCustomID()

	NewWorkflowManager().setClass(wf.name, &wf)
	return &wf
}

//todo: must call NewWorkflow to create this so that we can do validation. I'm not sure how to do that given that we must export fields for the worker library
type Workflow struct {
	name    string
	handler interfaces.Handler
	//// must be Exported for worker library
	//OnEvent func(string, interface{}) // todo: have reflect checks on this input type and the data type on Send()
	////todo: in Client.js it says that ID could either be a function or a field
	////todo: ask gilles why id is a function?
	id        string
	canonical string
}

func (wf Workflow) GetName() string { return wf.name }

//todo: change
func (wf Workflow) GetData() interface{} { return wf.handler }

func (wf Workflow) Handle() (interface{}, error) {
	wf.handler.Handle()
	return "", nil
}

func (wf Workflow) Async() error {
	return client.NewClient(false).StartWorkflow(wf.name, wf.canonical, wf.id, wf.handler)
}

func (wf Workflow) OnStart(task *task.Task) {
	starter, ok := wf.handler.(interfaces.Starter)
	if ok {
		starter.Start()
	}
}

func (wf Workflow) OnEvent(eventName string, input interface{}) {
	eventer, ok := wf.handler.(interfaces.OnEventer)
	if ok {
		eventer.OnEvent(eventName, input)
	}
}

func (wf Workflow) OnSuccess(task *task.Task) {
	successer, ok := wf.handler.(interfaces.Successer)
	if ok {
		successer.Success()
	}
}

//todo: should this ever return anything? it seems to be only called form user code and all it does is send an http request. In php for example, there is no return here, just a an exception if the http request doesn't work
func (wf *Workflow) Dispatch2() error {
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
	err := serializer.Decode(encodedData, wf.handler)
	if err != nil {
		panic(err)
	}
	return nil
}

func (wf *Workflow) GetCustomID() string {
	ider, ok := wf.handler.(interfaces.IDer)
	if ok {
		return ider.ID()
	}
	return ""
}

func (wf *Workflow) WhereID(id string) *Builder2 {
	return NewBuilder2(wf).WhereID(id)
}

func (wf *Workflow) Kill() *Builder2 {
	return NewBuilder2(wf).Kill()
}

func (wf *Workflow) Pause() *Builder2 {
	return NewBuilder2(wf).Pause()
}

func (wf *Workflow) Resume() *Builder2 {
	return NewBuilder2(wf).Resume()
}
