package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
)

var builderInstance2 *Builder

type Builder struct {
	WorkflowClass string
	ID            string
	Client        *client.Client
}

func NewBuilder(workflow *WorkflowType) *Builder {
	if builderInstance2 == nil {
		builderInstance2 = &Builder{
			Client:        client.NewClient(false),
			WorkflowClass: workflow.name,
		}
	}
	return builderInstance2
}

func (b *Builder) WhereID(id string) *Builder {
	b.ID = id
	return b
}

func (b *Builder) Find() (*Workflow, error) {
	output, err := b.Client.FindWorkflow(b.WorkflowClass, b.ID)
	if err != nil {
		return nil, err
	}

	properties := output["data"]["properties"].(string)
	name := output["data"]["name"].(string)

	return NewWorkflowManager().GetWorkflow(name, properties), nil
}

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
	b.Client.KillWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Pause a workflow instance
 */

func (b *Builder) Pause() *Builder {
	b.Client.PauseWorkflow(b.WorkflowClass, b.ID)
	return b
}

/**
* Resume a workflow instance
 */

func (b *Builder) Resume() *Builder {
	b.Client.ResumeWorkflow(b.WorkflowClass, b.ID)
	return b
}
