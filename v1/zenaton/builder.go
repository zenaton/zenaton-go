package zenaton

var builderInstance *Builder

type Builder struct {
	WorkflowClass string
	ID            string
	Client        *Client
}

func NewBuilder(workflowClass string) *Builder {
	if builderInstance == nil {
		builderInstance = &Builder{
			Client:        NewClient(false),
			WorkflowClass: workflowClass,
		}
	}
	return builderInstance
}

// do we want to have a different method for each type? or use this empty interface?
func (b *Builder) Send(eventName string, eventData interface{}) {
	b.Client.SendEvent(b.WorkflowClass, b.ID, eventName, eventData)
}

func (b *Builder) WhereID(id string) *Builder {
	b.ID = id
	return b
}
