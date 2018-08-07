package zenaton

import "time"

type Wait struct {
	task   *Task
	event  string
	buffer []duration
}

type duration struct {
	kind   string
	amount int64
}

func NewWait() *Wait {
	time.Now().Add(1*time.Minute + 1*time.Second)
	return &Wait{
		task: NewTask(TaskParams{
			Name:       "Wait",
			HandleFunc: func() {},
		}),
	}
}

func (w *Wait) WithEvent(event string) *Wait {
	w.event = event
	return w
}

func (w *Wait) Event() string {
	return w.event
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
	w.push(duration{
		kind:   "seconds",
		amount: value,
	})

	return w
}

func (w *Wait) Minutes(value int64) *Wait {
	w.push(duration{
		kind:   "minutes",
		amount: value,
	})

	return w
}

func (w *Wait) Hours(value int64) *Wait {
	w.push(duration{
		kind:   "hours",
		amount: value,
	})

	return w
}

func (w *Wait) Days(value int64) *Wait {
	w.push(duration{
		kind:   "days",
		amount: value,
	})

	return w
}

func (w *Wait) Weeks(value int64) *Wait {
	w.push(duration{
		kind:   "weeks",
		amount: value,
	})

	return w
}

func (w *Wait) months(value int64) *Wait {
	w.push(duration{
		kind:   "months",
		amount: value,
	})

	return w
}

func (w *Wait) Years(value int64) *Wait {
	w.push(duration{
		kind:   "years",
		amount: value,
	})

	return w
}

func (w *Wait) push(data duration) {
	w.buffer = append(w.buffer, data)
}

//
//func (w *Wait) initNowThen() {
//	// get setted or current time zone
//	var tz = undefined !== w.constructor._timezone ? w.constructor._timezone : moment.tz.guess();
//	var now = moment().tz(tz);
//	var then = moment(now);
//	return [now, then];
//},
//
//
//func (w *Wait) getTimestampOrDuration() {
//
//	//todo: is this right? should it be nil or and empty slice?
//	if w.buffer == nil {
//		//todo:
//		//return [null, null]
//	}
//
//	var initNowThen = w.initNowThen(),
//	initNowThen2 = slicedToArray(initNowThen, 2),
//	now = initNowThen2[0],
//	then = initNowThen2[1]
//
//	w.mode = null // apply buffered methods
//
//	w.buffer.forEach(function (call) {
//	then = w.apply(call[0], call[1], now, then)
//	}) // has user used a method by timestamp?
//
//
//	var isTimestamp = w.mode !== null // remove attribute to avoid having it in linearization
//
//	delete w.mode //return
//
//	if isTimestamp {
//	return [then.unix(), null]
//	} else {
//	return [null, then.diff(now, 'seconds')]
//}
//},
//
//timestamp(value) {
//w.push(['timestamp', value])
//
//return w
//},
//
//at(value) {
//w.push(['at', value])
//
//return w
//},
//
//dayOfMonth(value) {
//w.push(['dayOfMonth', value])
//
//return w
//},
//
//monday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['monday', value])
//
//return w
//},
//
//tuesday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['tuesday', value])
//
//return w
//},
//
//wednesday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['wednesday', value])
//
//return w
//},
//
//thursday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['thursday', value])
//
//return w
//},
//
//friday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['friday', value])
//
//return w
//},
//
//saturday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['saturday', value])
//
//return w
//},
//
//sunday() {
//var value = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : 1
//
//w.push(['sunday', value])
//
//return w
//},
//
//timestamp(timestamp) {
//w.setMode(MODETIMESTAMP)
//
//return moment.unix(timestamp)
//},
//
//at(time, now, then) {
//w.setMode(MODEAT)
//
//var segments = time.split(':')
//var h = parseInt(segments[0])
//var m = segments.length > 1 ? parseInt(segments[1]) : 0
//var s = segments.length > 2 ? parseInt(segments[2]) : 0
//then = then.set({
//'hour': h,
//'minute': m,
//'second': s
//})
//
//if now.isAfter(then) {
//switch (w.mode) {
//case MODEAT:
//then = then.add(1, 'days')
//break
//
//case MODEWEEKDAY:
//then = then.add(1, 'weeks')
//break
//
//case MODEMONTHDAY:
//then = then.add(1, 'months')
//break
//
//default:
//throw new InternalZenatonError('Unknown mode: ' + w.mode)
//}
//}
//
//return then
//},
//
//dayOfMonth(day, now, then) {
//w.setMode(MODEMONTHDAY)
//
//then = then.set('date', day)
//
//if now.isAfter(then) {
//then = then.add(1, 'months')
//}
//
//return then
//},
//
//weekDay(n, day, then) {
//w.setMode(MODEWEEKDAY)
//
//var d = then.isoWeekday()
//then = then.add(day - d, 'days')
//then = d > day ? then.add(n, 'weeks') : then.add(n - 1, 'weeks')
//return then
//},
//
//apply(method, value, now, then) {
//switch (method) {
//case 'timestamp':
//return w.timestamp(value, then)
//
//case 'at':
//return w.at(value, now, then)
//
//case 'dayOfMonth':
//return w.dayOfMonth(value, now, then)
//
//case 'monday':
//return w.weekDay(value, 1, then)
//
//case 'tuesday':
//return w.weekDay(value, 2, then)
//
//case 'wednesday':
//return w.weekDay(value, 3, then)
//
//case 'thursday':
//return w.weekDay(value, 4, then)
//
//case 'friday':
//return w.weekDay(value, 5, then)
//
//case 'saturday':
//return w.weekDay(value, 6, then)
//
//case 'sunday':
//return w.weekDay(value, 7, then)
//
//default:
//return w.applyDuration(method, value, then)
//}
//},
//
//setMode(mode) {
//// can not apply twice the same method
//if mode == w.mode {
//throw new ExternalZenatonError('Incompatible definition in Wait methods')
//} // timestamp can only be used alone
//
//
//if null !== w.mode && MODETIMESTAMP == mode || MODETIMESTAMP == w.mode {
//throw new ExternalZenatonError('Incompatible definition in Wait methods')
//} // other mode takes precedence to MODEAT
//
//
//if null == w.mode || MODEAT == w.mode {
//w.mode = mode
//}
//}
