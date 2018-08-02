package zenaton

import "errors"

type Service struct {
	Client          *Client
	Engine          *Engine
	Serializer      *Serializer
	WorkflowManager *WorkflowManager
	TaskManager     *TaskManager
	Errors          Errors
}

var ScheduledBoxError = errors.New("ScheduledBoxError")

type Errors struct {
	ScheduledBoxError error
}

func NewService() *Service {
	return &Service{
		Client:          NewClient(true),
		Engine:          NewEngine(),
		Serializer:      &Serializer{},
		WorkflowManager: NewWorkflowManager(),
		TaskManager:     NewTaskManager(),
		Errors: Errors{
			ScheduledBoxError: ScheduledBoxError,
		},
	}
}
