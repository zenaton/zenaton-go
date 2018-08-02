package zenaton

type Wait struct {
	task   *Task
	buffer []duration
}

type duration struct {
	kind   string
	amount int64
}

func NewWait(data interface{}) *Wait {
	return &Wait{
		task: &Task{
			name:       "_Wait",
			data:       data,
			handleFunc: func() interface{} { return nil },
		},
	}
}

// todo: shouldn't these path through to the task? maybe because Task is unexported?
func (w *Wait) Handle() (interface{}, error) {
	return w.task.Handle()
}

//func (w *Wait) AsyncHandle(myChan chan interface{}) {
//	w.task.AsyncHandle(myChan)
//}
func (w *Wait) GetName() string {
	return w.task.GetName()
}
func (w *Wait) GetData() interface{} {
	return w.task.GetData()
}

func (w *Wait) Execute() (interface{}, error) {
	return w.task.Execute()
}

func (w *Wait) Seconds(value int64) *Wait {
	w._push(duration{
		kind:   "seconds",
		amount: value,
	})

	return w
}

func (w *Wait) Minutes(value int64) *Wait {
	w._push(duration{
		kind:   "minutes",
		amount: value,
	})

	return w
}

func (w *Wait) Hours(value int64) *Wait {
	w._push(duration{
		kind:   "hours",
		amount: value,
	})

	return w
}

func (w *Wait) Days(value int64) *Wait {
	w._push(duration{
		kind:   "days",
		amount: value,
	})

	return w
}

func (w *Wait) Weeks(value int64) *Wait {
	w._push(duration{
		kind:   "weeks",
		amount: value,
	})

	return w
}

func (w *Wait) months(value int64) *Wait {
	w._push(duration{
		kind:   "months",
		amount: value,
	})

	return w
}

func (w *Wait) Years(value int64) *Wait {
	w._push(duration{
		kind:   "years",
		amount: value,
	})

	return w
}

func (w *Wait) _push(data duration) {
	w.buffer = append(w.buffer, data)
}
