package workflow

import "fmt"

var workflowManagerInstance *WorkflowManager

type WorkflowManager struct {
	workflows map[string]*WorkflowType
}

func NewWorkflowManager() *WorkflowManager {
	if workflowManagerInstance == nil {
		workflowManagerInstance = &WorkflowManager{
			workflows: make(map[string]*WorkflowType),
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

	workflow.defaultWorkflow.SetDataByEncodedString(encodedData)

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
	return workflow.defaultWorkflow
}

func (wfm *WorkflowManager) GetClass(name string) *WorkflowType {
	//fmt.Println("wfm.workflows", wfm.workflows)
	return wfm.workflows[name]
}

func (wfm *WorkflowManager) setClass(name string, workflow *WorkflowType) {
	if wfm.GetClass(name) != nil {
		panic(fmt.Sprint("Workflow definition with name '", name, "' already exists"))
	}
	wfm.workflows[name] = workflow
}
