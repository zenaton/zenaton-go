package task

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

var taskManagerInstance *TaskManager

type TaskManager struct {
	tasks map[string]*taskType
}

//todo: problem, This shouldn't be accessible to the user
func NewTaskManager() *TaskManager {
	if taskManagerInstance == nil {
		taskManagerInstance = &TaskManager{
			tasks: make(map[string]*taskType),
		}
	}
	return taskManagerInstance
}

func (tm *TaskManager) setClass(name string, tt *taskType) {
	// check that this task does not exist yet
	//todo: is this right?
	if tm.GetClass(name) != nil {
		panic(fmt.Sprint("Task definition with name '", name, "' already exists"))
	}
	tm.tasks[name] = tt

}

func (tm *TaskManager) GetClass(name string) *taskType {
	return tm.tasks[name]
}

func (tm *TaskManager) GetTask(name, encodedData string) *Task {
	// get task class
	tt := tm.GetClass(name)
	// unserialize data

	err := serializer.Decode(encodedData, tt.defaultTask.Handler)
	if err != nil {
		panic(err)
	}

	//todo: what is this:?
	//// do not use construct function to set data
	//taskClass._useInit = false
	//// get new task instance
	//const task = new taskClass(data)
	//// avoid side effect
	//taskClass._useInit = true
	//// return task

	return tt.defaultTask
}
