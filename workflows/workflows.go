package workflows

import (
	"github.com/zenaton/zenaton-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var (
	SequentialWorkflow = workflow.Workflow{
		Name: "SequentialWorkflow",
		Handle: func() {
			tasks.TaskA.Execute()
			tasks.TaskB.Execute()
		},
	}

	//ParallelWorkflow = workflow.Workflow{
	//	Name: "ParallelWorkflow",
	//	Handle: func() {
	//		runParallel := task.Tasks{
	//			tasks.TaskA,
	//			tasks.TaskB,
	//		}
	//		runParallel.Execute()
	//		tasks.TaskC.Execute()
	//	},
	//}

	AsynchronousWorkflow = workflow.Workflow{
		Name: "AsynchronousWorkflow",
		Handle: func() {
			tasks.TaskA.Dispatch()
			tasks.TaskB.Execute()
		},
	}

	EventWorkflow = workflow.Workflow{
		Name: "EventWorkflow",
		Handle: func() {
			tasks.TaskA.Execute()
			tasks.TaskA.Execute()
		},
		OnEvent: func(eventName string, eventData interface{}) {
			if eventName == "MyEvent" {
				tasks.TaskC.Execute()
			}
		},
		ID: func() string {
			return "MyID"
		},
	}
)
