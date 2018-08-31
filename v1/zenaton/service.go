package zenaton

import (
	"errors"

	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

type Service struct {
	Serializer      *serializer.Serializer
	Client          *client.Client
	Engine          *engine.Engine
	WorkflowManager *workflow.WorkflowManager
	TaskManager     *task.TaskManager
	Errors          Errors
}

var ScheduledBoxError = errors.New("ScheduledBoxError")

type Errors struct {
	ScheduledBoxError error
}

func NewService() *Service {
	return &Service{
		Client:          client.NewClient(true),
		Engine:          engine.NewEngine(),
		Serializer:      &serializer.Serializer{},
		WorkflowManager: workflow.NewWorkflowManager(),
		TaskManager:     task.NewTaskManager(),
		Errors: Errors{
			ScheduledBoxError: ScheduledBoxError,
		},
	}
}

func InitClient(appID, apiToken, appEnv string) {
	client.InitClient(appID, apiToken, appEnv)
}

type Workflow = workflow.Workflow
type Task = task.Task
type Wait = task.WaitTask
type Job = interfaces.Job
type Processor = engine.Processor
