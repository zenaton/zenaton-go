package workflow

//todo:
func Version(name string, workflowTypes []*WorkflowType) *WorkflowType {
	for _, wf := range workflowTypes {
		wf.SetCanonical(name)
	}
	return workflowTypes[len(workflowTypes)-1]
}
