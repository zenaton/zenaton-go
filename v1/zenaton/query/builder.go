package query

import "github.com/zenaton/zenaton-go/v1/zenaton/client"

var instance *Builder

type Builder struct {
	WorkflowClass string
	ID            string
	Client        *client.Client
}

func New(workflowClass string) *Builder {
	if instance == nil {
		instance = &Builder{
			Client:        client.New(false),
			WorkflowClass: workflowClass,
		}
	}
	return instance
}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder) Send(eventName string, eventData interface{}) {
	b.Client.SendEvent(b.WorkflowClass, b.ID, eventName, eventData)
}

func (b *Builder) WhereID(id string) *Builder {
	b.ID = id
	return b
}
