package wait

import "github.com/zenaton/zenaton-go/v1/zenaton/task"

type Wait struct {
	task   *task.Task
	buffer []duration
}

type duration struct {
	kind   string
	amount int64
}

func New(data interface{}) *Wait {
	return &Wait{
		task: &task.Task{
			Name: "_Wait",
			Data: data,
		},
	}
}

func (w *Wait) Execute() interface{} {
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
