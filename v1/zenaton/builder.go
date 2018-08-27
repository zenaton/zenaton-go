package zenaton

var builderInstance *Builder

type Builder struct {
	workflow      *Workflow
	WorkflowClass string
	ID            string
	Client        *Client
}

func NewBuilder(workflow *Workflow) *Builder {
	if builderInstance == nil {
		builderInstance = &Builder{
			Client:        NewClient(false),
			WorkflowClass: workflow.name,
			workflow:      workflow,
		}
	}
	return builderInstance
}

func (b *Builder) WhereID(id string) *Builder {
	b.ID = id
	return b
}

//todo:
//func (b *Builder) find() {
//	return b.Client.findWorkflow(b.WorkflowClass, b.ID)
//}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder) Send(eventName string, eventData interface{}) {
	//onEventType := reflect.TypeOf(b.workflow.OnEvent)
	//fmt.Println("onEventType.In(1): ", onEventType.In(1))
	//fmt.Println("reflect.TypeOf(eventData) ", reflect.TypeOf(eventData))
	//if onEventType.In(1) != reflect.TypeOf(eventData) {
	//todo:
	//panic("incompatible types")
	//}
	b.Client.SendEvent(b.WorkflowClass, b.ID, eventName, eventData)
}

/**
 * Kill a workflow instance
 */

func (b *Builder) Kill() *Builder {
	b.Client.killWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Pause a workflow instance
 */

func (b *Builder) Pause() *Builder {
	b.Client.pauseWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Resume a workflow instance
 */

func (b *Builder) Resume() *Builder {
	b.Client.resumeWorkflow(b.WorkflowClass, b.ID)
	return b
}
