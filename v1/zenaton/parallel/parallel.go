package parallel

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/internal/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

// Jobs is just a slice of Jobs that can be run in parallel with either Execute or Dispatch. A Job is either a task
// Instance or a workflow Instance.
type Jobs []engine.Job

// Execute will execute a Jobs (a slice of Instances) in parallel and wait for their completion.
// Execute returns an Execution which can be used to get the outputs and errors of the tasks executed.
//
// For parallel tasks, you will receive a slice of errors. This slice will be nil if no error occurred. If there was
// an error in one of the parallel tasks, you will receive a slice of the same length as the input tasks, and the index
// of the task that produced an error will be the same index as the non-nil err in the slice of errors
//
//  var a int
//	var b int
//
//	errs := task.Jobs{
//	    tasks.A.New(),
//	    tasks.B.New(),
//	}.Execute().Output(&a, &b)
//
//	if errs != nil {
//	    if errs[0] != nil {
//	        //  tasks.A error
//	    }
//	    if errs[1] != nil {
//	        //  tasks.B error
//	    }
//	}
//
// Here, tasks A and B will be executed in parallel, and we wait for all of them to end before continuing. You can
// retrieve the outputs of these tasks by passing pointers to .Output()
func (js Jobs) Execute() Execution {

	e := engine.NewEngine()
	values, serializedValues, errors := e.Execute(js)

	return Execution{
		outputValues:     values,
		serializedValues: serializedValues,
		errors:           errors,
	}
}

// Dispatch will launch the the tasks in parallel and not wait for them to complete before moving on. Thus:
//
//		task.Jobs{
//			tasks.A.New(),
//			tasks.B.New(),
//		}.Dispatch()
//
// should be equivalent to:
//
// 		tasks.A.New().Dispatch()
// 		tasks.B.New().Dispatch()
//
func (js Jobs) Dispatch() {
	e := engine.NewEngine()

	var launchInfos []engine.LaunchInfo
	for _, j := range js {
		switch v := j.(type) {
		case *workflow.Instance:
			li := engine.LaunchInfo{
				Type:      "workflow",
				Name:      v.GetName(),
				Canonical: v.GetCanonical(),
				ID:        v.GetCustomID(),
				Data:      v.GetData(),
			}
			launchInfos = append(launchInfos, li)
		case *task.Instance:
			li := engine.LaunchInfo{
				Type:      "task",
			}
			launchInfos = append(launchInfos, li)
		default:
		}
	}

	e.Dispatch(js, launchInfos)
}

// Execution represents the outputs and errors of the Jobs tasks.
// To get the output, use Execution.Output()
type Execution struct {
	outputValues     []interface{}
	serializedValues []string
	errors           []error
}

// Output gets the output of a Jobs execution
//
//
// For parallel tasks, you will receive a slice of errors. This slice will be nil if no error occurred. If there was
// an error in one of the parallel tasks, you will receive a slice of the same length as the input tasks, and the index
// of the task that produced an error will be the same index as the non-nil err in the slice of errors
//
//  var a int
//	var b int
//
//	errs := task.Jobs{
//	    tasks.A.New(),
//	    tasks.B.New(),
//	}.Execute().Output(&a, &b)
//
//	if errs != nil {
//	    if errs[0] != nil {
//	        //  tasks.A error
//	    }
//	    if errs[1] != nil {
//	        //  tasks.B error
//	    }
//	}
//
// Here, tasks A and B will be executed in parallel, and we wait for all of them to end before continuing. You can
// retrieve the outputs of these tasks by passing pointers to .Output()
func (e Execution) Output(values ...interface{}) []error {

	if len(values) != len(e.outputValues) && len(values) != len(e.serializedValues) {
		panic(fmt.Sprint("task: number of parallel tasks and return value pointers do not match"))
	}

	if len(values) == 0 {
		values = make([]interface{}, len(e.errors))
	}

	var errs []error

	if e.serializedValues != nil {
		for i := range e.serializedValues {
			err := outputFromSerialized(values[i], e.serializedValues[i])
			errs = append(errs, err)
		}
	} else {

		for i := range e.outputValues {
			if values[i] != nil {
				value := values[0]
				outputFromInterface(value, e.outputValues[i])
			}
		}

		errs = e.errors
	}

	for _, err := range errs {
		if err != nil {
			return e.errors
		}
	}
	return nil
}

func outputFromSerialized(to interface{}, from string) error {

	rv := reflect.ValueOf(to)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprint("must pass a non-nil pointer to task.Output"))
	}

	var combinedOutput map[string]json.RawMessage

	err := serializer.Decode(from, &combinedOutput)
	if err != nil {
		panic(err)
	}

	if to != nil {
		err = serializer.Decode(string(combinedOutput["output"]), to)
		if err != nil {
			panic(err)
		}
	}

	if combinedOutput["error"] != nil {
		return errors.New(string(combinedOutput["error"]))
	}
	return nil
}

func outputFromInterface(to, from interface{}) {
	rv := reflect.ValueOf(to)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprint("must pass a non-nil pointer to task.Output"))
	}

	if from != nil && to != nil {
		outV := reflect.ValueOf(from)
		if outV.IsValid() {
			rv.Elem().Set(outV)
		}
	}
}
