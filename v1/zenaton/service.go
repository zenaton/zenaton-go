package zenaton

type Service struct{
	Client *Client
	Engine *Engine
	Serializer *Serializer
	WorkflowManager *WorkflowManager
	TaskManager *TaskManager
}

func NewService () *Service {
	return &Service{
		Client: NewClient(true),
		Engine: NewEngine(),
		Serializer: &Serializer{},
		WorkflowManager: NewWorkflowManager(),
		TaskManager: NewTaskManager(),
	}
}
