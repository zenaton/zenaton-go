package version

import "github.com/zenaton/zenaton-go/v1/zenaton/workflow"

type Version struct {
	workflows []*workflow.Workflow
}

func New(name string, workflows []*workflow.Workflow) *workflow.Workflow {
	for _, wf := range workflows {
		wf.SetCanonical(name)
	}
	return workflows[len(workflows)-1]
}
