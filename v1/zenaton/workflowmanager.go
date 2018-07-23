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

//func (wfm *WorkflowManager) GetWorkflow (name, properties string) *Workflow {
//	// unserialize data
//	const data = serializer.decode(encodedData)
//	// get workflow class
//	var workflowClass = this.getClass(name)
//	// if Version => the workflow was versioned meanwhile => get the initial class
//	if ('VersionClass' === workflowClass.name) {
//	workflowClass = workflowClass.getInitialClass()
//	}
//	// do not use init function to set data
//	workflowClass._useInit = false
//	// return new workflow instance
//	// Object.create(workflowClass);
//	const workflow = new workflowClass(data)
//	// avoid side effect
//	workflowClass._useInit = true
//	// return workflow
//	return workflow
//}
