package workflow

import "github.com/zenaton/zenaton-go/v1/zenaton/client"

var builderInstance2 *Builder2

type Builder2 struct {
	workflow      *Workflow
	WorkflowClass string
	ID            string
	Client        *client.Client
}

func NewBuilder2(workflow *Workflow) *Builder2 {
	if builderInstance2 == nil {
		builderInstance2 = &Builder2{
			Client:        client.NewClient(false),
			WorkflowClass: workflow.GetName(),
			workflow:      workflow,
		}
	}
	return builderInstance2
}

func (b *Builder2) WhereID(id string) *Builder2 {
	b.ID = id
	return b
}

//todo:
//func (b *Builder) find() {
//	return b.Client.findWorkflow(b.WorkflowClass, b.ID)
//}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder2) Send(eventName string, eventData interface{}) {
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

func (b *Builder2) Kill() *Builder2 {
	b.Client.KillWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Pause a workflow instance
 */

func (b *Builder2) Pause() *Builder2 {
	b.Client.PauseWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Resume a workflow instance
 */

func (b *Builder2) Resume() *Builder2 {
	b.Client.ResumeWorkflow(b.WorkflowClass, b.ID)
	return b
}
