package workflow

import "fmt"

var workflowManagerInstance *WorkflowManager

type WorkflowManager struct {
	workflows2 map[string]*Workflow
}

func NewWorkflowManager() *WorkflowManager {
	if workflowManagerInstance == nil {
		workflowManagerInstance = &WorkflowManager{
			workflows2: make(map[string]*Workflow),
		}
	}

	return workflowManagerInstance
}

func (wfm *WorkflowManager) GetWorkflow(name, encodedData string) *Workflow {

	// get workflow class
	workflow := wfm.GetClass(name)

	if workflow == nil {
		panic(fmt.Sprint("unknown task: ", name,
			". Check that you registered the task with task.Register(&", name, "{}"))
	}

	if encodedData == `""` {
		encodedData = "{}"
	}

	workflow.SetDataByEncodedString(encodedData)

	//todo: figure out this version stuff
	//// if Version => the workflow was versioned meanwhile => get the initial class
	//if "VersionClass" == workflow.name {
	//	workflowClass = workflowClass.getInitialClass()
	//}
	// do not use init function to set data
	//workflowClass._useInit = false
	// return new workflow instance
	// Object.create(workflowClass);
	//const workflow = new workflowClass(data)
	// avoid side effect
	//workflowClass._useInit = true
	// return workflow
	return workflow
}

func (wfm *WorkflowManager) GetClass(name string) *Workflow {
	//fmt.Println("wfm.workflows2", wfm.workflows2)
	return wfm.workflows2[name]
}

func (wfm *WorkflowManager) setClass(name string, workflow *Workflow) {
	if wfm.GetClass(name) == nil {
		wfm.workflows2[name] = workflow
	}
}
