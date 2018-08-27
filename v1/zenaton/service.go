package zenaton

import (
	"errors"

	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

type Service struct {
	Serializer      *serializer.Serializer
	Client          *Client
	Engine          *Engine
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
		Serializer:      &serializer.Serializer{},
		WorkflowManager: NewWorkflowManager(),
		TaskManager:     NewTaskManager(),
		Errors: Errors{
			ScheduledBoxError: ScheduledBoxError,
		},
	}
}
