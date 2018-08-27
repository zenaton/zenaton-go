package zenaton

type Version struct {
	workflows []*Workflow
}

func NewVersion(name string, workflows []*Workflow) *Workflow {
	for _, wf := range workflows {
		wf.SetCanonical(name)
	}
	return workflows[len(workflows)-1]
}
