package zenaton

type Task struct {
	Name       string
	HandleFunc func() interface{}
	Data       interface{}
	ID         func() string
}

func (t *Task) Handle() interface{} {
	return t.HandleFunc()
}

func (t *Task) AsyncHandle(channel chan interface{}) {
	channel <- t.HandleFunc()
}

func (t *Task) Execute() interface{} {
	e := NewEngine()
	return e.Execute([]Job{t})
}

func (t *Task) Dispatch() chan interface{} {
	e := NewEngine()
	return e.Dispatch([]Job{t})[0]
}

func (ts Tasks) Dispatch() chan interface{} {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)[0]
}

type Tasks []*Task

func (ts Tasks) Execute() []interface{} {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetData() interface{} {
	return t.Data
}
