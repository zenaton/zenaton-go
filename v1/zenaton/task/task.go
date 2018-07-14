package task

type Task struct {
	Name   string
	Handle func() interface{}
}

func (t *Task) Execute() {
	t.Handle()
	//Engine().execute
	//e := engine.Engine{}
	//e.Dispatch([]engine.Job{
	//	{
	//		Name: t.Name,
	//		Canonical: t.GetCanonical(),
	//	},
	//})
}

func (t *Task) Dispatch() {
	go t.Handle()
	//Engine().execute
}
