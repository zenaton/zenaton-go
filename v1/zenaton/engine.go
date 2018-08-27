package zenaton

var instance *Engine

type Engine struct {
	client    *Client
	processor Processor
}

func NewEngine() *Engine {
	if instance == nil {
		instance = &Engine{
			client: NewClient(false),
		}
	}
	return instance
}

// todo: maybe I don't want this to be exported, so only tasks and workflows can implement this interface
type Job interface {
	Handle() (interface{}, error)
	GetName() string
	GetData() interface{}
}

type Processor interface {
	Process([]Job, bool) ([]interface{}, error)
}

type chanResult struct {
	result interface{}
	index  int
}

//todo: error handling
func (e *Engine) Execute(jobs []Job) ([]interface{}, error) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {

		var outputs []interface{}
		var output interface{}
		var err error

		for _, job := range jobs {
			output, err = job.Handle()
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, output)
		}

		return outputs, nil
	}

	return e.processor.Process(jobs, true)
}

func (e *Engine) Dispatch(jobs []Job) error {
	// local execution
	var chans []chan interface{}
	for range jobs {
		chans = append(chans, make(chan interface{}))
	}

	if e.processor == nil || len(jobs) == 0 {

		var err error

		for _, job := range jobs {
			switch v := job.(type) {
			case *Task:
				_, err = job.Handle()
				if err != nil {
					return err
				}
			case *Workflow:
				err = e.client.StartWorkflow(v.name, v.canonical, v.GetCustomID(), v.data)
			}
		}

		return nil
	}

	_, err := e.processor.Process(jobs, false)
	return err
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
