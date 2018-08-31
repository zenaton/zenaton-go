package workflow

type Versioner interface {
	Version() []*Workflow
}

func NewVersion2(name string, workflows []*Workflow) *Workflow {
	for _, wf := range workflows {
		wf.SetCanonical(name)
	}
	return workflows[len(workflows)-1]
}
