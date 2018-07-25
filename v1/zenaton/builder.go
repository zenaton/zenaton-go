package zenaton

import "reflect"

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
			workflow: workflow,
		}
	}
	return builderInstance
}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder) Send(eventName string, eventData interface{}) {
	onEventType := reflect.TypeOf(b.workflow.OnEvent)
	if onEventType.In(1) != reflect.TypeOf(eventData) {
		//todo:
		panic("incompatible types")
	}
	b.Client.SendEvent(b.WorkflowClass, b.ID, eventName, eventData)
}

func (b *Builder) WhereID(id string) *Builder {
	b.ID = id
	return b
}
