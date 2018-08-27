package zenaton

import "reflect"

type Task struct {
	name string
	//todo: would be nice if the handle func could take many arguments, instead of just one. would have to think how that would be done (maybe pass in argments into execute?)
	handleFunc        interface{}
	data              interface{}
	id                interface{}
	MaxProcessingTime func() int64
}

type TaskParams struct {
	Name       string
	HandleFunc interface{}
	//MaxProcessingTime is an optional function that sets the maximum allowed processing time (in seconds) for the task
	MaxProcessingTime func() int64
	Data              interface{}
	ID                interface{}
}

func NewTask(params TaskParams) *Task {
	validateTaskParams(params)

	task := &Task{
		name:       params.Name,
		handleFunc: params.HandleFunc,
		data:       params.Data,
		id:         params.ID,
	}

	NewTaskManager().setClass(params.Name, task)
	return task
}

//todo: should I panic here?
//todo: should I really require that you return an error?
func validateTaskParams(params TaskParams) {
	if params.Name == "" {
		panic("must set a Name for the task")
	}
	if params.HandleFunc == nil {
		panic("must set a HandleFunc for the task")
	}

	validateHandlerFunc(params.HandleFunc, params.Data)
	if params.ID != nil {
		validateIDFunc(params.ID, params.Data)
	}
}

func isError(inType reflect.Type) bool {
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	return inType.Implements(errorInterface)
}

//todo: I'm not sure if I actually need this. I don't think the output's are actually serialized.
func isValidResultType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return false
	}

	return true
}

func (t *Task) Handle() (interface{}, error) {

	handlFuncValue := reflect.ValueOf(t.handleFunc)
	handlFuncType := reflect.TypeOf(t.handleFunc)

	var in []reflect.Value
	if handlFuncType.NumIn() > 0 {
		in = []reflect.Value{reflect.ValueOf(t.data)}
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

//todo: would be great if we could take a pointer to execute and modify that like json.unmarshal does, but it's hard to figure out how they do it
func (t *Task) Execute() (interface{}, error) {
	e := NewEngine()
	output, err := e.Execute([]Job{t})
	//todo: make sure this is impossible to get index out of bounds
	if output == nil {
		return nil, err
	}
	return output[0], err
}

func (t *Task) Dispatch() error {
	e := NewEngine()
	err := e.Dispatch([]Job{t})
	return err
}

func (ts Tasks) Dispatch() error {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)
}

type Tasks []*Task

func (ts Tasks) Execute() ([]interface{}, error) {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) GetData() interface{} {
	return t.data
}

func (t *Task) SetData(data interface{}) *Task {
	t.data = data
	return t
}
