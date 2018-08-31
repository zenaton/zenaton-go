package task

import "github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"

var taskManagerInstance *TaskManager

type TaskManager struct {
	tasks map[string]*Task
}

//todo: problem, This shouldn't be accessible to the user
func NewTaskManager() *TaskManager {
	if taskManagerInstance == nil {
		taskManagerInstance = &TaskManager{
			tasks: make(map[string]*Task),
		}
	}
	return taskManagerInstance
}

func (tm *TaskManager) setClass(name string, task *Task) {
	// check that this task does not exist yet
	//todo: is this right?
	if tm.GetClass(name) == nil {
		tm.tasks[name] = task
	}

}

func (tm *TaskManager) GetClass(name string) *Task {
	return tm.tasks[name]
}

func (tm *TaskManager) GetTask(name, encodedData string) *Task {
	// get task class
	task := tm.GetClass(name)
	// unserialize data

	err := serializer.Decode(encodedData, task.handler)
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

	return task
}
