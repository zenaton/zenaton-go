package zenaton

var workflowManagerInstance *WorkflowManager

type WorkflowManager struct {
	workflows map[string]*Workflow
}

func NewWorkflowManager() *WorkflowManager {
	if workflowManagerInstance == nil {
		workflowManagerInstance = &WorkflowManager{
			workflows: make(map[string]*Workflow),
		}
	}

	return workflowManagerInstance
}

//todo: what's the difference between these two
func (wfm *WorkflowManager) GetWorkflow(name, encodedData string) *Workflow {
	// get workflow class
	workflow := wfm.GetClass(name)
	// unserialize string data and update the workflow data field

	if encodedData == "" {
		encodedData = "{}"
	}

	err := Serializer{}.Decode(encodedData, &workflow.data)

	if err != nil {
		panic(err)
	}
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
	return wfm.workflows[name]
}

func (wfm *WorkflowManager) setClass(name string, workflow *Workflow) {
	if wfm.GetClass(name) != nil {
		panic(`"` + name + `" workflow can not be defined twice`)
	}
	wfm.workflows[name] = workflow
}
