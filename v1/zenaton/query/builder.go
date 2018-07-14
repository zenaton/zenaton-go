package query

import "github.com/zenaton/zenaton-go/v1/zenaton/client"

type Builder struct {
	WorkflowClass string
	ID            string
	Client        *client.Client
}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder) Send(eventName string, eventData interface{}) {
	b.Client.SendEvent(b.WorkflowClass, b.ID, eventName, eventData)
}
